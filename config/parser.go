package config

import (
	"io/ioutil"
	"os"

	"github.com/ha666/golibs"
	"github.com/ha666/logs"
	"gopkg.in/yaml.v3"
)

var (
	defaultConfigFile = "./config.yaml"
	Conf              *root
)

//解析方法
func Parser() error {
	Conf = new(root)
	var (
		err      error
		yamlFile []byte
	)
	configFile := os.Getenv("sync_mysql_config_file")
	if golibs.Length(configFile) <= 0 {
		configFile = defaultConfigFile
	}
	logs.Info("加载配置文件:%s", configFile)
	yamlFile, err = ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		return err
	}
	return nil
}
