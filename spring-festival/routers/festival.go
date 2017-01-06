package routers

import (
	"festival-fun/controllers"

	"github.com/astaxie/beego"
)

func init() {
	/*
		beego.Router("/festival/user/status.json", &controllers.UserController{}, "*:GetUserStatus")
		beego.Router("/festival/user/is_bind.json", &controllers.TestController{}, "*:IsBindPhone")
		beego.Router("/festival/product/list.json", &controllers.TestController{}, "*:GetProductList")
		beego.Router("/festival/video/sub.json", &controllers.TestController{}, "*:SubscribeVideoInfo")
		beego.Router("/festival/vip/get.json", &controllers.TestController{}, "*:RechargeVIPMember")
		beego.Router("/festival/friends/list.json", &controllers.TestController{}, "*:GetFrendsHeadPicList")
		beego.Router("/festival/favor/add.json", &controllers.TestController{}, "*:FavorFriend")
		beego.Router("/festival/favor/exist.json", &controllers.TestController{}, "*:HasFavor")
		beego.Router("/festival/user/awards.json", &controllers.TestController{}, "*:GetAwardsList")
		beego.Router("/festival/address/save.json", &controllers.TestController{}, "*:AppendAwardsAddress")
		beego.Router("/festival/ip/save.json", &controllers.TestController{}, "*:LightUpIp")
		beego.Router("/festival/test/test.json", &controllers.FestivalController{}, "*:Test")
	*/
	//UserController
	beego.Router("/festival/user/status.json", &controllers.UserController{}, "*:GetUserStatus")
	beego.Router("/festival/user/is_bind.json", &controllers.UserController{}, "*:IsBindPhone")
	beego.Router("/festival/video/sub.json", &controllers.UserController{}, "*:SubscribeVideoInfo")
	beego.Router("/festival/vip/get.json", &controllers.UserController{}, "*:RechargeVIPMember")

	//FavorLightController
	beego.Router("/festival/favor/open.json", &controllers.FavorLightController{}, "*:OpenFavor")
	beego.Router("/festival/friends/list.json", &controllers.FavorLightController{}, "*:GetFrendsHeadPicList")
	beego.Router("/festival/favor/add.json", &controllers.FavorLightController{}, "*:FavorFriend")
	beego.Router("/festival/favor/exist.json", &controllers.FavorLightController{}, "*:HasFavor")
	beego.Router("/festival/ip/save.json", &controllers.FavorLightController{}, "*:LightUpIp")

	//ProductController
	beego.Router("/festival/product/list.json", &controllers.ProductController{}, "*:GetProductList")
	beego.Router("/festival/user/awards.json", &controllers.ProductController{}, "*:GetAwardsList")
	beego.Router("/festival/address/save.json", &controllers.ProductController{}, "*:AppendAwardsAddress")

}
