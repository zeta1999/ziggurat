package kstream

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gojekfarm/ziggurat-go/pkg/z"
	"github.com/gojekfarm/ziggurat-go/pkg/zlogger"
	"github.com/gojekfarm/ziggurat-go/pkg/zmw"
	"sync"
	"time"
)

const defaultPollTimeout = 100 * time.Millisecond
const brokerRetryTimeout = 2 * time.Second

var startConsumer = func(app z.App, h z.MessageHandler, consumer *kafka.Consumer, route string, instanceID string, wg *sync.WaitGroup) {
	go func(instanceID string) {
		doneCh := app.Context().Done()
		worker := NewWorker(10)
		sendCh, _ := worker.run(app, func(message *kafka.Message) {
			processor(message, route, consumer, h, app)
		})
		for {
			select {
			case <-doneCh:
				close(sendCh)
				wg.Done()
				return
			default:
				msg, err := readMessage(consumer, defaultPollTimeout)
				if err != nil && err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				} else if err != nil && err.(kafka.Error).Code() == kafka.ErrAllBrokersDown {
					zlogger.LogError(err, "retrying broker...", nil)
					time.Sleep(brokerRetryTimeout)
					continue
				}
				if msg != nil {
					sendCh <- msg
				}
			}
		}
	}(instanceID)
}

var StartConsumers = func(app z.App, consumerConfig *kafka.ConfigMap, topicEntity string, topics []string, instances int, h z.MessageHandler, wg *sync.WaitGroup) []*kafka.Consumer {
	consumers := make([]*kafka.Consumer, 0, instances)
	for i := 0; i < instances; i++ {
		consumer := createConsumer(consumerConfig, topics)
		consumers = append(consumers, consumer)
		groupID, _ := consumerConfig.Get("group.id", "")
		instanceID := fmt.Sprintf("%s_%s_%d", topicEntity, groupID, i)
		wg.Add(1)
		startConsumer(app, zmw.DefaultTerminalMW(h), consumer, topicEntity, instanceID, wg)
	}
	return consumers
}
