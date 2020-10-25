package service

import (
	"fmt"
	"gitea.com/ha666/sync-mysql/config"
	"gitea.com/ha666/sync-mysql/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ha666/logs"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

var (
	sourceEngine    *xorm.Engine
	targetEngine    *xorm.Engine
	sourceSchemaMap map[string]*schemas.Table
	targetSchemaMap map[string]*schemas.Table
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
		sourceEngine, err = xorm.NewEngine("mysql", connString)
		if err != nil {
			logs.Emergency("源库连接失败:%s", err.Error())
		}
		if err = sourceEngine.Ping(); err != nil {
			logs.Emergency("源库Ping失败:%s", err.Error())
			return
		}
		sourceEngine.SetMaxIdleConns(2)
		sourceEngine.SetMaxOpenConns(50)
		logs.Info("初始化源库(%s)成功", config.Conf.DataBases.Source.Name)
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
		targetEngine, err = xorm.NewEngine("mysql", connString)
		if err != nil {
			logs.Emergency("目标库连接失败:%s", err.Error())
		}
		if err = targetEngine.Ping(); err != nil {
			logs.Emergency("目标库Ping失败:%s", err.Error())
			return
		}
		targetEngine.SetMaxIdleConns(2)
		targetEngine.SetMaxOpenConns(50)
		logs.Info("初始化目标库(%s)成功", config.Conf.DataBases.Target.Name)
	}
	//endregion

}

func CheckDataBases() {

	//region 检查源数据库
	{
		logs.Info("开始检查源数据库")
		var err error
		sourceSchemaMap, err = getTableSchemaList(model.DataBaseSource)
		if err != nil {
			logs.Emergency("查询源数据库结构出错:%s", err.Error())
		}
		if sourceSchemaMap == nil || len(sourceSchemaMap) <= 0 {
			logs.Emergency("查询源数据库结构出错:空的")
		}
		logs.Info("完成检查源数据库")
	}
	//endregion

	//region 检查目标数据库
	{
		logs.Info("开始检查目标数据库")
		var err error
		targetSchemaMap, err = getTableSchemaList(model.DataBaseTarget)
		if err != nil {
			logs.Emergency("查询目标数据库结构出错:%s", err.Error())
		}
		if targetSchemaMap == nil || len(targetSchemaMap) <= 0 {
			logs.Emergency("查询目标数据库结构出错:空的")
		}
		logs.Info("完成检查目标数据库")
	}
	//endregion

	//region 检查数据库匹配程度
	{
		for sn, ss := range sourceSchemaMap {
			logs.Info("开始检查表:%s", sn)
			if ts, ok := targetSchemaMap[sn]; !ok {
				logs.Warn("目标数据库中不存在表:%s", sn)
			} else {
				for _, si := range ss.Columns() {
					isExist := false
					for _, ti := range ts.Columns() {
						if si.Name == ti.Name {
							if si.SQLType.Name != ti.SQLType.Name ||
								si.SQLType.DefaultLength != ti.SQLType.DefaultLength ||
								si.SQLType.DefaultLength2 != ti.SQLType.DefaultLength2 {
								logs.Warn("字段:%s字段类型不一致,%s!=%s,%d!=%d,%d!=%d",
									si.Name,
									si.SQLType.Name, ti.SQLType.Name,
									si.SQLType.DefaultLength, ti.SQLType.DefaultLength,
									si.SQLType.DefaultLength2, ti.SQLType.DefaultLength2)
							}
							isExist = true
							break
						}
					}
					if !isExist {
						logs.Warn("表:%s,目标表中不存在字段:%s", ss.Name, si.Name)
					}
				}
			}
		}
	}
	//endregion

}
