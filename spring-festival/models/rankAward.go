package models

import (
	"festival-fun/consts"
	"festival-fun/redis"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RankAwardMap struct {
	Id         int64 `orm:"pk;auto"`
	Rank       string
	ProductId  string
	CreateTime string
}

func RankAwardTableName() string {
	return beego.AppConfig.String("mysql.prefix") + "award_map"
}

func GetRankAwarkMapByRank(rank string) (RankAwardMap, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_RANKAWARDMAP_RANK, rank)
	var rankAwardMap RankAwardMap
	var err error
	cacheResult := redis.Get(cacheKey, &rankAwardMap)
	if cacheResult && rankAwardMap.Id > 0 {
		beego.Info("rankAwardMap info exist in cache", rankAwardMap)
		err = nil
	} else {
		o := orm.NewOrm()
		rankAwardMap = RankAwardMap{Rank: rank}
		err = o.Read(&rankAwardMap)
		redis.Set(cacheKey, rankAwardMap, consts.EXPIRE_ONE_DAY)
	}
	return rankAwardMap, err
}
