package service

import (
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

func getTableSchemaList(e *xorm.Engine) (schemaMap map[string][]*schemas.Column, err error) {
	var (
		schemaSlice []*schemas.Table
	)
	schemaSlice, err = e.DBMetas()
	if err != nil || schemaSlice == nil || len(schemaSlice) <= 0 {
		return
	}
	schemaMap = make(map[string][]*schemas.Column, len(schemaSlice))
	for _, s := range schemaSlice {
		columns := make([]*schemas.Column, 0)
		for _, c := range s.Columns() {
			columns = append(columns, c)
		}
		schemaMap[s.Name] = columns
	}
	return
}
