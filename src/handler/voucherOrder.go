package handler

import (
	"hmdp/src/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VoucherOrderHandler struct {
}

var voucherOrderHandler *VoucherOrderHandler

func (*VoucherOrderHandler) SeckillVoucher(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Fail[string]("the function is not finished"))
}
