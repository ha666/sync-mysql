package main

import (
	"gitea.com/ha666/sync-mysql/config"
	"gitea.com/ha666/sync-mysql/model"
	"gitea.com/ha666/sync-mysql/service"
	"github.com/ha666/logs"
)

func main() {
	logs.Info("源库:%s", config.Conf.DataBases.Source.Name)
	logs.Info("目标库:%s", config.Conf.DataBases.Target.Name)
	logs.Info(service.GetTableList(model.DataBaseSource))
	logs.Info(service.GetTableList(model.DataBaseTarget))
}
