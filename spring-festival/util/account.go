package util

import (
	"encoding/json"
	"festival-fun/consts"
	"festival-fun/redis"
	"fmt"
	"time"

	"github.com/astaxie/beego"
)

type UserInfo struct {
	Status   int    `json:"status"`
	Uid      int64  `json:"uid"`
	NickName string `json:"nickname"`
	SmallImg string `json:"smallimg"`
	BigImg   string `json:bigimg`
}

func GetAccountByPassport(passport string) (UserInfo, error) {
	cacheKey := fmt.Sprintf(consts.CACHE_ACCOUNT_INFO, passport)
	var userInfo UserInfo
	var err error
	cacheResult := redis.Get(cacheKey, &userInfo)
	if cacheResult && userInfo.Uid > 0 {
		beego.Info("account info exist in cache,passport:", passport, ",userinfo:", userInfo)
		err = nil
	} else {
		params := "p=" + passport + "&coding=utf8"
		data := HttpGet(beego.AppConfig.String("getaccount") + "?" + params)
		beego.Info("account info from http,passport:", passport, ",userinfo:", data)
		err = json.Unmarshal([]byte(data), &userInfo)
		if userInfo.Uid > 0 {
			redis.Set(cacheKey, userInfo, consts.EXPIRE_ONE_DAY)
		}
	}
	return userInfo, err
}
func CreateAccount(passport, ip string) (UserInfo, error) {
	var userInfo UserInfo
	var err error
	params := "p=" + passport + "&coding=utf8"
	data := HttpGet(beego.AppConfig.String("createaccount") + "?" + params)
	beego.Info("account info ,passport:", passport, ",userinfo:", data)
	err = json.Unmarshal([]byte(data), &userInfo)
	if userInfo.Uid > 0 {
		cacheKey := fmt.Sprintf(consts.CACHE_ACCOUNT_INFO, passport)
		redis.Set(cacheKey, userInfo, consts.EXPIRE_ONE_DAY)
	} else {
		beego.Error("create account info failed:passport:", passport, ",ip:", ip)
	}
	return userInfo, err
}

type Attachment struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}
type CheckLoginJson struct {
	Status     int        `json:"status"`
	Message    string     `json:"message"`
	Attachment Attachment `json:"attachment"`
}

func CheckLogin(passport string, token string) bool {
	beego.Info("check Login, passport =" + passport + ", token =" + token)
	login := false

	if passport != "" && token != "" {
		cacheKey := fmt.Sprintf(consts.CACHE_ACCOUNT_LOGIN, passport)
		cacheResult := redis.Get(cacheKey, &login)
		if cacheResult {
			beego.Info("login status exist in cache,passport:", passport, ",login status:", login)
		} else {
			params := "passport=" + passport + "&token=" + token + "&poid=123&plat=17&partner=1&sysver=1&sver=1.0&api_key=f351515304020cad28c92f70f002261c"
			if data := HttpGet(beego.AppConfig.String("checklogin") + "?" + params); data != "" {
				beego.Info("login status from http,passport:", passport, ",logininfo:", data)
				var checkLoginJson CheckLoginJson
				err := json.Unmarshal([]byte(data), &checkLoginJson)
				if err == nil {
					if checkLoginJson.Status == 200 && checkLoginJson.Attachment.Status == 0 {
						login = true
					}
				}
			}
			redis.Set(cacheKey, login, consts.EXPIRE_ONE_MINUTE)
		}
	}
	return login
}

type BindPhone struct {
	Mobile string `json:"mobile"`
}

func IsBindPhone(passport string) (BindPhone, error) {
	var bindPhone BindPhone
	var err error
	cacheKey := fmt.Sprintf(consts.CACHE_ACCOUNT_PHONE, passport)
	cacheResult := redis.Get(cacheKey, &bindPhone)
	if cacheResult && bindPhone.Mobile != "" {
		beego.Info("phoneinfo exist in cache,passport:", passport, ",phoneinfo:", bindPhone)
		err = nil
	} else {
		verifyPhone := beego.AppConfig.String("verify_phone")
		params := "p=" + passport
		data := HttpGet(verifyPhone + "?" + params)
		beego.Info("phoneinfo from http,passport:", passport, ",verify result:", data)
		err = json.Unmarshal([]byte(data), &bindPhone)
		if bindPhone.Mobile != "" {
			redis.Set(cacheKey, bindPhone, consts.EXPIRE_ONE_DAY)
		}
	}
	return bindPhone, err
}

type VipVideo struct {
	Status     int        `json:"status"`
	StatusText string     `json:"statusText"`
	Privileges Privileges `json:"data"`
}
type Privileges struct {
	Privilege []Privilege `json:"privileges"`
}
type Privilege struct {
	Amount int    `json:"amount"`
	CateId int    `json:"cateId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Unit   string `json:"unit"`
}

// vip 充值
// activityId 活动ID
// channelId 激活入口
func VipVideoRecharge(passport, activityId, channelId string) (bool, string) {
	params := map[string]string{
		"activity_id": activityId,
		"channel_id":  channelId,
		"passport":    passport,
	}
	headers := map[string]string{
		"app_id": "1",
		"plat":   "16",
		"gid":    "abc",
	}

	reqResult := HttpPost(beego.AppConfig.String("vip_rechargeurl"), params, headers, time.Second*5)

	var vip VipVideo
	json.Unmarshal([]byte(reqResult), &vip)
	beego.Info("VipVideoRecharge result=", reqResult)
	if vip.Status != 200 {
		return false, vip.StatusText
	}
	return true, ""
}
