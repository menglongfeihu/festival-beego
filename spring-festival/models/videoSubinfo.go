package models

import (
	"festival-fun/consts"
	"festival-fun/redis"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type VideoSubinfo struct {
	Id         int64 `orm:"pk;auto"`
	Passport   string
	Vname      string `orm:"column(vname)"`
	CreateTime string
}

func VideoSubinfoTableName() string {
	return beego.AppConfig.String("mysql.prefix") + "video_subinfo"
}

func SaveVideoSubInfo(info *VideoSubinfo) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(info)
	redis.Set(fmt.Sprintf(consts.CACHE_VIDEOINFO_LIST_PASSPORT, info.Passport), info, consts.EXPIRE_ONE_DAY)
	return id, err
}

func GetVideoSubinfoByPassportAndVname(passport string, vname string) (VideoSubinfo, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_VIDEOINFO_PASSPORT_VNAME, passport, vname)
	var videoSubinfo VideoSubinfo
	var err error
	cacheResult := redis.Get(cacheKey, &videoSubinfo)
	if cacheResult && videoSubinfo.Id > 0 {
		beego.Info("subvideo info exist in cache", videoSubinfo)
		err = nil
	} else {
		o := orm.NewOrm()
		err = o.QueryTable(VideoSubinfoTableName()).Filter("Passport", passport).Filter("Vname", vname).One(&videoSubinfo)
		redis.Set(cacheKey, videoSubinfo, consts.EXPIRE_ONE_WEEK)
	}

	return videoSubinfo, err
}

func ListVideoSubinfoByPassport(passport string) ([]VideoSubinfo, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_VIDEOINFO_LIST_PASSPORT, passport)
	var videoSubinfos []VideoSubinfo
	var err error
	cacheResult := redis.Get(cacheKey, &videoSubinfos)
	if cacheResult && len(videoSubinfos) > 0 {
		beego.Info("subvideo info list exist in cache", videoSubinfos)
		err = nil
	} else {
		o := orm.NewOrm()
		_, err = o.QueryTable(VideoSubinfoTableName()).Filter("Passport", passport).All(&videoSubinfos)
		redis.Set(cacheKey, videoSubinfos, consts.EXPIRE_ONE_WEEK)
	}

	return videoSubinfos, err
}
