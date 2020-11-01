package service

import (
	"fmt"
	"gitea.com/ha666/sync-mysql/plugin/kafka"
	"github.com/ha666/golibs"
	"strings"

	"gitea.com/ha666/sync-mysql/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ha666/logs"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

var (
	sourceEngine        *xorm.Engine                   //源库
	targetEngines       []*xorm.Engine                 //目标库
	sourceSchemaColumns map[string][]*schemas.Column   //源库数据库结构
	targetSchemaColumns []map[string][]*schemas.Column //目标库数据库结构
	sourceKafkaName     string                         //源Kafka名称
	targetKafkaName     string                         //目标Kafka名称
)

//初始化数据库
func InitDataBases() {

	//region 初始化源库
	{
		if config.Conf.Source.Database != nil {
			var err error
			connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Asia%%2fShanghai",
				config.Conf.Source.Database.Account,
				config.Conf.Source.Database.Password,
				config.Conf.Source.Database.Address,
				config.Conf.Source.Database.Port,
				config.Conf.Source.Database.Name)
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
			logs.Info("初始化源库(%s)成功", config.Conf.Source.Database.Name)
		}
	}
	//endregion

	//region 初始化目标库
	{
		targetEngines = make([]*xorm.Engine, 0)
		if config.Conf.Target.Databases != nil && len(config.Conf.Target.Databases) > 0 {
			for i, d := range config.Conf.Target.Databases {
				connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Asia%%2fShanghai",
					d.Account,
					d.Password,
					d.Address,
					d.Port,
					d.Name)
				tmpEngine, err := xorm.NewEngine("mysql", connString)
				if err != nil {
					logs.Emergency("目标库(%d)连接失败:%s", i, err.Error())
				}
				if err = tmpEngine.Ping(); err != nil {
					logs.Emergency("目标库(%d)Ping失败:%s", i, err.Error())
					return
				}
				tmpEngine.SetMaxIdleConns(2)
				tmpEngine.SetMaxOpenConns(50)
				targetEngines = append(targetEngines, tmpEngine)
				logs.Info("初始化目标库(%d)-(%s)成功", i, d.Name)
			}
		}
	}
	//endregion

}

//检查数据库
func CheckDataBases() {

	//region 检查源数据库
	{
		if config.Conf.Source.Database != nil {
			logs.Info("开始检查源数据库")
			var err error
			sourceSchemaColumns, err = getTableSchemaList(sourceEngine)
			if err != nil {
				logs.Emergency("查询源数据库结构出错:%s", err.Error())
			}
			if sourceSchemaColumns == nil || len(sourceSchemaColumns) <= 0 {
				logs.Emergency("查询源数据库结构出错:空的")
			}
		}
	}
	//endregion

	//region 检查目标数据库
	{
		logs.Info("开始检查目标数据库")
		targetSchemaColumns = make([]map[string][]*schemas.Column, 0)
		if len(targetEngines) > 0 {
			for i, e := range targetEngines {
				tmpColumns, err := getTableSchemaList(e)
				if err != nil {
					logs.Emergency("查询目标数据库(%d)结构出错:%s", i, err.Error())
				}
				if tmpColumns == nil || len(tmpColumns) <= 0 {
					logs.Emergency("查询目标数据库(%d)结构出错:空的", i)
				}
				targetSchemaColumns = append(targetSchemaColumns, tmpColumns)
			}
		}
		logs.Info("完成检查目标数据库")
	}
	//endregion

	//region 检查数据库匹配程度
	{
		for sn, ss := range sourceSchemaColumns {
			for tsi, tsc := range targetSchemaColumns {
				if ts, ok := tsc[sn]; !ok {
					logs.Warn("目标数据库(%d)中不存在表:%s", tsi, sn)
				} else {
					for _, si := range ss {
						isExist := false
						for _, ti := range ts {
							if si.Name == ti.Name {
								if si.SQLType.Name != ti.SQLType.Name ||
									si.SQLType.DefaultLength != ti.SQLType.DefaultLength ||
									si.SQLType.DefaultLength2 != ti.SQLType.DefaultLength2 {
									logs.Warn("目标数据库(%d),字段:%s字段类型不一致,%s!=%s,%d!=%d,%d!=%d",
										tsi,
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
							logs.Warn("目标数据库(%d),表:%s,目标表中不存在字段:%s", tsi, sn, si.Name)
						}
					}
				}
			}
		}
	}
	//endregion

}

//初始化Kafka
func InitKafkas() {

	//初始化源Kafka
	{
		if config.Conf.Source.Kafka != nil {
			sourceKafkaName = golibs.Md5(strings.Join(config.Conf.Source.Kafka.Addresses, ",") + config.Conf.Source.Kafka.Topic)
			kafka.InitConsumer(sourceKafkaName,
				config.Conf.Source.Kafka.Addresses,
				config.Conf.Source.Kafka.Topic,
				config.Conf.Source.Kafka.Version,
				config.Conf.Source.Kafka.Consumer)
		}
	}
	//endregion

	//初始化目标Kafka
	{
		if config.Conf.Target.Kafka != nil {
			targetKafkaName = golibs.Md5(strings.Join(config.Conf.Target.Kafka.Addresses, ",") + config.Conf.Target.Kafka.Topic)
			if sourceKafkaName == targetKafkaName {
				logs.Emergency("源kafka和目标kafka配置不能相同")
			}
			kafka.InitProducer(targetKafkaName,
				config.Conf.Source.Kafka.Addresses)
		}
	}
	//endregion

}
