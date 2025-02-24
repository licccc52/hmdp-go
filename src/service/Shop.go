package service

import (
	"hmdp/src/model"
)

type ShopService struct {

}

var ShopManager *ShopService

func (*ShopService) QueryShopById(id int64) (model.Shop , error) {
	var shop model.Shop
	shop.Id = id
	err := shop.QueryShopById()
	return shop , err 
}

func (*ShopService) SaveShop(shop *model.Shop) error {
	err := shop.SaveShop()
	return err	
}

func (*ShopService) UpdateShop(shop *model.Shop) error {
	err := shop.UpdateShop()
	return err
}

func (*ShopService) QueryByType(typeId int , current int) ([]model.Shop , error) {
	var shopUtils model.Shop
	shops , err := shopUtils.QueryShopByType(typeId , current)
	return shops , err
}

func (*ShopService) QueryByName(name string , current int) ([]model.Shop , error) {
	var shopUtils model.Shop
	shops , err := shopUtils.QueryShopByName(name , current)
	return shops , err
}
