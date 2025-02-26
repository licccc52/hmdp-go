package handler

import (
	"hmdp/src/dto"
	"hmdp/src/middleware"
	"hmdp/src/model"
	"hmdp/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BlogHandler struct {
}

var blogHandler *BlogHandler

// @Description: save the blog
// @Router:  /blog [POST]
func (*BlogHandler) SaveBlog(c *gin.Context) {
	var blog model.Blog
	err := c.ShouldBindJSON(&blog)
	if err != nil {
		logrus.Error("[Blog handler] bind json failed!")
		c.JSON(http.StatusOK, dto.Fail[string]("insert failed!"))
		return
	}
	var result dto.Result[int64]
	result, err = service.BlogManager.SaveBlog(&blog)
	if err != nil {
		logrus.Error("[Blog handler] insert data into database failed!")
		c.JSON(http.StatusOK, dto.Fail[string]("insert failed!"))
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Description: modify the number of linked
// @Router:  /blog/like/:id  [PUT]
func (*BlogHandler) LikeBlog(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		logrus.Error("[Blog Handler] Give a empty string")
		c.JSON(http.StatusOK, dto.Fail[string]("like blog failed!"))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail[string]("type transform failed!"))
		return
	}
	err = service.BlogManager.LikeBlog(id)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail[string]("like failed!"))
		return
	}
	c.JSON(http.StatusOK, dto.Ok[string]())
}

// @Description: query my blog
// @Router: /blog/of/me [GET]
func (*BlogHandler) QueryMyBlog(c *gin.Context) {
	var current string
	current = c.Query("current")

	if current == "" {
		current = "1"
	}

	currentPage, err := strconv.Atoi(current)
	if err != nil {
		logrus.Error("type transform failed!")
		c.JSON(http.StatusOK, dto.Fail[string]("type transform failed!"))
		return
	}

	user, err := middleware.GetUserInfo(c)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail[string]("get user info failed!"))
		return
	}

	blogs, err := service.BlogManager.QueryMyBlog(user.Id, currentPage)
	if err != nil {
		logrus.Error("page query failed!")
		c.JSON(http.StatusOK, dto.Fail[string]("page query failed!"))
		return
	}
	c.JSON(http.StatusOK, dto.OkWithData[[]model.Blog](blogs))
}

// @Description: query the hot blog
// @Router: /blog/hot [GET]
func (*BlogHandler) QueryHotBlog(c *gin.Context) {
	var currentStr = "1"
	currentStr = c.Query("current")
	if currentStr == "" {
		currentStr = "1"
	}
	current, err := strconv.Atoi(currentStr)
	if err != nil {
		logrus.Error("transform type failed!")
		c.JSON(http.StatusOK, dto.Fail[string]("transform type failed!"))
		return
	}
	blogs, err := service.BlogManager.QueryHotBlogs(current)
	if err != nil {
		logrus.Error("query hot blogs failed!")
		c.JSON(http.StatusOK, dto.Fail[string]("query hot blogs failed!"))
		return
	}
	c.JSON(http.StatusOK, dto.OkWithData[[]model.Blog](blogs))
}

// @Description: Get Blog by id
// @Router: /blog/:id  [GET]
func (*BlogHandler) GetBlogById(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		logrus.Error("id str is empty")
		c.JSON(http.StatusOK, dto.Fail[string]("id str is empty"))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail[string]("type transform is failed!"))
		return
	}
	blog, err := service.BlogManager.GetBlogById(id)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusOK, dto.Fail[string]("get blog by id failed!"))
		return
	}
	c.JSON(http.StatusOK, dto.OkWithData(blog))
}
