package zig

import (
	"sync"
	"time"
)

type MessageEvent struct {
	MessageKey        interface{}
	MessageValue      interface{}
	MessageValueBytes []byte
	MessageKeyBytes   []byte
	Topic             string
	TopicEntity       string
	KafkaTimestamp    time.Time
	TimestampType     string
	Attributes        map[string]interface{}
}

func (m MessageEvent) GetMessageAttribute(key string) interface{} {
	return m.Attributes[key]
}

func (m *MessageEvent) SetMessageAttribute(key string, value interface{}) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	m.Attributes[key] = value
}
