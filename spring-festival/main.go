package main

import (
	"festival-fun/job"
	_ "festival-fun/redis"
	_ "festival-fun/routers"
	"festival-fun/util"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

func main() {
	env := os.Getenv("GOMODE")
	if env != "" && env == "pro" {
		beego.BConfig.RunMode = "pro"
	}
	beego.Info("go runmode :", env)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	util.InitDB()
	job.StatJob()
	toolbox.StartTask()
	defer toolbox.StopTask()
	beego.Run()
}
