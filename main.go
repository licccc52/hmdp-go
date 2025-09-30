package main

import (
	"github.com/gin-gonic/gin"
	"hmdp/src/config"
	"hmdp/src/handler"
	"hmdp/src/service"
)

func process1() {
	r := gin.Default()

	config.Init()
	handler.ConfigRouter(r)
	service.InitOrderHandler()

	r.Run(":8081")
}

func process2() {
	r := gin.Default()

	config.Init()
	handler.ConfigRouter(r)
	service.InitOrderHandler()

	r.Run(":8082")
}

func main() {
	// 模拟集群
	go process2()
	process1()
}
