package service

import (
	"fmt"
	"gitea.com/ha666/sync-mysql/model"
	"xorm.io/xorm/schemas"
)

func getTableSchemaList(dataBaseType model.DataBaseType) (schemaMap map[string]*schemas.Table, err error) {
	var (
		schemaSlice []*schemas.Table
	)
	switch dataBaseType {
	default:
		return schemaMap, fmt.Errorf("不支持的数据库类型:%v", dataBaseType)
	case model.DataBaseSource:
		schemaSlice, err = sourceEngine.DBMetas()
	case model.DataBaseTarget:
		schemaSlice, err = targetEngine.DBMetas()
	}
	if err != nil || schemaSlice == nil || len(schemaSlice) <= 0 {
		return
	}
	schemaMap = make(map[string]*schemas.Table, len(schemaSlice))
	for _, s := range schemaSlice {
		schemaMap[s.Name] = s
	}
	return
}
