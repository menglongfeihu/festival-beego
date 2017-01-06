package controllers

import (
	"festival-fun/consts"
	"festival-fun/models"
	"festival-fun/util"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type TestController struct {
	beego.Controller
	Passport string
	Callback string
	Token    string
	/*
		IsLogin      bool
		UserUserId   int64
		UserUsername string
		UserAvatar   string
	*/
}

func (this *TestController) Prepare() {
	this.Passport = strings.Trim(this.GetString("passport"), " ")
	this.Token = strings.Trim(this.GetString("token"), " ")
	this.Callback = strings.Trim(this.GetString("callback"), " ")
	/*
		accountInfo, err := util.GetAccountByPassport(this.Passport)
		if nil == err {
			this.UserUserId = accountInfo.Uid
			this.UserUsername = accountInfo.NickName
			this.UserAvatar = accountInfo.SmallImg
		} else {
			beego.Error("GetAccountByPassport error", err)
			util.ResponseAccountNotExist(this.Ctx.Output, this.Callback)
			return
		}
		isLogin := util.CheckLogin(this.Passport, this.Token)
		if isLogin {
			this.IsLogin = true
		} else {
			this.IsLogin = false
		}
	*/
}

// 判断用户是否绑定手机号
func (this *TestController) IsBindPhone() {
	if "" == this.Passport {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	accountInfo, _ := util.GetAccountByPassport(this.Passport)
	if accountInfo.Uid == 0 {
		util.ResponseUserNotExist(this.Ctx.Output, this.Callback)
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
}

// 获取实物奖品列表
func (this *TestController) GetProductList() {
	data, err := models.ListProductByType(consts.PRODUCT_TYPE_ENTITY)
	if nil != err {
		beego.Error("Get Product list failed:", err)
		util.ResponseSystemErr(this.Ctx.Output, this.Callback)
		return
	}
	util.ResponseObj(this.Ctx.Output, data, this.Callback)
}

// 订阅专辑
func (this *TestController) SubscribeVideoInfo() {
	vname := strings.Trim(this.GetString("vname"), " ")
	if "" == this.Passport || "" == vname || "" == this.Token {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
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
	videoSubinfo, _ := models.GetVideoSubinfoByPassportAndVname(this.Passport, vname)
	if videoSubinfo.Id > 0 {
		util.ResponseObj(this.Ctx.Output, "你已经订阅《"+vname+"》", this.Callback)
		return
	} else {
		video := new(models.VideoSubinfo)
		video.Passport = this.Passport
		video.Vname = vname
		video.CreateTime = util.GetDateMH(time.Now().Unix())
		models.SaveVideoSubInfo(video)
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
}

//充值搜狐视频黄金会员
func (this *TestController) RechargeVIPMember() {
	vipType, _ := this.GetInt("type", 0)

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
	count64, _ := models.CountUserLogByPassportAndType(this.Passport, consts.USER_LOG_FAVOR)
	count, _ := strconv.Atoi(strconv.FormatInt(count64, 10))
	if count < 3 {
		util.ResponseObj(this.Ctx.Output, "集福数还不够，无法领取搜狐视频黄金会员", this.Callback)
		return
	}
	var proid int64
	switch vipType {
	case consts.VIP_TYPE_14:
		beego.Info("recharge 14 vip")
		proid = 3
	case consts.VIP_TYPE_7:
		beego.Info("recharge 7 vip")
		proid = 2
	case consts.VIP_TYPE_3:
		beego.Info("recharge 3 vip")
		proid = 1
	default:
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	userAward, _ := models.GetUseAwardByPassportAndProId(this.Passport, proid)

	if userAward.Id > 0 && userAward.Status == 0 {
		util.ResponseObj(this.Ctx.Output, "你已经领取搜狐视频"+strconv.Itoa(vipType)+"天黄金会员", this.Callback)
		return
	}
	userAward.Status = 0
	userAward.UpdateTime = util.GetDateMH(time.Now().Unix())
	userAward.ReceiveTime = userAward.UpdateTime
	models.UpdateUseAward(&userAward)
	util.ResponseOK(this.Ctx.Output, this.Callback)
}

// 获取集福好友头像列表
func (this *TestController) GetFrendsHeadPicList() {
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
}

// 向好友集福
func (this *TestController) FavorFriend() {
	fpassport := strings.Trim(this.GetString("fpassport"), " ")
	if "" == this.Passport && "" == fpassport {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}
	// 判断用户的好友是否已经点亮6个ip，只有点亮6个ip才可以进行集福
	count64, _ := models.CountUserLogByPassportAndType(fpassport, consts.USER_LOG_IP)
	count, _ := strconv.Atoi(strconv.FormatInt(count64, 10))
	if count < consts.COUNT_IP {
		util.ResponseObj(this.Ctx.Output, "你的好友还差"+strconv.Itoa(consts.COUNT_IP-count)+"个IP,才可以集福~", this.Callback)
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
		if count >= 3 {
			info := new(models.UserAward)
			info.Passport = this.Passport
			info.UserName = faccount.NickName
			info.Type = 0
			info.Status = 1
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
		user, _ := models.GetUserByPassport(fpassport)
		if user.Id > 0 {
			user.UpdateTime = util.GetDateMH(time.Now().Unix())
			user.FavorNum = count
			user.Status = consts.USER_STATUS_FAVOR
			models.UpdateUser(&user)
		}
	} else {
		util.ResponseObj(this.Ctx.Output, "你已经为你的好友:"+faccount.NickName+"集福啦", this.Callback)
		return
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
}

// 判断是否已经向好友集福
func (this *TestController) HasFavor() {
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

// 获取我的奖品列表
func (this *TestController) GetAwardsList() {
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
	user, err := models.GetUserByPassport(this.Passport)
	if err != nil {
		util.ResponseObj(this.Ctx.Output, "需要先点亮IP，集福，才有机会获得奖品", this.Callback)
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
}

// 填写领取实物奖品地址信息
func (this *TestController) AppendAwardsAddress() {
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
	pid, _ := this.GetInt64("pid", 0)
	name := strings.Trim(this.GetString("name"), " ")
	addr := strings.Trim(this.GetString("addr"), " ")
	phone := strings.Trim(this.GetString("phone"), " ")
	qq := strings.Trim(this.GetString("qq"), " ")
	email := strings.Trim(this.GetString("email"), " ")
	if "" == this.Passport || 0 == pid {
		util.ResponseLackParam(this.Ctx.Output, this.Callback)
		return
	}

	userAward, _ := models.GetUseAwardByPassportAndProId(this.Passport, pid)

	if userAward.Id > 0 {
		userAward.Address = addr
		userAward.UserName = name
		userAward.Phone = phone
		userAward.Email = email
		userAward.QQ = qq
		userAward.UpdateTime = util.GetDateMH(time.Now().Unix())
		models.UpdateUseAward(&userAward)
	} else {
		util.ResponseSystemErr(this.Ctx.Output, this.Callback)
		return
	}
	util.ResponseOK(this.Ctx.Output, this.Callback)
}

func (this *TestController) Test() {
	models.GetVideoSubinfoByPassportAndVname("o8u7xju3grT9c7YuaqI6bmDhVM6w@wechat.sohu.com", "法医秦明")
}

// 点亮ip
func (this *TestController) LightUpIp() {
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
		//用户第一次参加，保存用户
		info := new(models.User)
		info.Headpic = accountInfo.SmallImg
		info.Nickname = accountInfo.NickName
		info.Userid = accountInfo.Uid
		info.Passport = this.Passport
		info.Status = consts.USER_STATUS_IP
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
		util.ResponseObj(this.Ctx.Output, "你今天已经点亮", this.Callback)
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
}

// 获取RealIp
func (this *TestController) GetRealIP() string {
	request := this.Ctx.Request
	ip := request.Header.Get("ip")
	if "" == ip {
		ip = request.Header.Get("X-Forwarded-For")
	}
	if "" == ip {
		ip = request.Header.Get("X-Real-IP")
	}
	if "" == ip {
		ip = request.RemoteAddr
	}
	beego.Info("ip=" + ip)
	return ip
}
