package job

import (
	"festival-fun/consts"
	"festival-fun/models"
	"festival-fun/util"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

func StatJob() {
	job1 := toolbox.NewTask("user_stat", "0 01 00 */1 * *", statFavor)
	toolbox.AddTask("user_stat", job1)
}

/**
  获取所有参与集福用户，刷新榜单
*/
func statFavor() error {
	//获取所有集福用户
	var day = util.GetDateInt(time.Now().Unix())
	var dayTime = util.GetDateFormat(time.Now().Unix(), "2006-01-02 15:04:01")
	beego.Info("execute user rank stat day:", day, "time:", dayTime)
	users, err := models.ListUsersByStatus(consts.USER_STATUS_FAVORING, "-FavorNum", 100)
	if err != nil {
		beego.Error("ListUsersByStatus error", err)
		return err
	}
	if users != nil && len(users) > 0 {
		userRanks := models.MakeUserRanksByUsers(users, day, dayTime)
		if userRanks != nil && len(userRanks) > 0 {
			//删除库里原有该天的数据
			num, err := models.DeleteUserRank(day)
			if err != nil {
				beego.Error("DeleteUserRank day:", day, ",error:", err)
				return err
			} else {
				beego.Info("DeleteUserRank num:", strconv.FormatInt(num, 10))
				//保存新排行到库里
				err := models.BatchSaveUserRank(userRanks)
				if err != nil {
					beego.Error("BatchSaveUser day:", day, ",error:", err)
					return err
				}
			}
		}
	}
	return nil
}
