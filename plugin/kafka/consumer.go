package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ha666/logs"
	"time"
)

type Msg struct {
	Key   string
	Topic string
	Data  interface{}
}

func SendKMsg(name string, msg *Msg) error {
	b, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(msg.Key),
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(b),
	}
	pub, err := getKPub(name)
	if err != nil {
		return err
	}
	_, _, err = pub.SendMessage(m)
	return err
}

func initKSub(name string, opt *kafkaOpts) error {
	var err error
	version, err := sarama.ParseKafkaVersion(opt.Version)
	if err != nil {
		return fmt.Errorf("Error parsing Kafka version: %v", err)
	}
	config := sarama.NewConfig()
	config.Version = version
	switch opt.Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		return fmt.Errorf("Unrecognized consumer group partition assignor: %s", opt.Assignor)
	}
	if opt.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	client, err := sarama.NewConsumerGroup(opt.Brokers, opt.Group, config)
	if err != nil {
		return fmt.Errorf("Error creating consumer group client: %v", err)
	}
	kafkaConsumers[getKey(name, opt.Topic)] = client
	return nil
}

func getKSub(name string, topic string) (sarama.ConsumerGroup, error) {
	if v := kafkaConsumers[getKey(name, topic)]; v != nil {
		return v, nil
	}
	return nil, errors.New("must first init")
}

func KConsume(ctx context.Context, name, topic string, fn func([]byte) error) error {
	c, err := getKSub(name, topic)
	if err != nil {
		return err
	}
	topics := []string{topic}
	for {
		if err := c.Consume(ctx, topics, &kConsumer{
			Handle: fn,
		}); err != nil {
			logs.Error("Error from consumer: %v", err)
			time.Sleep(time.Second)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return nil
		}
	}
}

// Consumer represents a Sarama consumer group consumer
type kConsumer struct {
	Handle func([]byte) error
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *kConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	logs.Info("Sarama consumer up and running!...")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *kConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *kConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := c.Handle(message.Value); err != nil {
			logs.Info("Message claimed: value = %s, timestamp = %v, topic = %s, err = %v", string(message.Value), message.Timestamp, message.Topic, err)
			return err
		}
		session.MarkMessage(message, "")
	}
	return nil
}
