package config

//根
type root struct {

	//应用配置
	App app `yaml:"app"`

	//来源配置
	Source source `yaml:"source"`

	//目标配置
	Target target `yaml:"target"`
}

//应用配置
type app struct {

	//每页记录数
	PageSize int `yaml:"page_size"`
}

//来源配置
type source struct {

	//来源类型,可选值:db、kafka
	Type string `yaml:"type"`

	//数据库配置
	DataBase database `yaml:"database"`
}

//目标配置
type target struct {

	//目标类型,可选值:db、kafka
	Type string `yaml:"type"`

	//数据库配置
	DataBase database `yaml:"database"`
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
