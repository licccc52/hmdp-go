package service

import (
	"hmdp/src/config/mysql"
	"hmdp/src/config/redis"
	"hmdp/src/model"
	"sync"
	"testing"
)

// 建立热点 key
// func TestHotKey(t *testing.T) {
// 	redis.Init()
// 	mysql.Init()
//
// 	var id int64 = 16
// 	var shopInfo model.Shop
// 	err := shopInfo.QueryShopById(id)
// 	if err != nil {
// 		t.Fatal(err.Error())
// 		return
// 	}
//
// 	var redisData utils.RedisData[model.Shop]
// 	redisData.Data = shopInfo
// 	redisData.ExpireTime = time.Now().Add(time.Second * 20)
// 	redisValue , err := json.Marshal(redisData)
// 	if err != nil {
// 		t.Fatal(err.Error())
// 		return
// 	}
//
// 	ctx , cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	redisKey := utils.CACHE_SHOP_KEY + strconv.FormatInt(id,10)
// 	err = redis.GetRedisClient().Set(ctx , redisKey , string(redisValue) , 0).Err()
// 	if err != nil {
// 		t.Fatal(err.Error())
// 		return
// 	}
// 	t.Log("测试成功")
// }

func TestHotKeyLogicExpired(t *testing.T) {
	redis.Init()
	mysql.Init()

	var id int64 = 16

	var wg sync.WaitGroup

	wg.Add(10)

	for i := range 10 {
		go func(x int) {
			defer wg.Done()
			var shopInfo model.Shop
			shopInfo, err := ShopManager.QueryShopByIdWithLogicExpire(id)
			if err != nil {
				return
			}
			t.Log("goroutine: ", i, shopInfo)
		}(i)
	}
	wg.Wait()
}
