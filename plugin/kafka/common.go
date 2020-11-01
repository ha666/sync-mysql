package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ha666/golibs"
	"github.com/ha666/logs"
	"sync"
)

var (
	mutex          sync.RWMutex
	kafkaPubs      = make(map[string]sarama.SyncProducer)
	kafkaConsumers = make(map[string]sarama.ConsumerGroup)
)

type kafkaOpts struct {
	Brokers  []string
	Group    string
	Version  string
	Topic    string
	Assignor string
	Oldest   bool
	Verbose  bool
}

func getKey(name string, key string) string {
	if name == "" {
		name = "default"
	}
	return fmt.Sprintf("%s-%v", name, key)
}

func newOpts(kafkaVersion string, fns ...func(*kafkaOpts)) *kafkaOpts {
	logs.Info("【newOpts】kafka.version：%s", kafkaVersion)
	opt := kafkaOpts{
		Version:  kafkaVersion,
		Assignor: "range",
		Oldest:   true,
		Verbose:  false,
	}
	for _, fn := range fns {
		fn(&opt)
	}

	if len(opt.Brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}
	if golibs.Length(opt.Topic) == 0 {
		panic("no topics given to be consumed, please set the -topic flag")
	}
	if len(opt.Group) == 0 {
		panic("no Kafka consumer group defined, please set the -group flag")
	}
	return &opt
}

//初始化Kafka生产者
func InitProducer(name string, addresses []string) {
	if golibs.Length(name) <= 0 {
		name = "default"
	}
	err := initKPub(name, addresses)
	if err != nil {
		logs.Emergency("【InitConsumer】err:%s", err.Error())
	}
	logs.Info("【InitConsumer】初始化kafka生产者%s成功", name)
}

//初始化Kafka消费者
func InitConsumer(name string, addresses []string, topic, kafkaVersion, consumerName string) {
	if golibs.Length(name) <= 0 {
		name = "default"
	}
	if err := initKSub(name, newOpts(kafkaVersion, func(opt *kafkaOpts) {
		opt.Topic = topic
		opt.Brokers = addresses
		opt.Group = consumerName
	})); err != nil {
		logs.Emergency("【InitConsumer】err:%s", err.Error())
	}
	logs.Info("【InitConsumer】初始化kafka消费者%s,%s成功", name, consumerName)
}
