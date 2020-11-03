package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitea.com/ha666/sync-mysql/model"
	"gitea.com/ha666/sync-mysql/plugin/kafka"
	"strconv"
	"time"

	"gitea.com/ha666/sync-mysql/config"
	"github.com/ha666/golibs"
	"github.com/ha666/logs"
	"xorm.io/xorm/schemas"
)

type sqlAndArgs struct {
	Sql  string
	Args []interface{}
}

func StartDBWrite() {
	if config.Conf.Source.Database == nil {
		return
	}
	for tn, _ := range sourceSchemaColumns {
		logs.Info("开始读取源表:%s", tn)
		offset := 0
		for {
			if len(receiveQueue) > 20000 {
				logs.Info("队列长度达到%d，暂停10秒", len(receiveQueue))
				time.Sleep(10 * time.Second)
			}
			logs.Info("查询表:%s,offset:%d,limit:%d", tn, offset, config.Conf.App.PageSize)
			result, err := sourceEngine.QueryInterface(fmt.Sprintf("select * from %s limit %d,%d", tn, offset, config.Conf.App.PageSize))
			if err != nil {
				logs.Emergency("查询表:%s出错:%s", tn, err.Error())
			}
			rows, err := parseSourceSchema(tn, result)
			if err != nil {
				logs.Emergency("解析数据出错:%s", err.Error())
			}
			if rows <= 0 {
				break
			}
			offset = offset + config.Conf.App.PageSize
		}
		logs.Info("完成读取源表:%s", tn)
	}
}

func StartKafkaWrite() {
	if config.Conf.Source.Kafka == nil {
		return
	}
	ctx := context.TODO()
	for {
		err := kafka.KConsume(ctx, sourceKafkaName, config.Conf.Source.Kafka.Topic, func(data []byte) error {
			logs.Info("消息内容:%s", golibs.SliceByteToString(data))
			var msg model.DataChangeMsg
			if err := json.Unmarshal(data, &msg); err != nil {
				logs.Error("解析消息错误:%s", err.Error())
				return err
			}
			return parseDataChangeMsg(msg)
		})
		if err != nil {
			logs.Error("接收消息错误:%s", err.Error())
		}
	}
}

func parseDataChangeMsg(msg model.DataChangeMsg) error {
	if msg.Type != "insert" && msg.Type != "update" {
		return errors.New("暂时不支持的消息类型:" + msg.Type)
	}
	if msg.Database != config.Conf.Source.Kafka.DatabaseName {
		return nil
	}

	sliceOrder := make([]string, 0)
	for k, _ := range msg.Data {
		sliceOrder = append(sliceOrder, k)
	}

	data := golibs.NewStringBuilder()
	data.Append("insert into ").Append(msg.Table)
	data.Append("(")

	index := 0
	for _, o := range sliceOrder {
		if _, ok := msg.Data[o]; ok {
			if index > 0 {
				data.Append(",")
			}
			data.Append("`").Append(o).Append("`")
			index++
		}
	}

	data.Append(") values(")
	index = 0
	args := make([]interface{}, 0)
	for _, o := range sliceOrder {
		if v, ok := msg.Data[o]; ok {
			if index > 0 {
				data.Append(",")
			}
			data.Append("?")
			args = append(args, v)
			index++
		}
	}

	data.Append(")")
	index = 0

	if len(args) > 0 {
		receiveQueue <- &sqlAndArgs{
			Sql:  data.ToString(),
			Args: args,
		}
	}
	return nil
}

func parseSourceSchema(tableName string, result []map[string]interface{}) (int, error) {
	if result == nil || len(result) <= 0 {
		return 0, nil
	}
	count := 0
	if columns, ok1 := sourceSchemaColumns[tableName]; !ok1 {
		return 0, fmt.Errorf("不支持的表结构:%s", tableName)
	} else {
		if columns == nil || len(columns) <= 0 {
			return 0, fmt.Errorf("无效的表结构:%s", tableName)
		}
		count = len(result)
		for _, vs := range result {
			data := golibs.NewStringBuilder()
			data.Append("insert into ").Append(tableName)
			data.Append("(")
			index := 0
			for _, col := range columns {
				if _, ok2 := vs[col.Name]; ok2 {
					if index > 0 {
						data.Append(",")
					}
					data.Append("`").Append(col.Name).Append("`")
					index++
				}
			}
			data.Append(") values(")
			index = 0
			for _, col := range columns {
				if _, ok2 := vs[col.Name]; ok2 {
					if index > 0 {
						data.Append(",")
					}
					data.Append("?")
					index++
				}
			}
			data.Append(")")
			index = 0
			args := make([]interface{}, 0)
			for _, col := range columns {
				if cs, ok2 := vs[col.Name]; ok2 {
					switch col.SQLType.Name {
					default:
						logs.Emergency("表:%s,字段:%s,不支持的类型:%s", tableName, col.Name, col.SQLType.Name)
					case schemas.Bit, schemas.TinyInt, schemas.SmallInt, schemas.MediumInt, schemas.Int, schemas.Integer, schemas.BigInt:
						switch m := cs.(type) {
						default:
							logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, cs, cs)
						case []uint8:
							obj, err := strconv.ParseInt(golibs.SliceByteToString(m), 10, 64)
							if err != nil {
								logs.Emergency("表:%s,字段:%s,值:%v,%T,错误:%s", tableName, col.Name, cs, cs, err.Error())
							}
							args = append(args, obj)
						case nil:
							args = append(args, 0)
						}
					case schemas.Char, schemas.Varchar, schemas.NChar, schemas.NVarchar, schemas.TinyText, schemas.Text, schemas.NText, schemas.Clob, schemas.MediumText, schemas.LongText, schemas.Uuid, schemas.UniqueIdentifier, schemas.SysName:
						switch m := cs.(type) {
						default:
							logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, cs, cs)
						case []uint8:
							args = append(args, golibs.SliceByteToString(m))
						case nil:
							args = append(args, nil)
						}
					case schemas.DateTime:
						switch m := cs.(type) {
						default:
							logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, cs, cs)
						case time.Time:
							args = append(args, m.Format(golibs.Time_TIMEStandard))
						case nil:
							args = append(args, nil)
						}
					case schemas.TimeStamp:
						switch m := cs.(type) {
						default:
							logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, cs, cs)
						case time.Time:
							args = append(args, m.Format(golibs.Time_TIMEStandard))
						case nil:
							args = append(args, nil)
						}
					}
				}
			}
			if len(args) > 0 {
				receiveQueue <- &sqlAndArgs{
					Sql:  data.ToString(),
					Args: args,
				}
			}
		}
	}
	return count, nil
}
