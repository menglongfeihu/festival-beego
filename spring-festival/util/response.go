package util

import (
	"github.com/astaxie/beego/context"
)

func JsonResponse(output *context.BeegoOutput, status int, data interface{}, callback string) {
	output.ContentType("application/json; charset=UTF-8")
	output.Header("Cache-Control", "no-cache")
	output.Header("Pragma", "no-cache")
	var result = map[string]interface{}{"status": status, "message": data}
	if callback == "" {
		output.JSON(result, false, true)
	} else {
		output.JSONP(result, false)
	}

}

func ResponseOK(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 200, "success", callback)
}

func ResponseSystemErr(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 104, "system error", callback)
}

func ResponseLackParam(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 101, "lack param", callback)
}

func ResponseObj(output *context.BeegoOutput, data interface{}, callback string) {
	JsonResponse(output, 200, data, callback)
}

func ResponseUserNotExist(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 201, "user not exist", callback)
}

func ResponseAccountNotExist(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 202, "account not exist", callback)
}

func ResponseUserNoLogin(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 203, "user not login", callback)
}

func ResponseVipHasObtain(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 204, "user has obtain vip", callback)
}
func ResponseUserHasSubVideo(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 205, "user already sub this video", callback)
}

func ResponseUserHasAward(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 206, "user has awarded", callback)
}

func ResponseUserHasLightUp(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 207, "user has light up", callback)
}

func ResponseLightIPNotEnough(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 208, "user light up not enough", callback)
}

func ResponseUserNotFavor(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 209, "user not in favor status", callback)
}

func ResponseFavorNotEnough(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 210, "favor not enough", callback)
}

func ResponseUserHasFavor(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 211, "user has favor", callback)
}

func ResponseFavorLimit(output *context.BeegoOutput, callback string) {
	JsonResponse(output, 212, "favor up limit", callback)
}
