package models

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserRank struct {
	Id         int64 `orm:"pk;auto"`
	Passport   string
	Rank       int64
	Day        string
	CreateTime string `orm:"auto_now_add;type(datetime)"`
}

func UserRankTableName() string {
	return beego.AppConfig.String("mysql.prefix") + "user_rank"
}
func GetUserRank(passport string) (int64, error) {
	var userRank UserRank
	o := orm.NewOrm()
	err := o.QueryTable(UserRankTableName()).Filter("passport", passport).OrderBy("-day").One(&userRank, "rank")
	if err != nil {
		return -1, err
	}
	return userRank.Rank, nil
}

func BatchSaveUserRank(users []UserRank) error {
	o := orm.NewOrm()
	_, err := o.InsertMulti(1, &users)
	return err
}

func DeleteUserRank(day string) (int64, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable(UserRankTableName()).Filter("day", day).Delete()
	return num, err
}

func MakeUserRanksByUsers(users []*User, day string, dayTime string) []UserRank {
	var ranks []UserRank
	if users != nil && len(users) > 0 {
		for i, n := 0, len(users); i < n; i++ {
			var userRank UserRank
			userRank.Passport = users[i].Passport
			userRank.Day = day
			userRank.CreateTime = dayTime
			userRank.Rank = int64(i + 1)
			beego.Info("passport:" + userRank.Passport + ",i:" + strconv.FormatInt(userRank.Rank, 10))
			ranks = append(ranks, userRank)
		}
	}
	return ranks
}
