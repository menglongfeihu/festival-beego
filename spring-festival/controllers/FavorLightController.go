package controllers

import (
	"festival-fun/consts"
	"festival-fun/models"
	"festival-fun/redis"
	"festival-fun/util"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type FavorLightController struct {
	BaseController
}

//开启集福
func (this *FavorLightController) OpenFavor() {
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
	user, _ := models.GetUserByPassport(this.Passport)
	if user.Id == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}
	// 判断用户是否已经点亮6个ip
	if user.IpNum < consts.COUNT_IP {
		util.ResponseLightIPNotEnough(this.Ctx.Output, this.Callback)
		return
	}
	if user.Status == consts.USER_STATUS_IPING {
		user.Status = consts.USER_STATUS_FAVORING
		user.UpdateTime = util.GetDateMH(time.Now().Unix())
		models.UpdateUser(&user)
	} else if user.Status == consts.USER_STATUS_FAVORING {
		util.ResponseUserHasFavor(this.Ctx.Output, this.Callback)
		return
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
	return
}

// 获取集福好友头像列表
func (this *FavorLightController) GetFrendsHeadPicList() {
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

	userLogs, _ := models.ListUserLogByPassportAndType(this.Passport, consts.USER_LOG_FAVOR)
	headPics := make([]map[string]interface{}, len(userLogs))
	count := 1
	if len(userLogs) > 0 {
		for _, log := range userLogs {
			account, err := util.GetAccountByPassport(log.Passport)
			if err == nil {
				headPics[count-1] = map[string]interface{}{"nickname": account.NickName, "headpic": account.SmallImg}
			}
			count += 1
			if count > 10 {
				break
			}
		}
	} else {
		headPics = make([]map[string]interface{}, 0)
	}
	util.ResponseObj(this.Ctx.Output, headPics, this.Callback)
	return
}

// 向好友集福
func (this *FavorLightController) FavorFriend() {
	fpassport := strings.Trim(this.GetString("fpassport"), " ")
	if "" == this.Passport && "" == fpassport {
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
	fuser, _ := models.GetUserByPassport(fpassport)
	if fuser.Id == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
		return
	}

	//判断好友是否已经进入集福状态
	if fuser.Status != consts.USER_STATUS_FAVORING {
		util.ResponseUserNotFavor(this.Ctx.Output, this.Callback)
		return
	}

	//判断用户是否可以向好友集福，防刷
	allow := this.AllowFavor(this.Passport, fpassport)
	if allow == false {
		beego.Info("not allow favor")
		util.ResponseFavorLimit(this.Ctx.Output, this.Callback)
		return
	}

	// 再次查看该用户是否已经向好友集福
	userLog, _ := models.GetUserLogByPassportAndFPassport(this.Passport, fpassport)

	faccount, _ := util.GetAccountByPassport(fpassport)
	if userLog.Id == 0 {
		//保存集福记录
		info := new(models.UserLog)
		info.Passport = this.Passport
		info.Fpassport = fpassport
		info.Day = util.GetDateInt(time.Now().Unix())
		info.Type = consts.USER_LOG_FAVOR
		info.CreateTime = util.GetDateMH(time.Now().Unix())
		info.Ip = this.GetRealIP()
		models.SaveUserLog(info)

		// 获取集福数量，判断是否可以获取黄金会员，并自动插入个人奖品
		count, _ := models.CountUserLogByPassportAndType(this.Passport, consts.USER_LOG_FAVOR)
		if count >= consts.FAVOR_VIP_3 {
			info := new(models.UserAward)
			info.Passport = this.Passport
			info.UserName = faccount.NickName
			info.Type = consts.PRODUCT_TYPE_VIP
			info.Status = consts.AWARD_STATUS_UNDRAW
			info.CreateTime = util.GetDateMH(time.Now().Unix())
			info.ReceiveTime = info.CreateTime
			info.UpdateTime = info.CreateTime
			if count >= consts.FAVOR_VIP_14 {
				info.ProductId = 3
			} else if count >= consts.FAVOR_VIP_7 {
				info.ProductId = 2
			} else if count >= consts.FAVOR_VIP_3 {
				info.ProductId = 1
			}
			models.SaveUserAward(info)
		}
		//更新user favor_num
		if fuser.Id > 0 {
			fuser.UpdateTime = util.GetDateMH(time.Now().Unix())
			fuser.FavorNum = count
			models.UpdateUser(&fuser)
		}
	} else {
		util.ResponseUserHasFavor(this.Ctx.Output, this.Callback)
		return
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
	return
}

// 判断是否已经向好友集福
func (this *FavorLightController) HasFavor() {
	fpassport := strings.Trim(this.GetString("fpassport"), " ")
	if "" == this.Passport && "" == fpassport {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	userLog, _ := models.GetUserLogByPassportAndFPassport(this.Passport, fpassport)
	if userLog.Id > 0 {
		data := map[string]interface{}{"favor": 1}
		util.ResponseObj(this.Ctx.Output, data, this.Callback)
		return
	} else {
		data := map[string]interface{}{"favor": 0}
		util.ResponseObj(this.Ctx.Output, data, this.Callback)
		return
	}
}

// 点亮ip
func (this *FavorLightController) LightUpIp() {
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

	user, _ := models.GetUserByPassport(this.Passport)
	if user.Id == 0 {
		//用户第一次参加，保存用户
		info := new(models.User)
		info.Headpic = accountInfo.SmallImg
		info.Nickname = accountInfo.NickName
		info.Userid = accountInfo.Uid
		info.Passport = this.Passport
		info.Status = consts.USER_STATUS_IPING
		info.IpNum = 1
		info.FavorNum = 0
		phoneInfo, _ := util.IsBindPhone(this.Passport)
		info.Phone = phoneInfo.Mobile
		info.CreateTime = util.GetDateMH(time.Now().Unix())
		info.UpdateTime = info.CreateTime
		models.SaveUser(info)
	}
	// 获取点亮日志记录
	now := time.Now().Unix()
	day := util.GetDateInt(now)
	userLog, _ := models.GetUserSignLog(this.Passport, consts.USER_LOG_IP, day)

	if userLog.Id > 0 {
		util.ResponseUserHasLightUp(this.Ctx.Output, this.Callback)
		return
	} else {
		if user.Id > 0 {
			//更新用户点亮IP数量
			user.UpdateTime = util.GetDateMH(time.Now().Unix())
			user.IpNum = user.IpNum + 1
			models.UpdateUser(&user)
		}
		// 保存点亮日志记录
		info := new(models.UserLog)
		info.Passport = this.Passport
		info.Day = day
		info.Type = consts.USER_LOG_IP
		info.CreateTime = util.GetDateMH(now)
		info.Ip = this.GetRealIP()
		models.SaveUserLog(info)
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
	return
}

func (this *FavorLightController) AllowFavor(passport string, fpassport string) bool {
	ip := this.GetRealIP()
	remaid := util.GetRemainSecond()
	remaidSecond, _ := strconv.Atoi(strconv.FormatInt(remaid, 10))
	//判断集福人的限制情况
	var dayLimitFavor int64
	dayLimitKeyFavor := fmt.Sprintf(consts.CACHE_ACCOUNT_FAVOR_LIMIT_DAY, passport)
	var ipLimitFavor int64
	ipLimitKeyFavor := fmt.Sprintf(consts.CACHE_ACCOUNT_FAVOR_LIMIT_IP, passport, ip)

	redis.Get(dayLimitKeyFavor, &dayLimitFavor)

	if dayLimitFavor == 0 || dayLimitFavor < consts.DAY_LIMIT {
		redis.Set(dayLimitKeyFavor, dayLimitFavor+1, remaidSecond)
	} else {
		beego.Info("Favor cache day request up limit")
		return false
	}
	redis.Get(ipLimitKeyFavor, &ipLimitFavor)
	if ipLimitFavor == 0 || ipLimitFavor < consts.IP_LIMIT {
		redis.Set(ipLimitKeyFavor, ipLimitFavor+1, remaidSecond)
	} else {
		beego.Info("Favor cache ip request up limit")
		return false
	}
	//判断被集福人的限制情况
	var dayLimitBeFavor int64
	dayLimitKeyBeFavor := fmt.Sprintf(consts.CACHE_ACCOUNT_FAVOR_LIMIT_DAY, fpassport)
	var ipLimitBeFavor int64
	ipLimitKeyBeFavor := fmt.Sprintf(consts.CACHE_ACCOUNT_FAVOR_LIMIT_IP, fpassport, ip)

	redis.Get(dayLimitKeyBeFavor, &dayLimitBeFavor)

	if dayLimitBeFavor == 0 || dayLimitBeFavor < consts.DAY_LIMIT {
		redis.Set(dayLimitKeyBeFavor, dayLimitBeFavor+1, remaidSecond)
	} else {
		beego.Info("BeFavor cache day request up limit")
		return false
	}
	redis.Get(ipLimitKeyBeFavor, &ipLimitBeFavor)
	if ipLimitBeFavor == 0 || ipLimitBeFavor < consts.IP_LIMIT {
		redis.Set(ipLimitKeyBeFavor, ipLimitBeFavor+1, remaidSecond)
	} else {
		beego.Info("BeFavor cache ip request up limit")
		return false
	}
	return true
}
