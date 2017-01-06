package controllers

import (
	"festival-fun/consts"
	"festival-fun/models"
	"festival-fun/util"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type ProductController struct {
	BaseController
}

// 获取实物奖品列表
func (this *ProductController) GetProductList() {
	data, err := models.ListProductByType(consts.PRODUCT_TYPE_ENTITY)
	if nil != err {
		beego.Error("Get Product list failed:", err)
		util.ResponseSystemErr(this.Ctx.Output, this.Callback)
		return
	}
	util.ResponseObj(this.Ctx.Output, data, this.Callback)
}

// 获取我的奖品列表
func (this *ProductController) GetAwardsList() {
	if "" == this.Passport || "" == this.Token {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	accountInfo, _ := util.GetAccountByPassport(this.Passport)
	if accountInfo.Uid == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}
	isLogin := util.CheckLogin(this.Passport, this.Token)
	if isLogin == false {
		util.ResponseUserNoLogin(this.Ctx.Output, this.Callback)
		return
	}
	user, _ := models.GetUserByPassport(this.Passport)
	if user.Id == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}
	rankInt, _ := models.GetUserRank(this.Passport)
	data := make(map[string]interface{})
	data["status"] = user.Status
	data["ip_num"] = user.IpNum
	data["favor_num"] = user.FavorNum
	data["nickname"] = user.Nickname
	data["headpic"] = user.Headpic
	data["rank"] = rankInt
	awards, _ := models.ListUserAwardByPassport(this.Passport)
	data["awards"] = awards
	util.ResponseObj(this.Ctx.Output, data, this.Callback)
	return
}

// 填写领取实物奖品地址信息
func (this *ProductController) AppendAwardsAddress() {
	if "" == this.Passport || "" == this.Token {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	accountInfo, _ := util.GetAccountByPassport(this.Passport)
	if accountInfo.Uid == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}
	isLogin := util.CheckLogin(this.Passport, this.Token)
	if isLogin == false {
		util.ResponseUserNoLogin(this.Ctx.Output, this.Callback)
		return
	}
	user, _ := models.GetUserByPassport(this.Passport)
	if user.Id == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}
	pid, _ := this.GetInt64("pid", 0)
	name := strings.Trim(this.GetString("name"), " ")
	addr := strings.Trim(this.GetString("addr"), " ")
	phone := strings.Trim(this.GetString("phone"), " ")
	qq := strings.Trim(this.GetString("qq"), " ")
	email := strings.Trim(this.GetString("email"), " ")

	userAward, _ := models.GetUseAwardByPassportAndProId(this.Passport, pid)

	if userAward.Id > 0 {
		if userAward.Status == consts.AWARD_STATUS_UNDRAW {
			userAward.Address = addr
			userAward.UserName = name
			userAward.Phone = phone
			userAward.Email = email
			userAward.QQ = qq
			userAward.Status = consts.AWARD_STATUS_DRAW
			userAward.UpdateTime = util.GetDateMH(time.Now().Unix())
			userAward.ReceiveTime = userAward.UpdateTime
			models.UpdateUseAward(&userAward)
		} else {
			util.ResponseUserHasAward(this.Ctx.Output, this.Callback)
			return
		}

	} else {
		util.ResponseSystemErr(this.Ctx.Output, this.Callback)
		return
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
	return
}
