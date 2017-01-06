// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"festival-fun/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/festival/time/range.json", &controllers.BaseController{}, "*:FesvivalTimeRange")

	//UserController
	beego.Router("/festival/user/status.json", &controllers.UserController{}, "*:GetUserStatus")
	beego.Router("/festival/user/is_bind.json", &controllers.UserController{}, "*:IsBindPhone")
	beego.Router("/festival/video/sub.json", &controllers.UserController{}, "*:SubscribeVideoInfo")
	beego.Router("/festival/vip/get.json", &controllers.UserController{}, "*:RechargeVIPMember")
	//test
	//beego.Router("/festival/get.json", &controllers.UserController{}, "get:Get")

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
