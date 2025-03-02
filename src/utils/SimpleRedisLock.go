package utils

import (
	"bytes"
	"context"
	redisConfig "github.com/redis/go-redis/v9"
	redisClient "hmdp/src/config/redis"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SimpleRedisLock struct {
	Name string
}

const (
	KEY_PREFIX = "lock:"
)

var ID_PREFIX = uuid.New().String() + "-"

func (s *SimpleRedisLock) TryLock(expireTime int64) bool {
	key := KEY_PREFIX + s.Name
	id := ID_PREFIX + strconv.FormatUint(GetGoroutinueID(), 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	flag, err := redisClient.GetRedisClient().SetNX(ctx, key, id, time.Duration(expireTime)*time.Second).Result()

	if err != nil {
		return false
	}

	return flag
}

// func (s *SimpleRedisLock) Unlock() {
// 	realID := ID_PREFIX + strconv.FormatUint(GetGoroutinueID(), 10)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	id, err := redis.GetRedisClient().Get(ctx, KEY_PREFIX+s.Name).Result()
// 	if err != nil {
// 		return
// 	}
//
// 	if realID == id {
// 		redis.GetRedisClient().Del(ctx, KEY_PREFIX+s.Name)
// 	}
// }

// get the atomic operation
func (s *SimpleRedisLock) Unlock() {
	var script = redisConfig.NewScript(`
		local id = redis.call('get' , KEYS[1])
		if(id == ARGV[1]) then
			return redis.call('del' , KEYS[1])
		end
	`)

	keys := []string{KEY_PREFIX + s.Name}
	realId := ID_PREFIX + strconv.FormatUint(GetGoroutinueID(), 10)
	var values []interface{}
	values = append(values, realId)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	script.Run(ctx, redisClient.GetRedisClient(), keys, values...)
}

func GetGoroutinueID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
