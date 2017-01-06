package controllers

import (
	"festival-fun/consts"
	"festival-fun/models"
	"festival-fun/util"
	"os"
	"strings"

	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	BaseController
}

func (this *UserController) Get() {
	//	id := this.GetString("id")
	callback := this.GetString("callback", "")
	env := os.Getenv("GOMODE")
	beego.Info("env:", env)

	//	beego.Info("id:" + id + ", callback:" + callback)
	//	if id == "" {
	//		util.ResponseLackParam(this.Ctx.Output, callback)
	//	}
	//	idInt, err := strconv.ParseInt(id, 10, 64)
	//	if err != nil {
	//		beego.Error("parse int error,id:" + id)
	//		util.ResponseSystemErr(this.Ctx.Output, callback)
	//	}
	//	ob, err := models.GetUserById(int(idInt))
	//	if err == nil {
	//		if ob.Id == 0 {
	//			util.ResponseUserNotExist(this.Ctx.Output, callback)
	//		} else {
	//			util.ResponseObj(this.Ctx.Output, ob, callback)
	//		}
	//	}

	//	users, err := models.ListUsersByStatus(consts.USER_STATUS_FAVORING, "-favor_num", 100)
	//	users, err := models.ListUsersByStatus(-1, "-favor_num", 100)
	//	if err != nil {
	//		util.ResponseSystemErr(this.Ctx.Output, callback)
	//	}
	util.ResponseObj(this.Ctx.Output, env, callback)
	//	util.ResponseSystemErr(this.Ctx.Output, callback)
}

// 获取用户信息
func (this *UserController) GetUserStatus() {
	if "" == this.Passport || "" == this.Token {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
	}
	accountInfo, _ := util.GetAccountByPassport(this.Passport)
	if accountInfo.Uid == 0 {
		util.ResponseAccountNotExist(this.Ctx.Output, this.Callback)
		return
	}
	isLogin := util.CheckLogin(this.Passport, this.Token)
	if isLogin == false {
		util.ResponseUserNoLogin(this.Ctx.Output, this.Callback)
		return
	}
	user, err := models.GetUserByPassport(this.Passport)
	if err != nil {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}
	//获取日期

	var ipToday = 0
	userLog, err := models.GetUserSignLog(this.Passport, consts.USER_LOG_IP, util.GetDateInt(time.Now().Unix()))
	if err != nil {
		beego.Error("GetUserSignLog error", err)
		if err == orm.ErrNoRows {
			beego.Error("No GetUserSignLog find")
		}
	}
	if userLog.Id > 0 {
		ipToday = 1
	}

	rankInt, err := models.GetUserRank(this.Passport)
	if err != nil {
		beego.Error("GetUserRank error", err)
		rankInt = -1
	}
	data := make(map[string]interface{})
	data["status"] = user.Status
	data["ip_num"] = user.IpNum
	data["favor_num"] = user.FavorNum
	data["nickname"] = accountInfo.NickName
	data["headpic"] = accountInfo.SmallImg

	data["rank"] = rankInt
	data["ip_today"] = ipToday
	util.ResponseObj(this.Ctx.Output, data, this.Callback)
}

// 判断用户是否绑定手机号
func (this *UserController) IsBindPhone() {
	if "" == this.Passport {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	accountInfo, _ := util.GetAccountByPassport(this.Passport)
	if accountInfo.Uid == 0 {
		util.ResponseAccountNotExist(this.Ctx.Output, this.Callback)
		return
	}

	result, _ := util.IsBindPhone(this.Passport)
	data := make(map[string]interface{})

	if "" != result.Mobile {
		data["is_bind"] = 1
	} else {
		data["is_bind"] = 0
	}
	util.ResponseObj(this.Ctx.Output, data, this.Callback)
	return
}

// 订阅专辑
func (this *UserController) SubscribeVideoInfo() {
	vname := strings.Trim(this.GetString("vname"), " ")
	if "" == this.Passport || "" == vname || "" == this.Token {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	accountInfo, _ := util.GetAccountByPassport(this.Passport)
	if accountInfo.Uid == 0 {
		util.ResponseAccountNotExist(this.Ctx.Output, this.Callback)
		return
	}
	isLogin := util.CheckLogin(this.Passport, this.Token)
	if isLogin == false {
		util.ResponseUserNoLogin(this.Ctx.Output, this.Callback)
		return
	}
	userInfo, _ := models.GetUserByPassport(this.Passport)
	if userInfo.Id == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}

	videoSubinfo, _ := models.GetVideoSubinfoByPassportAndVname(this.Passport, vname)
	if videoSubinfo.Id > 0 {
		util.ResponseUserHasSubVideo(this.Ctx.Output, this.Callback)
		return
	} else {
		video := new(models.VideoSubinfo)
		video.Passport = this.Passport
		video.Vname = vname
		video.CreateTime = util.GetDateMH(time.Now().Unix())
		models.SaveVideoSubInfo(video)
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
	return
}

//充值搜狐视频黄金会员
func (this *UserController) RechargeVIPMember() {
	vipType, _ := this.GetInt("type", 0)

	if "" == this.Passport || "" == this.Token {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	accountInfo, _ := util.GetAccountByPassport(this.Passport)
	if accountInfo.Uid == 0 {
		util.ResponseAccountNotExist(this.Ctx.Output, this.Callback)
		return
	}
	isLogin := util.CheckLogin(this.Passport, this.Token)
	if isLogin == false {
		util.ResponseUserNoLogin(this.Ctx.Output, this.Callback)
		return
	}
	userInfo, _ := models.GetUserByPassport(this.Passport)
	if userInfo.Id == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}

	count64, _ := models.CountUserLogByPassportAndType(this.Passport, consts.USER_LOG_FAVOR)
	//count, _ := strconv.Atoi(strconv.FormatInt(count64, 10))
	if count64 < consts.FAVOR_VIP_3 {
		util.ResponseFavorNotEnough(this.Ctx.Output, this.Callback)
		return
	}
	var proid int64
	var activityId string
	switch vipType {
	case consts.VIP_TYPE_14:
		beego.Info("recharge 14 vip")
		proid = 2
		activityId = beego.AppConfig.String("vip_activiId_14")
	case consts.VIP_TYPE_7, consts.VIP_TYPE_3:
		beego.Info("recharge 7 vip")
		proid = 1
		activityId = beego.AppConfig.String("vip_activiId_7")
	default:
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	userAward, _ := models.GetUseAwardByPassportAndProId(this.Passport, proid)
	if userAward.Id > 0 && userAward.Status == 0 {
		util.ResponseVipHasObtain(this.Ctx.Output, this.Callback)
		return
	}
	//充值会员
	recharge, msg := util.VipVideoRecharge(this.Passport, activityId, beego.AppConfig.String("vip_channelId"))
	if recharge == false {
		beego.Error("recharge failed：" + msg)
		util.ResponseSystemErr(this.Ctx.Output, this.Callback)
		return
	} else {
		userAward.Status = 0
		userAward.UpdateTime = util.GetDateMH(time.Now().Unix())
		userAward.ReceiveTime = userAward.UpdateTime
		models.UpdateUseAward(&userAward)
		util.ResponseOK(this.Ctx.Output, this.Callback)
		return
	}

}
