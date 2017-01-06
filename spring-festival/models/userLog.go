package models

import (
	"festival-fun/consts"
	"festival-fun/redis"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserLog struct {
	Id         int64 `orm:"pk;auto"`
	Passport   string
	Fpassport  string
	Day        string
	Type       int
	Ip         string
	CreateTime string
}

func UserLogTableName() string {
	return beego.AppConfig.String("mysql.prefix") + "user_log"
}

func GetUserSignLog(passport string, ttype int, day string) (UserLog, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USERLOG_PASSPORT_TYPE_DAY, passport, ttype, day)
	var userLog UserLog
	var err error
	cacheResult := redis.Get(cacheKey, &userLog)
	if cacheResult && userLog.Id > 0 {
		beego.Info("userlog info exist in cache", userLog)
		err = nil
	} else {
		o := orm.NewOrm()
		err = o.QueryTable(UserLogTableName()).Filter("Passport", passport).Filter("day", day).Filter("type", ttype).One(&userLog)
		redis.Set(cacheKey, userLog, consts.EXPIRE_ONE_DAY)
	}

	return userLog, err
}

func GetUserLogByPassportAndFPassport(passport string, fpassport string) (UserLog, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USERLOG_FAVOR_PASSPORT_FPASSPORT, passport, fpassport)
	var userLog UserLog
	var err error
	cacheResult := redis.Get(cacheKey, &userLog)
	if cacheResult && userLog.Id > 0 {
		beego.Info("userlog info exist in cache", userLog)
		err = nil
	} else {
		o := orm.NewOrm()
		err = o.QueryTable(UserLogTableName()).Filter("Passport", passport).Filter("Fpassport", fpassport).Filter("Type", consts.USER_LOG_FAVOR).One(&userLog)
		redis.Set(cacheKey, userLog, consts.EXPIRE_ONE_DAY)
	}
	return userLog, err
}

func ListUserLogByPassportAndType(passport string, kind int) ([]UserLog, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USERLOG_LIST_PASSPORT_TYPE, passport, kind)
	var userLogs []UserLog
	var err error
	cacheResult := redis.Get(cacheKey, &userLogs)
	if cacheResult && len(userLogs) > 0 {
		beego.Info("userlog list info exist in cache", userLogs)
		err = nil
	} else {
		o := orm.NewOrm()
		var field string
		if kind == consts.USER_LOG_FAVOR {
			field = "Fpassport"
		} else {
			field = "Passport"
		}
		_, err = o.QueryTable(UserLogTableName()).Filter(field, passport).Filter("Type", kind).OrderBy("-CreateTime").All(&userLogs)
		redis.Set(cacheKey, userLogs, consts.EXPIRE_ONE_DAY)
	}
	return userLogs, err
}

func CountUserLogByPassportAndType(passport string, kind int) (int64, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_USERLOG_COUNT_PASSPORT_TYPE, passport, kind)
	var count int64
	var err error
	cacheResult := redis.Get(cacheKey, &count)
	if cacheResult && count > 0 {
		beego.Info("userlog list info exist in cache", count)
		err = nil
	} else {
		o := orm.NewOrm()
		var field string
		if kind == consts.USER_LOG_FAVOR {
			field = "Fpassport"
		} else {
			field = "Passport"
		}
		count, err = o.QueryTable(UserLogTableName()).Filter(field, passport).Filter("Type", kind).Count()
		redis.Set(cacheKey, count, consts.EXPIRE_ONE_HOUR)
	}

	return count, err
}

func SaveUserLog(userLog *UserLog) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(userLog)
	if userLog.Type == consts.USER_LOG_FAVOR {
		redis.Set(fmt.Sprintf(consts.CACHE_USERLOG_FAVOR_PASSPORT_FPASSPORT, userLog.Passport, userLog.Fpassport), userLog, consts.EXPIRE_ONE_DAY)
		redis.Set(fmt.Sprintf(consts.CACHE_USERLOG_COUNT_PASSPORT_TYPE, userLog.Fpassport, userLog.Type), userLog, consts.EXPIRE_ONE_DAY)
	} else {
		redis.Set(fmt.Sprintf(consts.CACHE_USERLOG_COUNT_PASSPORT_TYPE, userLog.Passport, userLog.Type), userLog, consts.EXPIRE_ONE_DAY)
	}

	return id, err
}
