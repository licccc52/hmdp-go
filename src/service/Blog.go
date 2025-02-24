package service

import (
	"hmdp/src/dto"
	"hmdp/src/model"

	"github.com/sirupsen/logrus"
)

type BlogService struct{

}

var BlogManager *BlogService

func (*BlogService)SaveBlog(blog *model.Blog) (result dto.Result[int64],err error) {
	id , err := blog.SaveBlog()		
	if err != nil {
		logrus.Error("[Blog Service] failed to insert data!")
		result = dto.Fail[int64]("failed to insert data into database!")
		return 
	}
	result = dto.OkWithData[int64](id)
	return 
}

func (*BlogService) LikeBlog(id int64) (err error) {
	var blog model.Blog
	blog.Id = id
	err = blog.IncreseLike()
	return 
}

func (*BlogService) QueryMyBlog(id int64 , current int) ([]model.Blog , error) {
	var blog model.Blog
	blog.UserId = id
	blogs , err := blog.QueryBlogs(current)
	return blogs , err
}

func (*BlogService) QueryHotBlogs(current int) ([]model.Blog , error) {
	var blogUtils model.Blog
	blogs , err := blogUtils.QueryHots(current)
	if err != nil {
		return nil , err
	}
	for i := range blogs {
		id := blogs[i].UserId
		user,err := UserManager.GetUserById(id)
		if err != nil {
			return blogs , err
		}
		blogs[i].Icon = user.Icon
		blogs[i].Name = user.NickName
	}

	return blogs , nil
}

func (*BlogService) GetBlogById(id int64) (model.Blog ,  error) {
	var blog model.Blog	
	err := blog.GetBlogById(id)
	return blog , err
}
