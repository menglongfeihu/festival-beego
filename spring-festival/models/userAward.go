package models

import (
	"festival-fun/consts"
	"festival-fun/redis"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserAward struct {
	Id          int64  `orm:"pk;auto"`
	Passport    string `orm:"index"`
	UserName    string `orm:"column(username)"`
	ProductId   int64
	Type        int
	Status      int
	Address     string `orm:"null;size(512)"`
	Phone       string `orm:"null"`
	QQ          string `orm:"null;column(qq)"`
	Email       string `orm:"null"`
	ReceiveTime string
	CreateTime  string
	UpdateTime  string
}

func UserAwardTableName() string {
	return beego.AppConfig.String("mysql.prefix") + "user_award"
}

func GetUseAwardByPassportAndProId(passport string, proid int64) (UserAward, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USERAWARD_PASSPORT_PROID, passport, proid)
	var userAward UserAward
	var err error
	cacheResult := redis.Get(cacheKey, &userAward)
	if cacheResult && userAward.Id > 0 {
		beego.Info("userAward info exist in cache", userAward)
		err = nil
	} else {
		o := orm.NewOrm()
		err = o.QueryTable(UserAwardTableName()).Filter("Passport", passport).Filter("ProductId", proid).One(&userAward)
		redis.Set(cacheKey, userAward, consts.EXPIRE_ONE_DAY)
	}
	return userAward, err
}
func ListUserAwardByPassport(passport string) ([]map[string]interface{}, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USERAWARD_LIST_PASSPORT, passport)
	data := make([]map[string]interface{}, 0)
	var err error
	cacheResult := redis.Get(cacheKey, &data)
	if cacheResult && len(data) > 0 {
		beego.Info("userAward list info exist in cache", data)
		err = nil
	} else {
		o := orm.NewOrm()
		var userAwardList []UserAward
		_, err = o.QueryTable(UserAwardTableName()).Filter("Passport", passport).OrderBy("Id").All(&userAwardList)
		// 奖品数量
		count := len(userAwardList)
		data = make([]map[string]interface{}, count)
		if count > 0 {
			index := 0
			for _, userAward := range userAwardList {
				productInfo, _ := GetProductById(userAward.ProductId)
				award := map[string]interface{}{
					"short_name":  productInfo.ShortName,
					"name":        productInfo.Name,
					"detail":      productInfo.Detail,
					"product_pic": productInfo.ProductPic,
					"status":      userAward.Status,
					"id":          userAward.ProductId,
				}
				data[index] = award
				index++
			}
		}
		redis.Set(cacheKey, data, consts.EXPIRE_ONE_DAY)
	}
	data = removeDuplicatesAndEmpty(data)
	return data, err
}

func removeDuplicatesAndEmpty(source []map[string]interface{}) (ret []map[string]interface{}) {
	a_len := len(source)
	for i := 0; i < a_len; i++ {
		if len(source[i]) == 0 {
			continue
		}
		ret = append(ret, source[i])
	}
	return ret
}

func SaveUserAward(info *UserAward) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(info)
	redis.Remove(fmt.Sprintf(consts.CACHE_USERAWARD_LIST_PASSPORT, info.Passport))
	return id, err
}

func UpdateUseAward(info *UserAward) error {
	o := orm.NewOrm()
	_, err := o.Update(info, "UserName", "Type", "Status", "Address", "Phone", "Email", "QQ", "UpdateTime", "ReceiveTime")
	redis.Remove(fmt.Sprintf(consts.CACHE_USERAWARD_LIST_PASSPORT, info.Passport))
	redis.Set(fmt.Sprintf(consts.CACHE_USERAWARD_PASSPORT_PROID, info.Passport, info.ProductId), info, consts.EXPIRE_ONE_DAY)
	return err
}
