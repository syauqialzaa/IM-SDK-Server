package kafka

import (
	"gin-chat-svc/pkg/logger"
	"strings"

	"github.com/Shopify/sarama"
)

var consumer sarama.Consumer

type ConsumerCallback func(data []byte)

// initialize the consumer
func InitConsumer(hosts string) {
	config := sarama.NewConfig()
	client, err := sarama.NewClient(strings.Split(hosts, ","), config)
	if nil != err {
		logger.Logger.Error("init kafka consumer client error", logger.Any("init kafka consumer client error", err.Error()))
	}

	consumer, err = sarama.NewConsumerFromClient(client)
	if nil != err {
		logger.Logger.Error("init kafka consumer error", logger.Any("init kafka consumer error", err.Error()))
	}
}

// consume messages through callback functions
func ConsumerMsg(callBack ConsumerCallback) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if nil != err {
		logger.Logger.Error("ConsumePartition error", logger.Any("ConsumePartition error", err.Error()))
		return
	}

	defer partitionConsumer.Close()
	for {
		msg := <- partitionConsumer.Messages()
		if nil != callBack {
			callBack(msg.Value)
		}
	}
}

func CloseConsumer() {
	if nil != consumer {
		consumer.Close()
	}
}