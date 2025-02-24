package handler

import (
	"hmdp/src/dto"
	"hmdp/src/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ShopTypeHandler struct {

}

var shopTypeHandler *ShopTypeHandler


// @Description: query shop type list
// @Router: /shop-type/list  [GET]
func (*ShopTypeHandler) QueryShopTypeList(c *gin.Context) {
	shopTypeList , err := service.ShopTypeManager.QueryShopTypeList()
	if err != nil {
		logrus.Error("failed to get type list")
		c.JSON(http.StatusOK , dto.Fail[string]("failed to get type list"))
		return 
	}
	c.JSON(http.StatusOK , dto.OkWithData(shopTypeList))
}
