package util

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.SetLevel(beego.LevelInformational)
	beego.SetLogFuncCall(true)
	beego.SetLogger("file", `{"filename":"logs/stdout.log","maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	//	logs.SetLogger(logs.AdapterFile, `{"filename":"stdout.log","maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

}
