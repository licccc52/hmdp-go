package utils

import (
	"sync"
	"hmdp/src/dto"
)

type UserHolder struct {
	userInfo dto.UserDTO
	lock	 sync.Mutex	
}

var UserInfo *UserHolder

func(holder *UserHolder) GetUserInfo() dto.UserDTO {
	holder.lock.Lock()
	defer holder.lock.Unlock()
	return holder.userInfo
}

func (holder *UserHolder) SetUserInfo(u dto.UserDTO) {
	holder.lock.Lock()
	holder.userInfo = u
	holder.lock.Unlock()
}
