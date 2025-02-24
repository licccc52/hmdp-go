package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ConfigRouter(r *gin.Engine) {

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	blogController := r.Group("/blog")

	{
		blogController.POST("", blogHandler.SaveBlog)
		blogController.PUT("/like/:id", blogHandler.LikeBlog)
		blogController.GET("/of/me", blogHandler.QueryMyBlog)
		blogController.GET("/hot", blogHandler.QueryHotBlog)
		blogController.GET("/:id", blogHandler.GetBlogById)
	}

	// blogCommentsController := r.Group("/blog-comments")

	// followController := r.Group("/follow")

	shopController := r.Group("/shop")

	{
		shopController.GET("/:id", shopHandler.QueryShopById)
		shopController.POST("", shopHandler.SaveShop)
		shopController.PUT("", shopHandler.UpdateShop)
		shopController.GET("/of/type", shopHandler.QueryShopByType)
		shopController.GET("/of/name", shopHandler.QueryShopByName)
	}

	shopTypeController := r.Group("/shop-type")

	{
		shopTypeController.GET("/list", shopTypeHandler.QueryShopTypeList)
	}

	uploadController := r.Group("/upload")

	{
		uploadController.POST("/blog", uploadHandler.UploadImage)
		uploadController.GET("/blog/delete", uploadHandler.DeleteBlogImg)
	}

	userController := r.Group("/user")

	{
		userController.POST("/code", userHandler.SendCode)
		userController.POST("/login", userHandler.Login)
		userController.POST("/logout", userHandler.Logout)
		userController.GET("/me", userHandler.Me)
		userController.GET("/info/:id", userHandler.Info)
	}

	voucherController := r.Group("/voucher")

	{
		voucherController.POST("", voucherHandler.AddVoucher)
		voucherController.POST("/seckill", voucherHandler.AddSecKillVoucher)
		voucherController.GET("/list/:shopId", voucherHandler.QueryVoucherOfShop)
	}

	voucherOrderController := r.Group("/voucher-order")

	{
		voucherOrderController.POST("/seckill/:id", voucherOrderHandler.SeckillVoucher)
	}
}
