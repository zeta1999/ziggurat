package ziggurat

import (
	"context"
	"testing"
)

func TestPipeHandlers(t *testing.T) {
	mw1 := func(next Handler) Handler {
		return HandlerFunc(func(messageEvent Event, ) ProcessStatus {
			me := CreateMessageEvent(nil, []byte("foo"), nil, context.Background())
			return next.HandleMessage(me)
		})
	}
	mw2 := func(next Handler) Handler {
		return HandlerFunc(func(messageEvent Event) ProcessStatus {
			byteValue := append(messageEvent.Value(), []byte("-bar")...)
			me := CreateMessageEvent(nil, byteValue, nil, context.Background())
			return next.HandleMessage(me)
		})
	}
	actualHandler := HandlerFunc(func(event Event) ProcessStatus {
		if string(event.Value()) != "foo-bar" {
			t.Errorf("expected message to be %s,but got %s", "foo-bar", string(event.Value()))
		}
		return ProcessingSuccess
	})
	finalHandler := PipeHandlers(mw1, mw2)(actualHandler)
	finalHandler.HandleMessage(Message{})
}
