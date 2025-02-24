package model

import (
	"hmdp/src/config/mysql"
	"time"

	_ "github.com/jinzhu/gorm"
)

const VOUCHER_TABLE_NAME = "tb_voucher"

type Voucher struct {
	Id		int64		`gorm:"primary;AUTO_INCREMENT;column:id" json:"id"`
	ShopId	int64		`gorm:"column:shop_id" json:"shopId"`
	Title	string		`gorm:"column:title" json:"title"`
	SubTitlte	string	`gorm:"column:sub_title" json:"subTitle"`
	Rules	string		`gorm:"column:rules" json:"rules"`
	PayValue	int64	`gorm:"column:pay_value" json:"payValue"`
	ActualValue	int64	`gorm:"column:actual_value" json:"actualValue"`
	Type		int		`gorm:"column:type" json:"type"`
	Status		int		`gorm:"column:status" json:"status"`
	Stock		int		`gorm:"-" json:"stock"`
	BeginTime	time.Time	`gorm:"-" json:"bgeinTime"`
	EndTime		time.Time	`gorm:"-" json:"endTime"`
	CreateTime	time.Time	`gorm:"column:create_time" json:"createTime"`
	UpdateTime	time.Time	`gorm:"column:update_time" json:"updateTime"`
}

func (*Voucher) TableName() string {
	return VOUCHER_TABLE_NAME
}

func (voucher *Voucher) AddVoucher() error {
	err := mysql.GetMysqlDB().Table(voucher.TableName()).Create(voucher).Error
	return err
}

func (voucher *Voucher) QueryVoucherByShop(shopId int64) ([]Voucher , error) {
	var vouchers []Voucher
	err := mysql.GetMysqlDB().Table(voucher.TableName()).Where("shop_id = ?" , shopId).Find(&vouchers).Error
	return vouchers , err
} 
