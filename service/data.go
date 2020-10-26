package service

func insertData(args []interface{}) (err error) {
	_, err = targetEngine.Exec(args...)
	return err
}
