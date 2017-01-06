package models

import (
	"festival-fun/consts"
	"festival-fun/redis"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id         int    `json:"id"`
	Passport   string `json:"passport"`
	Headpic    string `json:"headpic"`
	Nickname   string `json:"nickname"`
	Userid     int64  `json:"userid"`
	Status     int    `json:"status"`
	CreateTime string `json:"createtime"`
	UpdateTime string `json:"updatetime"`
	Phone      string `json:"phone"`
	IpNum      int    `json:"ip_num"`
	FavorNum   int64  `json:"favor_num"`
}

func UserTableName() string {
	return beego.AppConfig.String("mysql.prefix") + "user"
}

func GetUserById(id int) (User, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USER_ID, id)
	var user User
	var err error
	cacheResult := redis.Get(cacheKey, &user)
	if cacheResult && user.Id > 0 {
		beego.Info("user info exist in cache", user)
		err = nil
	} else {
		o := orm.NewOrm()
		err = o.QueryTable(UserTableName()).Filter("id", id).One(&user)
		if err == nil || err == orm.ErrNoRows {
			redis.Set(cacheKey, user, consts.EXPIRE_ONE_DAY)
			return user, nil
		}
	}

	return user, err
}

func GetUserByPassport(passport string) (User, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USER_PASSPORT, passport)
	var user User
	var err error
	cacheResult := redis.Get(cacheKey, &user)
	if cacheResult && user.Id > 0 {
		beego.Info("user info exist in cache", user)
		err = nil
	} else {
		o := orm.NewOrm()
		var user User
		err := o.QueryTable(UserTableName()).Filter("passport", passport).One(&user)
		if err == nil || err == orm.ErrNoRows {
			redis.Set(cacheKey, user, consts.EXPIRE_ONE_DAY)
			return user, nil
		}
	}

	return user, err
}

/*order默认为"",limit为0不限制个数*/
func ListUsersByStatus(status int, orderCol string, limitCol int) ([]*User, error) {
	o := orm.NewOrm()
	var users []*User
	qs := o.QueryTable(UserTableName()).Filter("status", status)
	if orderCol != "" {
		qs = qs.OrderBy(orderCol)
	}
	if limitCol > 0 {
		qs = qs.Limit(limitCol)
	}
	_, err := qs.All(&users)
	if err != nil && err != orm.ErrNoRows {
		beego.Error("ListUsersByStatus err", err)
		return users, err
	}
	return users, nil
}

func SaveUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

func UpdateUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Update(user, "IpNum", "FavorNum", "Status", "UpdateTime")
	redis.Set(fmt.Sprintf(consts.CACHE_USER_PASSPORT, user.Passport), user, consts.EXPIRE_ONE_DAY)
	redis.Set(fmt.Sprintf(consts.CACHE_USER_ID, user.Id), user, consts.EXPIRE_ONE_DAY)
	return err
}
