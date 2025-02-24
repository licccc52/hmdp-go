package handler

import (
	"hmdp/src/dto"
	"hmdp/src/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {

}

var userHandler *UserHandler

// @Description: send the phone code
// @Router: /user/code [POST]
func (*UserHandler) SendCode(c *gin.Context) {
	// TODO 发送段新验证码并且保存验证码
	c.JSON(http.StatusOK , dto.Fail[string]("the function not finished"))
}

// @Description: user login in
// @Router: /user/login  [POST]
func (*UserHandler) Login(c *gin.Context) {
	// TODO 实现登录功能
	c.JSON(http.StatusOK , dto.Fail[string]("this function is not finished"))	
}

// @Description: user layout 
// @Router: /user/logout [POST]
func (*UserHandler) Logout(c *gin.Context) {
	// TODO 实现注销功能	
	c.JSON(http.StatusOK , dto.Fail[string]("this function is not finished"))	
}

// @Description: get the info of me
// @Router: /user/me [GET]
func (*UserHandler) Me(c *gin.Context) {
	// TODO 获取当前登录的用户
	c.JSON(http.StatusOK , dto.Fail[string]("this function is not finished"))	
}

// @Description: get the info of user by user Id
// @Router /user/info/:id [GET]
func (*UserHandler) Info(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		logrus.Error("id str is empty!")
		c.JSON(http.StatusOK , dto.Fail[string]("id str is empty!"))
		return 
	}
	id,err := strconv.ParseInt(idStr , 10 , 64)
	if err != nil {
		logrus.Error("parse int failed!")
		c.JSON(http.StatusOK , dto.Fail[string]("id parse failed!"))
		return 
	}
	userInfo , err := service.UserInfoManager.GetUserInfoById(id)
	if err != nil {
		logrus.Error("get user info failed!")
		c.JSON(http.StatusOK , dto.Fail[string]("get user info failed!"))
		return
	}
	userInfo.CreateTime = time.Time{}
	userInfo.UpdateTime = time.Time{}
	c.JSON(http.StatusOK , dto.OkWithData(userInfo))
}
