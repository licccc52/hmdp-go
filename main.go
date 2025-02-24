package main

import (
	"hmdp/src/config"
	"hmdp/src/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.Init()
	handler.ConfigRouter(r)

	r.Run(":8081")
}
