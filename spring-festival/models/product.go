package models

import (
	"festival-fun/consts"
	"festival-fun/redis"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Product struct {
	Id         int64 `orm:"pk;auto"`
	ShortName  string
	Name       string
	Detail     string `orm:"size(1024)"`
	ProductPic string `orm:"size(512)"`
	Num        int
	CreateTime string
	UpdateTime string
	ExpireTime string
	Type       int
	VipDays    int
}

func ProductTableName() string {
	return beego.AppConfig.String("mysql.prefix") + "product"
}

func GetProductById(pid int64) (Product, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_PRODUCT_ID, pid)
	var product Product
	var err error
	cacheResult := redis.Get(cacheKey, &product)
	if cacheResult && product.Id > 0 {
		beego.Info("product info exist in cache", product)
		err = nil
	} else {
		o := orm.NewOrm()
		product = Product{Id: pid}
		err = o.Read(&product)
		redis.Set(cacheKey, product, consts.EXPIRE_ONE_DAY)
	}

	return product, err
}

func UpdateProduct(product *Product) error {
	o := orm.NewOrm()
	_, err := o.Update(product, "ShortName", "Name", "ProductPic", "Num", "Detail", "UpdateTime", "ExpireTime", "Type", "VipDays")
	redis.Set(fmt.Sprintf(consts.CACHE_PRODUCT_ID, product.Id), product, consts.EXPIRE_ONE_DAY)
	redis.Remove(fmt.Sprintf(consts.CACHE_PRODUCT_LIST, ""))
	return err
}

func ListProductByType(kind int) ([]Product, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_PRODUCT_LIST, "")
	var products []Product
	var err error
	cacheResult := redis.Get(cacheKey, &products)
	if cacheResult && len(products) > 0 {
		beego.Info("product list info exist in cache", products)
		err = nil
	} else {
		o := orm.NewOrm()
		qs := o.QueryTable(ProductTableName())
		_, err = qs.Filter("Type", kind).All(&products)
		redis.Set(cacheKey, products, consts.EXPIRE_ONE_DAY)
	}

	return products, err
}

func SaveProduct(product *Product) error {
	o := orm.NewOrm()
	_, err := o.Insert(product)
	redis.Set(fmt.Sprintf(consts.CACHE_PRODUCT_ID, product.Id), product, consts.EXPIRE_ONE_DAY)
	return err
}
