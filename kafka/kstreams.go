package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gojekfarm/ziggurat"
	"github.com/gojekfarm/ziggurat/logger"
	"strings"
	"sync"
)

type ConsumerConfig struct {
	BootstrapServers string
	OriginTopics     string
	ConsumerGroupID  string
	ConsumerCount    int
}

type RouteGroup map[string]ConsumerConfig

type Streams struct {
	routeConsumerMap map[string][]*kafka.Consumer
	Logger           ziggurat.StructuredLogger
	RouteGroup
}

func (k *Streams) Stream(ctx context.Context, handler ziggurat.Handler) chan error {
	if k.Logger == nil {
		k.Logger = logger.NewJSONLogger("info")
	}
	var wg sync.WaitGroup
	k.routeConsumerMap = make(map[string][]*kafka.Consumer, len(k.RouteGroup))
	stopChan := make(chan error)
	for routeName, stream := range k.RouteGroup {
		consumerConfig := NewConsumerConfig(stream.BootstrapServers, stream.ConsumerGroupID)
		topics := strings.Split(stream.OriginTopics, ",")
		k.routeConsumerMap[routeName] = StartConsumers(ctx, consumerConfig, routeName, topics, stream.ConsumerCount, handler, k.Logger, &wg)
	}

	go func() {
		wg.Wait()
		k.stop()
		stopChan <- nil
	}()

	return stopChan
}

func (k *Streams) stop() {
	for _, consumers := range k.routeConsumerMap {
		for i, _ := range consumers {
			k.Logger.Error("error stopping consumer %v", consumers[i].Close(), nil)
		}
	}
}
