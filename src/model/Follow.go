package model

import (
	"time"
	_"github.com/jinzhu/gorm"
)

type Follow struct {
	Id  int64  `gorm:"primary;AUTO_INCREMENT;column:id" json:"id"`
	UserId int64  `gorm:"column:user_id" json:"userId"`
	FollowUserId int64  `gorm:"column:follow_user_id" json:"followUserId"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
}
