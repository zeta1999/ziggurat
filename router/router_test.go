package router

import (
	"github.com/gojekfarm/ziggurat"
	"reflect"
	"testing"
)

func TestDefaultRouter_HandleMessage(t *testing.T) {
	dr := New()
	expectedEvent := ziggurat.CreateMockEvent()
	expectedEvent.ValueFunc = func() []byte {
		return nil
	}
	expectedEvent.HeadersFunc = func() map[string]string {
		return map[string]string{ziggurat.HeaderMessageRoute: "bar"}
	}
	dr.HandleFunc("foo", func(event ziggurat.Event) error {
		if !reflect.DeepEqual(event, expectedEvent) {
			t.Errorf("expected event %+v, got %+v", expectedEvent, event)
		}
		return nil
	})
	dr.HandleEvent(ziggurat.MockEvent{
		ValueFunc: func() []byte {
			return nil
		},
		HeadersFunc: func() map[string]string {
			return map[string]string{ziggurat.HeaderMessageRoute: "bar"}
		},
	})
}

func TestDefaultRouter_NotFoundHandler(t *testing.T) {
	notFoundHandlerCalled := false
	dr := New(WithNotFoundHandler(func(event ziggurat.Event) error {
		notFoundHandlerCalled = true
		return nil
	}))

	dr.HandleFunc("foo", func(event ziggurat.Event) error {
		return nil
	})

	dr.HandleEvent(ziggurat.MockEvent{
		ValueFunc: func() []byte {
			return []byte{}
		},
		HeadersFunc: func() map[string]string {
			return map[string]string{ziggurat.HeaderMessageRoute: "bar"}
		},
	})

	if !notFoundHandlerCalled {
		t.Errorf("expected %v got %v", true, notFoundHandlerCalled)
	}

}
