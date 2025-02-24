package service

import "hmdp/src/model"


type ShopTypeService struct {
	
}

var ShopTypeManager *ShopTypeService

func (*ShopTypeService) QueryShopTypeList() ([]model.ShopType , error) {
	var shopTypeUtils model.ShopType
	shopTypeList , err := shopTypeUtils.QueryTypeList()
	return shopTypeList , err
}
