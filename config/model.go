package config

// https://yaml.to-go.online

type root struct {
	App     app     `yaml:"app"`
	Source  source  `yaml:"source"`
	Target  target  `yaml:"target"`
	Mapping mapping `yaml:"mapping"`
}

type app struct {
	PageSize    int `yaml:"page_size"`
	ThreadCount int `yaml:"thread_count"`
}

type database struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
}

type kafkaConsumer struct {
	Version      string   `yaml:"version"`
	Addresses    []string `yaml:"addresses"`
	Topic        string   `yaml:"topic"`
	Consumer     string   `yaml:"consumer"`
	DatabaseName string   `yaml:"databaseName"`
}

type kafkaProducer struct {
	Addresses []string `yaml:"addresses"`
	Topic     string   `yaml:"topic"`
}

type source struct {
	Database *database      `yaml:"database"`
	Kafka    *kafkaConsumer `yaml:"kafka"`
}

type target struct {
	Databases []*database    `yaml:"databases"`
	Kafka     *kafkaProducer `yaml:"kafka"`
}

type table struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type mapping struct {
	Tables map[string]string `yaml:"tables"`
}
