package service

import (
	"fmt"

	"github.com/ha666/golibs"
	"github.com/ha666/logs"
	"github.com/xwb1989/sqlparser"
	"xorm.io/xorm/schemas"
)

func StartRead() {
	for get := range queue {
		if get == nil {
			break
		}
		stmt, err := sqlparser.Parse(get.Sql)
		if err != nil {
			logs.Error("解析sql:%s,失败:%s", get.Sql, err.Error())
			continue
		}
		switch stmt.(type) {
		default:
			break
		case *sqlparser.Insert:
			if err = processInsert(get.Sql, get.Args, stmt); err != nil {
				logs.Error(err.Error())
			}
		}
	}
}

func processInsert(Sql string, Args []interface{}, stmt sqlparser.Statement) error {
	var (
		ok              bool
		insertStatement *sqlparser.Insert
	)
	if insertStatement, ok = stmt.(*sqlparser.Insert); !ok {
		logs.Error("解析sql:%s,失败", Sql)
	} else if insertStatement == nil {
		logs.Error("解析sql:%s,失败", Sql)
	}
	if insertStatement == nil {
		return fmt.Errorf("解析sql:%s,失败,类型:%T", Sql, stmt)
	}
	if insertStatement.Action != "insert" {
		return fmt.Errorf("解析sql:%s,action出错:%s", Sql, insertStatement.Action)
	}
	if golibs.Length(insertStatement.Table.Name.String()) <= 0 {
		return fmt.Errorf("解析sql:%s,table.name出错:%s", Sql, insertStatement.Table.Name.String())
	}
	if insertStatement.Columns == nil || len(insertStatement.Columns) <= 0 || len(Args) <= 0 || len(insertStatement.Columns) != len(Args) {
		return fmt.Errorf("解析sql:%s,columns出错:%+v,%+v", Sql, insertStatement.Columns, Args)
	}
	tableName := insertStatement.Table.Name.String()
	if columns, ok1 := targetSchemaColumns[tableName]; !ok1 {
		return fmt.Errorf("不支持的表:%s", tableName)
	} else {
		if columns == nil || len(columns) <= 0 {
			return fmt.Errorf("无效的表:%s", tableName)
		}
		targetColumns := make(map[string]string, 0)
		for _, ic := range insertStatement.Columns {
			if golibs.Length(ic.String()) > 0 {
				targetColumns[ic.String()] = ic.String()
			}
		}
		if len(targetColumns) <= 0 {
			return fmt.Errorf("缺少有效的字段:%s", tableName)
		}
		data := golibs.NewStringBuilder()
		data.Append("insert into ").Append(tableName)
		data.Append("(")
		index := 0
		for _, col := range columns {
			if _, ok2 := targetColumns[col.Name]; ok2 {
				if index > 0 {
					data.Append(",")
				}
				data.Append(col.Name)
				index++
			}
		}
		data.Append(") values(")
		index = 0
		for _, col := range columns {
			if _, ok2 := targetColumns[col.Name]; ok2 {
				if index > 0 {
					data.Append(",")
				}
				data.Append("?")
				index++
			}
		}
		data.Append(") on duplicate key update ")
		index = 0
		for _, col := range columns {
			if col.IsPrimaryKey {
				continue
			}
			if _, ok2 := targetColumns[col.Name]; ok2 {
				if index > 0 {
					data.Append(", ")
				}
				data.Append(col.Name).Append("=?")
				index++
			}
		}
		index = 0
		args := make([]interface{}, 0)
		args = append(args, data.ToString())
		for _, col := range columns {
			if _, ok2 := targetColumns[col.Name]; ok2 {

				for ici, ic := range insertStatement.Columns {
					if col.Name == ic.String() {
						obj := Args[ici]
						switch col.SQLType.Name {
						default:
							logs.Emergency("表:%s,字段:%s,不支持的类型:%s", tableName, col.Name, col.SQLType.Name)
						case schemas.Bit, schemas.TinyInt, schemas.SmallInt, schemas.MediumInt, schemas.Int, schemas.Integer, schemas.BigInt:
							switch obj.(type) {
							default:
								logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, obj, obj)
							case int64:
								args = append(args, obj)
							}
						case schemas.Char, schemas.Varchar, schemas.NChar, schemas.NVarchar, schemas.TinyText, schemas.Text, schemas.NText, schemas.Clob, schemas.MediumText, schemas.LongText, schemas.Uuid, schemas.UniqueIdentifier, schemas.SysName:
							switch obj.(type) {
							default:
								logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, obj, obj)
							case string:
								args = append(args, obj)
							case nil:
								args = append(args, nil)
							}
						case schemas.DateTime, schemas.TimeStamp:
							switch obj.(type) {
							default:
								logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, obj, obj)
							case string:
								args = append(args, obj)
							case nil:
								args = append(args, nil)
							}
						}
						break
					}
				}

			}
		}
		index = 0
		for _, col := range columns {
			if col.IsPrimaryKey {
				continue
			}
			if _, ok2 := targetColumns[col.Name]; ok2 {

				for ici, ic := range insertStatement.Columns {
					if col.Name == ic.String() {
						obj := Args[ici]
						switch col.SQLType.Name {
						default:
							logs.Emergency("表:%s,字段:%s,不支持的类型:%s", tableName, col.Name, col.SQLType.Name)
						case schemas.Bit, schemas.TinyInt, schemas.SmallInt, schemas.MediumInt, schemas.Int, schemas.Integer, schemas.BigInt:
							switch obj.(type) {
							default:
								logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, obj, obj)
							case int64:
								args = append(args, obj)
							}
						case schemas.Char, schemas.Varchar, schemas.NChar, schemas.NVarchar, schemas.TinyText, schemas.Text, schemas.NText, schemas.Clob, schemas.MediumText, schemas.LongText, schemas.Uuid, schemas.UniqueIdentifier, schemas.SysName:
							switch obj.(type) {
							default:
								logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, obj, obj)
							case string:
								args = append(args, obj)
							case nil:
								args = append(args, nil)
							}
						case schemas.DateTime, schemas.TimeStamp:
							switch obj.(type) {
							default:
								logs.Emergency("表:%s,字段:%s,值:%v,无效的类型:%T", tableName, col.Name, obj, obj)
							case string:
								args = append(args, obj)
							case nil:
								args = append(args, nil)
							}
						}
						break
					}
				}

			}
		}
		if err := insertData(args); err != nil {
			return fmt.Errorf("插入数据:%+v,出错:%s", args, err.Error())
		}
	}
	return nil
}
