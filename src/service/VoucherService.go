package service

import (
	"context"
	"hmdp/src/config/redis"
	"hmdp/src/model"
	"hmdp/src/utils"
	"strconv"
)

type VoucherService struct {
}

var VoucherManager *VoucherService

func (*VoucherService) AddVoucher(voucher *model.Voucher) error {
	err := voucher.AddVoucher()
	return err
}

func (*VoucherService) QueryVoucherOfShop(shopId int64) ([]model.Voucher, error) {
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
	seckillVoucher.CreateTime = voucher.CreateTime
	seckillVoucher.UpdateTime = voucher.UpdateTime

	err = seckillVoucher.AddSeckillVoucher()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = redis.GetRedisClient().Set(ctx, utils.SECKILL_STOCK_KEY+strconv.FormatInt(voucher.Id, 10), voucher.Stock, 0).Err()

	return err
}
