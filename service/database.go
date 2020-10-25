package service

import (
	"fmt"
	"gitea.com/ha666/sync-mysql/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ha666/logs"
	"xorm.io/xorm"
)

var (
	SourceEngine *xorm.Engine
	TargetEngine *xorm.Engine
)

func InitDataBases() {

	//region 初始化源库
	{
		var err error
		connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Asia%%2fShanghai",
			config.Conf.DataBases.Source.Account,
			config.Conf.DataBases.Source.Password,
			config.Conf.DataBases.Source.Address,
			config.Conf.DataBases.Source.Port,
			config.Conf.DataBases.Source.Name)
		SourceEngine, err = xorm.NewEngine("mysql", connString)
		if err != nil {
			logs.Emergency("源库连接失败:%s", err.Error())
		}
		if err = SourceEngine.Ping(); err != nil {
			logs.Emergency("源库Ping失败:%s", err.Error())
			return
		}
		SourceEngine.SetMaxIdleConns(2)
		SourceEngine.SetMaxOpenConns(50)
		logs.Info("初始化源库成功")
	}
	//endregion

	//region 初始化目标库
	{
		var err error
		connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Asia%%2fShanghai",
			config.Conf.DataBases.Target.Account,
			config.Conf.DataBases.Target.Password,
			config.Conf.DataBases.Target.Address,
			config.Conf.DataBases.Target.Port,
			config.Conf.DataBases.Target.Name)
		TargetEngine, err = xorm.NewEngine("mysql", connString)
		if err != nil {
			logs.Emergency("目标库连接失败:%s", err.Error())
		}
		if err = TargetEngine.Ping(); err != nil {
			logs.Emergency("目标库Ping失败:%s", err.Error())
			return
		}
		TargetEngine.SetMaxIdleConns(2)
		TargetEngine.SetMaxOpenConns(50)
		logs.Info("初始化目标库成功")
	}
	//endregion

}
