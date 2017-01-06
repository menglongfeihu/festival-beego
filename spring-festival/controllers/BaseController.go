package controllers

import (
	"festival-fun/util"
	"net"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	Passport string
	Callback string
	Token    string
}

func (this *BaseController) Prepare() {
	this.Passport = strings.Trim(this.GetString("passport", ""), " ")
	this.Token = strings.Trim(this.GetString("token", ""), " ")
	this.Callback = strings.Trim(this.GetString("callback", ""), " ")
	beego.Info("request params:passport=" + this.Passport + "|token=" + this.Token + "|callback=" + this.Callback)
	this.GetRealIP()
	if this.Passport != "" && strings.Index(this.Passport, "@") > 0 {
		userInfo, _ := util.GetAccountByPassport(this.Passport)
		if userInfo.Uid == 0 {
			ip := this.GetRealIP()
			userInfo, _ = util.CreateAccount(this.Passport, ip)
			beego.Info("new add account:", userInfo)
		}
	}
}

func (this BaseController) FesvivalTimeRange() {
	startTime := util.GetTimeParse(beego.AppConfig.String("festival.start"))
	endTime := util.GetTimeParse(beego.AppConfig.String("festival.end"))
	now := time.Now().Unix()
	flag := 0
	if now < startTime {
		flag = 1 // no begin
	}
	if now > endTime {
		flag = 2 // has over
	}
	data := map[string]interface{}{
		"time":   now,
		"status": flag,
	}
	util.ResponseObj(this.Ctx.Output, data, this.Callback)
}

// 获取RealIp
func (this *BaseController) GetRealIP() string {
	request := this.Ctx.Request
	ip := request.Header.Get("ip")
	if "" == ip {
		ip = request.Header.Get("X-Forwarded-For")
	}
	if "" == ip {
		ip = request.Header.Get("X-Real-IP")
	}
	if "" == ip {
		ip = request.Header.Get("Proxy-Client-IP")
	}
	if "" == ip {
		ip = request.Header.Get("WL-Proxy-Client-IP")
	}
	if "" == ip {
		ip = request.RemoteAddr //format  Ip:Port
		ip_port := strings.Split(ip, ":")
		if len(ip_port) > 1 {
			ip = ip_port[0]
		}
		if ip == "127.0.0.1" {
			addrs, err := net.InterfaceAddrs()
			if err == nil {
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							ip = ipnet.IP.String()
							break
						}
					}
				}
			}
		}
	}
	beego.Info("ip=" + ip)
	if len(ip) > 15 && strings.Index(ip, ",") > 0 {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			ip = ips[0]
		}
	}
	return ip
}
