package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
)

func initKPub(name string, brokers []string) error {
	var err error
	kc := sarama.NewConfig()
	mutex.Lock()
	defer mutex.Unlock()
	if kafkaPubs[name] != nil {
		return nil
	}
	kc.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	kafkaPub, err := sarama.NewSyncProducer(brokers, kc)
	if err != nil {
		return err
	}
	kafkaPubs[name] = kafkaPub
	return nil
}

func getKPub(name string) (sarama.SyncProducer, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	if p, ok := kafkaPubs[name]; ok {
		return p, nil
	}
	return nil, errors.New("must first init")
}
