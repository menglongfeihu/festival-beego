package job

import (
	"festival-fun/consts"
	"festival-fun/models"
	"festival-fun/util"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
	//	"github.com/astaxie/beego/orgs"
)

func StatJob() {
	i := 0
	c := cron.New()
	spec := "0 32 13 * * *"
	c.AddFunc(spec, func() {
		i++
		logs.Info("start", i)
//		execute()
	})
	c.Start()
	/*	select {}*/ //阻塞主线程不退出

}

/**
  获取所有参与集福用户，刷新榜单
*/
func execute() {
	//获取所有集福用户
	var day = util.GetDateInt(time.Now().Unix())
	var dayTime = util.GetDateFormat(time.Now().Unix(), "2006-01-02 15:04:01")
	logs.Info("execute user rank stat,day:%s,time:%s", day, dayTime)
	users, _ := models.ListUsersByStatus(consts.USER_STATUS_FAVORING, "-favor_num", 100)
	if users != nil && len(users) > 0 {
		var userRanks []*models.UserRank
		for i, n := 0, len(users); i < n; i++ {
			userRank := new(models.UserRank)
			userRank.Passport = users[i].Passport
			userRank.Day = day
			userRank.CreateTime = dayTime
			userRank.Rank = int64(i + 1)
			userRanks[i] = &userRank
		}
		if userRanks != nil && len(userRanks) > 0 {
			//删除库里原有该天的数据
			num, err := models.DeleteUserRank(day)
			if err != nil {
				logs.Error("DeleteUserRank error,day:%s", day, err)
			} else {
				logs.Info("DeleteUserRank num:%s", strconv.FormatInt(num, 10))
				//保存新排行到库里
				err := models.BatchSaveUserRank(userRanks)
				if err != nil {
					logs.Error("BatchSaveUser error,day:%s", day, err)
				}
			}
		}
	}
}
