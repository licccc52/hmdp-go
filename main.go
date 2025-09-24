package main

import (
	"hmdp/src/config"
	"hmdp/src/handler"
	"hmdp/src/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.Init()
	handler.ConfigRouter(r)
	service.InitOrderHandler()

	r.Run(":8081")
}

func test() {
	//
}
