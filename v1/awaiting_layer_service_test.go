package v1

import (
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
)

func TestAwaitingLayerService(t *testing.T) {
	service := NewAwaitingLayerService(1)

	t.Run("enqueue and dequeue", func(t *testing.T) {
		layout := types.Layout{Id: "enqueue this"}
		if err := service.Enqueue(layout); err != nil {
			t.Fatal(err)
		}
		got := service.Dequeue(0)
		AssertLayout(t, got, layout)
	})

	t.Run("queue is full", func(t *testing.T) {
		if err := service.Enqueue(types.Layout{}); err != nil {
			t.Fatal(err)
		}

		c0 := make(chan error)
		go func() { c0 <- service.Enqueue(types.Layout{}) }()
		select {
		case err := <-c0:
			if err != ErrQueueFull {
				t.Fatalf("expected %s", ErrQueueFull)
			}
		case <-time.After(time.Second):
			t.Fatal("timed out after one second")
		}

		c1 := make(chan *struct{})
		go func() { service.Dequeue(0); c1 <- nil }()
		select {
		case <-c1:
		case <-time.After(time.Second):
			t.Fatal("timed out after one second")
		}
	})
}

func gochan(f func()) (channel chan *struct{}) {
	go func() { f(); channel <- nil }()
	return
}
