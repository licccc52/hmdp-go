package service

import (
	"hmdp/src/model"
)

type UserService struct {
	
}

var UserManager *UserService

func (*UserService) GetUserById(id int64) (model.User , error) {
	var userUtils model.User
	user , err := userUtils.GetUserById(id)
	return user , err	
}
