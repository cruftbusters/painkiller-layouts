package v1

import (
	"sync"
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
)

func TestAwaitingLayerService(t *testing.T) {
	queueSize := 2
	service := NewAwaitingLayerService(queueSize)

	t.Run("enqueue and dequeue", func(t *testing.T) {
		layout := types.Layout{Id: "enqueue this"}
		if err := service.Enqueue(layout); err != nil {
			t.Fatal(err)
		}
		got := service.Dequeue(0)
		AssertLayout(t, got, layout)
	})

	t.Run("queue is full", func(t *testing.T) {
		for i := 0; i < queueSize; i++ {
			if err := service.Enqueue(types.Layout{}); err != nil {
				t.Fatal(err)
			}
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
		for i := 0; i < queueSize; i++ {
			go func() { service.Dequeue(0); c1 <- nil }()
			select {
			case <-c1:
			case <-time.After(time.Second):
				t.Fatal("timed out after one second")
			}
		}
	})

	t.Run("workers specify priority", func(t *testing.T) {
		layout0, layout1 := types.Layout{Id: "layout0"}, types.Layout{Id: "layout1"}

		var done sync.WaitGroup
		done.Add(3)
		go func() {
			defer done.Done()
			got1 := service.Dequeue(1)
			if got1 != layout1 {
				t.Errorf("got %+v want %+v", got1, layout1)
			}
		}()
		go func() {
			defer done.Done()
			got0 := service.Dequeue(0)
			if got0 != layout0 {
				t.Errorf("got %+v want %+v", got0, layout0)
			}
		}()
		go func() {
			defer done.Done()
			time.Sleep(125 * time.Millisecond)
			if err := service.Enqueue(layout0); err != nil {
				t.Error(err)
			} else if err := service.Enqueue(layout1); err != nil {
				t.Error(err)
			}
		}()
		done.Wait()
	})
}
