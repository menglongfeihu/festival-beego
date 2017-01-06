package util

import (
	"festival-fun/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() {
	// get the connection config of mysql
	beego.Info("======== begin init mysql ========")
	mysqlurl := beego.AppConfig.String("mysql.url")

	//	if "dev" == beego.AppConfig.String("runmode") {
	//		orm.Debug = true
	//	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", mysqlurl)
	if nil != err {
		panic("init mysql failed :" + err.Error())
	}
	//register data model
	mysqlprefix := beego.AppConfig.String("mysql.prefix")
	orm.RegisterModelWithPrefix(mysqlprefix,
		new(models.User),
		new(models.Product),
		new(models.UserAward),
		new(models.UserLog),
		new(models.UserRank),
		new(models.VideoSubinfo),
		new(models.RankAwardMap),
	)

	// get the config of the connection
	openCons, _ := beego.AppConfig.Int("mysql.maxOpenCons")
	idleCons, _ := beego.AppConfig.Int("mysql.maxIdleCons")
	beego.Info("init max_cons:", openCons, ",max_idles:", idleCons)
	orm.SetMaxOpenConns("default", openCons)
	orm.SetMaxIdleConns("default", idleCons)

	beego.Info("======== finish init mysql ========")
}
