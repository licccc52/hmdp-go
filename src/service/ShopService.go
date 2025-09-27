package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hmdp/src/config/mysql"
	redisClient "hmdp/src/config/redis"
	"hmdp/src/model"
	"hmdp/src/utils"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	redisConfig "github.com/redis/go-redis/v9"
)

type ShopService struct {
}

var ShopManager *ShopService

const (
	MAX_REDIS_DATA_QUEUE = 10
)

var redisDataQueue chan int64

func init() {
	redisDataQueue = make(chan int64, MAX_REDIS_DATA_QUEUE)
	go ShopManager.SyncUpdateCache()
}

func (*ShopService) QueryShopById(id int64) (model.Shop, error) {
	var shop model.Shop
	shop.Id = id
	err := shop.QueryShopById(id)
	return shop, err
}

func (*ShopService) SaveShop(shop *model.Shop) error {
	err := shop.SaveShop()
	return err
}

func (*ShopService) UpdateShop(shop *model.Shop) error {
	err := shop.UpdateShop(mysql.GetMysqlDB())
	return err
}

func (*ShopService) QueryByType(typeId int, current int) ([]model.Shop, error) {
	var shopUtils model.Shop
	shops, err := shopUtils.QueryShopByType(typeId, current)
	return shops, err
}

func (*ShopService) QueryByName(name string, current int) ([]model.Shop, error) {
	var shopUtils model.Shop
	shops, err := shopUtils.QueryShopByName(name, current)
	return shops, err
}

func (*ShopService) QueryShopByIdWithCache(id int64) (model.Shop, error) {
	redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(id, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shopInfo, err := redisClient.GetRedisClient().Get(ctx, redisKey).Result()
	if err == nil {
		var shop model.Shop
		err = json.Unmarshal([]byte(shopInfo), &shop)
		if err != nil {
			return model.Shop{}, err
		}
		return shop, nil
	}

	if err == redisConfig.Nil {
		var shop model.Shop
		shop.Id = id
		err = shop.QueryShopById(id)
		if err != nil {
			return model.Shop{}, err
		}

		redisValue, err := json.Marshal(shop)
		if err != nil {
			return model.Shop{}, err
		}

		// 超时剔除策略
		err = redisClient.GetRedisClient().Set(ctx, redisKey, string(redisValue), time.Duration(time.Minute)).Err()

		if err != nil {
			return model.Shop{}, err
		}
		return shop, nil
	}

	return model.Shop{}, err
}

// 缓存更新的最佳实践方法
func (*ShopService) UpdateShopWithCacheCallBack(db *gorm.DB, shop *model.Shop) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := shop.QueryShopById(shop.Id)
		if err != nil {
			return err
		}

		// update the database
		err = shop.UpdateShop(tx)
		if err != nil {
			return err
		}

		// delete the cache
		redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(shop.Id, 10)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err = redisClient.GetRedisClient().Del(ctx, redisKey).Err()

		if err != nil {
			return err
		}

		return nil
	})
}

func (*ShopService) UpdateShopWithCache(shop *model.Shop) error {
	return ShopManager.UpdateShopWithCacheCallBack(mysql.GetMysqlDB(), shop)
}

// 缓存穿透的解决方法: 缓存空对象
func (*ShopService) QueryShopByIdWithCacheNull(id int64) (model.Shop, error) {
	redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(id, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shopInfoStr, err := redisClient.GetRedisClient().Get(ctx, redisKey).Result()

	if err == nil {
		var shopInfo model.Shop
		if shopInfoStr == "" {
			return model.Shop{}, nil
		}
		err = json.Unmarshal([]byte(shopInfoStr), &shopInfo)
		if err != nil {
			return model.Shop{}, err
		}
		return shopInfo, nil
	}

	if err == redisConfig.Nil {
		var shopInfo model.Shop
		shopInfo.Id = id
		err = shopInfo.QueryShopById(id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = redisClient.GetRedisClient().Set(ctx, redisKey, "", time.Duration(time.Minute)).Err()
			if err != nil {
				return model.Shop{}, err
			}
			return model.Shop{}, nil
		}

		redisValue, err := json.Marshal(shopInfo)
		if err != nil {
			return model.Shop{}, err
		}
		err = redisClient.GetRedisClient().Set(ctx, redisKey, string(redisValue), time.Duration(time.Minute)).Err()
		if err != nil {
			return model.Shop{}, err
		}
		return shopInfo, nil
	}
	return model.Shop{}, nil
}

// 利用互斥锁的方式解决缓存击穿的问题
func (*ShopService) QueryShopByIDWithMutex(id int64, retry int) (model.Shop, error) {
	if retry > 3 {
		return model.Shop{}, errors.New("获取锁重试次数超限")
	}

	redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(id, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	info, err := redisClient.GetRedisClient().Get(ctx, redisKey).Result()
	if err == nil {
		var shopInfo model.Shop
		err = json.Unmarshal([]byte(info), &shopInfo)
		if err != nil {
			return model.Shop{}, err
		}
		return shopInfo, nil
	}
	if err == redisConfig.Nil {
		lock := utils.CACHE_LOCK_KEY + strconv.FormatInt(id, 10)
		flag, _ := redisClient.GetRedisClient().SetNX(ctx, lock, "1", 1*time.Second).Result()
		if !flag {
			time.Sleep(50 * time.Millisecond)
			return ShopManager.QueryShopByIDWithMutex(id, retry+1)
		}
		defer func() {
			//删除mutex的key
			err = redisClient.GetRedisClient().Del(ctx, lock).Err()
			if err != nil {
				fmt.Println("ShopService Error,  delete cache_lock_key failed ", err.Error())
			}
		}()

		var shopInfo model.Shop
		err = shopInfo.QueryShopById(id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = redisClient.GetRedisClient().Set(ctx, redisKey, "", 10*time.Second).Err()
			return model.Shop{}, err
		}

		respStr, err := json.Marshal(shopInfo)
		if err != nil {
			return model.Shop{}, err
		}

		err = redisClient.GetRedisClient().Set(ctx, redisKey, string(respStr), 10*time.Second).Err()
		if err != nil {
			return model.Shop{}, err
		}
		return shopInfo, nil

	}

	return model.Shop{}, nil
}

// 利用互斥锁的方式解决缓存击穿的问题
func (*ShopService) QueryShopByIDWithLogicExpire(id int64) (model.Shop, error) {

	redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(id, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	info, err := redisClient.GetRedisClient().Get(ctx, redisKey).Result()
	var shopInfo model.Shop
	if err == nil {
		err = json.Unmarshal([]byte(info), &shopInfo)
		if err != nil {
			return model.Shop{}, err
		}

		if shopInfo.ExpireTime.Before(time.Now()) {
			//缓存已经 逻辑过期
			lock := utils.CACHE_LOCK_KEY + strconv.FormatInt(id, 10)
			flag, _ := redisClient.GetRedisClient().SetNX(ctx, lock, "1", 1*time.Second).Result()
			if !flag {
				return shopInfo, nil
			}
			defer func() {
				err = redisClient.GetRedisClient().Del(ctx, lock).Err()
				if err != nil {
					fmt.Println("ShopService Error,  delete cache_lock_key failed ", err.Error())
				}
			}()

			go func() { //根据id查询数据库

				goroutineCtx := context.Background()
				err = shopInfo.QueryShopById(id)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = redisClient.GetRedisClient().Set(goroutineCtx, redisKey, "", 10*time.Second).Err()
				}

				if err != nil {
					return
				}
				shopInfo.ExpireTime = time.Now().Add(10 * time.Second)

				shopInfoStr, err := json.Marshal(shopInfo)
				if err != nil {
					return
				}
				err = redisClient.GetRedisClient().Set(goroutineCtx, redisKey, string(shopInfoStr), 10*time.Minute).Err()
				if err != nil {
					return
				}
				return
			}()
		}
		return shopInfo, nil
	}

	if err == redisConfig.Nil {
		var shopInfo model.Shop
		err = shopInfo.QueryShopById(id)
		shopInfo.ExpireTime = time.Now().Add(10 * time.Second)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = redisClient.GetRedisClient().Set(ctx, redisKey, "", 10*time.Minute).Err()
		}
		if err != nil {
			return model.Shop{}, err
		}
		shopInfoStr, err := json.Marshal(shopInfo)
		if err != nil {
			return model.Shop{}, err
		}
		_ = redisClient.GetRedisClient().Set(ctx, redisKey, string(shopInfoStr), 10*time.Minute).Err()
		return shopInfo, nil
	}

	return model.Shop{}, nil
}

// @Description: use the logic expire to deal with the cache pass through
func (*ShopService) QueryShopByIdWithLogicExpire(id int64) (model.Shop, error) {
	redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(id, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisDataStr, err := redisClient.GetRedisClient().Get(ctx, redisKey).Result()

	// hot key is in redis
	if err == redisConfig.Nil {
		return model.Shop{}, nil
	}

	if err == nil {
		if redisDataStr == "" {
			return model.Shop{}, nil
		}

		var redisData utils.RedisData[model.Shop]
		err = json.Unmarshal([]byte(redisDataStr), &redisData)
		if err != nil {
			return model.Shop{}, err
		}

		if redisData.ExpireTime.After(time.Now()) {
			return redisData.Data, nil
		}
		// 否则过期,需要重新建立缓存

		lockKey := utils.CACHE_LOCK_KEY + strconv.FormatInt(id, 10)
		flag := utils.RedisUtil.TryLock(lockKey)

		// if not get the lock
		if !flag {
			return redisData.Data, nil
		}

		// if get the lock
		defer utils.RedisUtil.ClearLock(lockKey)
		redisDataQueue <- id
		// go func() {
		// 	var shopInfo model.Shop
		// 	err = shopInfo.QueryShopById(id)
		// 	if err != nil {
		// 		return
		// 	}
		// 	var redisDataToSave utils.RedisData[model.Shop]
		//
		// 	redisDataToSave.Data = shopInfo
		// 	// the time of hot key exists
		// 	redisDataToSave.ExpireTime = time.Now().Add(time.Second * utils.HOT_KEY_EXISTS_TIME)
		//
		// 	redisValue,err := json.Marshal(redisDataToSave)
		// 	err = redisClient.GetRedisClient().Set(ctx , redisKey , string(redisValue) , 0).Err()
		// 	if err != nil {
		// 		return
		// 	}
		// 	return
		// }()
		//
		return redisData.Data, nil
	}

	return model.Shop{}, err
}

func (*ShopService) SyncUpdateCache() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		id := <-redisDataQueue

		redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(id, 10)

		var shopInfo model.Shop
		err := shopInfo.QueryShopById(id)

		if err != nil {
			continue
		}

		var redisDataToSave utils.RedisData[model.Shop]

		redisDataToSave.Data = shopInfo
		// the time of hot key exists
		redisDataToSave.ExpireTime = time.Now().Add(time.Second * utils.HOT_KEY_EXISTS_TIME)

		redisValue, err := json.Marshal(redisDataToSave)
		err = redisClient.GetRedisClient().Set(ctx, redisKey, string(redisValue), 0).Err()
		if err != nil {
			continue
		}
	}
}
