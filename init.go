package main

import (
	"runtime"

	"gitea.com/ha666/sync-mysql/config"
	"gitea.com/ha666/sync-mysql/service"
	"github.com/ha666/golibs"
	"github.com/ha666/logs"
)

const version = "2020.1028.1015"

func init() {
	initLog()
	initEnv()
	initConfig()
	service.InitDataBases()
	service.CheckDataBases()
}

func initLog() {
	logs.SetLogger(logs.AdapterConsole, `{"level":7}`)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.Info("初始化日志成功")
}

func initEnv() {

	//region 输出当前版本
	logs.Info("当前版本:<%s>---<%s>", version, runtime.Version())
	//endregion

	//region 输出系统信息
	logs.Info("os:%s", runtime.GOOS)
	logs.Info("cpu:%d", runtime.NumCPU())
	//endregion

	//region 输出网络信息
	logs.Info("ip:%s", golibs.GetCurrentIntranetIP())
	//endregion

	//region 输出应用信息
	logs.Info("path:%s", golibs.GetCurrentDirectory())
	//endregion

}

func initConfig() {
	err := config.Parser()
	if err != nil {
		logs.Emergency("读取配置文件出错:%s", err.Error())
	}
	if config.Conf.App.PageSize < 0 || config.Conf.App.PageSize > 10000 {
		logs.Emergency("配置出错:page_size范围是1~10000")
	}
	if config.Conf.App.ThreadCount < 0 || config.Conf.App.ThreadCount > 16 {
		logs.Emergency("配置出错:thread_count范围是1~16")
	}
	logs.Info("初始化配置成功")
}
