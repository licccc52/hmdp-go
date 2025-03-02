package model

import (
	"hmdp/src/config/mysql"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
)

const (
	NOTPAYED = 1   // 未支付
	PAYED    = 2   // 已支付
	USED     = 3   // 已核销
	CANCELED = 4   // 已取消
	RETURN	 = 5   // 退款中
	RETURNED = 6   // 已退款
)


const (
	EXTRAPAY = 1    // 余额支付
	ALIPAY   = 2    // 支付宝支付
	WEIXINPAY = 3   // 微信支付
)

type VoucherOrder struct {
	Id		int64	`gorm:"primary;column:id" json:"id"`
	UserId	int64	`gorm:"column:user_id" json:"userId"`
	VoucherId	int64	`gorm:"column:voucher_id" json:"voucherId"`
	PayType		int		`gorm:"column:pay_type" json:"payType"`
	Status		int		`gorm:"column:status" json:"status"`
	CreateTime	time.Time	`gorm:"column:create_time" json:"create_time"`
	PayTime		time.Time	`gorm:"column:pay_time" json:"payTime"`
	UseTime		time.Time	`gorm:"column:use_time" json:"useTime"`
	RefundTime	time.Time	`gorm:"column:refund_time" json:"refundTime"`
	UpdateTime	time.Time	`gorm:"column:update_time" json:"updateTime"`
}

func (*VoucherOrder) TableName() string {
	return "tb_voucher_order"
}

// func (vo *VoucherOrder) QueryVocherOrderByUserId(userId int64) error {
// 	err := mysql.GetMysqlDB().Table(vo.TableName()).Where("user_id = ?" , userId).First(vo).Error
// 	return err
// }

func (vo *VoucherOrder) CreateVoucherOrder(tx *gorm.DB) error {
	err := tx.Table(vo.TableName()).Create(vo).Error
	return err
}

func (vo *VoucherOrder) QueryVoucherOrderByUserId(userId int64) (int64 , error) {
	var count int64
	err := mysql.GetMysqlDB().Table(vo.TableName()).Where("user_id = ?" , userId).Count(&count).Error
	return count , err
}
