package config

import (
	"hmdp/src/config/mysql"
	"hmdp/src/config/redis"
)

func Init() {
	mysql.Init()
	redis.Init()
}
