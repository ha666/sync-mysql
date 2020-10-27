package service

func insertData(dbId int, args []interface{}) (err error) {
	_, err = targetEngines[dbId].Exec(args...)
	return err
}
