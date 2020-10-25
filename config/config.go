package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var (
	Conf *root
)

//根
type root struct {

	//数据库
	DataBases databases `yaml:"databases"`
}

type databases struct {

	//源数据库
	Source database `yaml:"source"`

	//目标数据库
	Target database `yaml:"target"`
}

//数据库
type database struct {

	//数据库名称
	Name string `yaml:"name"`

	//数据库实例地址
	Address string `yaml:"address"`

	//数据库实例端口号
	Port int `yaml:"port"`

	//数据库帐号
	Account string `yaml:"account"`

	//数据库密码
	Password string `yaml:"password"`
}

//解析方法
func Parser() error {
	Conf = new(root)
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		return err
	}
	return nil
}
