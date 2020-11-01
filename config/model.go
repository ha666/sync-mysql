package config

// https://yaml.to-go.online

type root struct {
	App    app    `yaml:"app"`
	Source source `yaml:"source"`
	Target target `yaml:"target"`
}

type app struct {
	PageSize    int    `yaml:"page_size"`
	ThreadCount uint64 `yaml:"thread_count"`
}

type database struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
}

type kafka struct {
	Version   string   `yaml:"version"`
	Addresses []string `yaml:"addresses"`
	Topic     string   `yaml:"topic"`
	Consumer  string   `yaml:"consumer"`
}

type source struct {
	Database *database `yaml:"database"`
	Kafka    *kafka    `yaml:"kafka"`
}

type target struct {
	Databases []*database `yaml:"databases"`
	Kafka     *kafka      `yaml:"kafka"`
}
