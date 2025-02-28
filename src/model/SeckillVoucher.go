package model

import (
	"hmdp/src/config/mysql"
	"time"

	_ "github.com/jinzhu/gorm"
)

const SECKILL_VOUCHER_NAME = "tb_seckill_voucher"

type SecKillVoucher struct {
	VoucherId  int64     `gorm:"primary;AUTO_INCREMENT;column:voucher_id" json:"voucherId"`
	Stock      int       `gorm:"column:stock" json:"stock"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	BeginTime  time.Time `gorm:"column:begin_time" json:"beginTime"`
	EndTime    time.Time `gorm:"column:end_time" json:"endTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (*SecKillVoucher) TableName() string {
	return SECKILL_VOUCHER_NAME
}

func (sec *SecKillVoucher) AddSeckillVoucher() error {
	return mysql.GetMysqlDB().Table(sec.TableName()).Create(sec).Error
}
