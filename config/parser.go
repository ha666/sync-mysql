package config

import (
	"io/ioutil"
	"os"

	"github.com/ha666/golibs"
	"gopkg.in/yaml.v3"
)

var (
	Conf *root
)

//解析方法
func Parser() error {
	Conf = new(root)
	var (
		err      error
		yamlFile []byte
	)
	configFile := os.Getenv("sync_mysql_config_file")
	if golibs.Length(configFile) > 0 {
		yamlFile, err = ioutil.ReadFile(configFile)
	} else {
		yamlFile, err = ioutil.ReadFile("./config.yaml")
	}
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		return err
	}
	return nil
}
