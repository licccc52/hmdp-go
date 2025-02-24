package service

import (
	"hmdp/src/model"
)
type VoucherService struct {

}

var VoucherManager *VoucherService

func (*VoucherService) AddVoucher(voucher *model.Voucher) (error) {
	err := voucher.AddVoucher()
	return err
}

func (*VoucherService) QueryVoucherOfShop(shopId int64) ([]model.Voucher , error) {
	var vocherUtils model.Voucher
	return vocherUtils.QueryVoucherByShop(shopId)
}

func (*VoucherService) AddSeckillVoucher(voucher *model.Voucher) error {
	err := voucher.AddVoucher()
	if err != nil {
		return err
	}
	var seckillVoucher model.SecKillVoucher
	seckillVoucher.VoucherId = voucher.Id
	seckillVoucher.Stock = voucher.Stock
	seckillVoucher.BeginTime = voucher.BeginTime
	seckillVoucher.EndTime = voucher.EndTime
	err = seckillVoucher.AddSeckillVoucher()
	return err 
}
