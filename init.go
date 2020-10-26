package main

import (
	"gitea.com/ha666/sync-mysql/config"
	"gitea.com/ha666/sync-mysql/service"
	"github.com/ha666/golibs"
	"github.com/ha666/logs"
	"runtime"
)

const version = "2020.1026.1708"

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
	logs.Info("初始化配置成功")
}
