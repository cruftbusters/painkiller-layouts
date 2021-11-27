package v1

import (
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

func TestAwaitingLayers(t *testing.T) {
	awaitingHeightmap := new(MockAwaitingLayerService)
	awaitingHillshade := new(MockAwaitingLayerService)
	controller := &AwaitingLayersController{awaitingHeightmap, awaitingHillshade}
	httpBaseURL, wsBaseURL := TestController(controller)
	client := ClientV2{BaseURL: httpBaseURL}
	instances := []struct {
		string
		*MockAwaitingLayerService
	}{
		{"/v1/awaiting_heightmap", awaitingHeightmap},
		{"/v1/awaiting_hillshade", awaitingHillshade},
	}

	t.Run("enqueue one", func(t *testing.T) {
		for _, instance := range instances {
			t.Run(instance.string, func(t *testing.T) {
				layout := types.Layout{Id: "enqueue me"}
				instance.MockAwaitingLayerService.On("Enqueue", layout).Return(nil)

				if err := client.EnqueueLayoutExpect(instance.string, layout, 201); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("enqueue one when queue is full", func(t *testing.T) {
		for _, instance := range instances {
			t.Run(instance.string, func(t *testing.T) {
				layout := types.Layout{Id: "im not gunna fit"}
				instance.MockAwaitingLayerService.On("Enqueue", layout).Return(ErrQueueFull).Once()

				if err := client.EnqueueLayoutExpect(instance.string, layout, 500); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("dequeue one", func(t *testing.T) {
		for _, instance := range instances {
			t.Run(instance.string, func(t *testing.T) {
				layout := types.Layout{Id: "rabid dequeueing"}
				instance.MockAwaitingLayerService.On("Dequeue", 0).Return(layout).Once()

				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+instance.string, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer conn.Close()

				got, err := BeginDequeueLayout(conn, 0)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, layout)
				conn.WriteMessage(websocket.BinaryMessage, nil)
			})
		}
	})

	t.Run("requeue work unfinished by closed workers", func(t *testing.T) {
		for _, instance := range instances {
			t.Run(instance.string, func(t *testing.T) {
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+instance.string, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer conn.Close()

				layout := types.Layout{Id: "requeue me"}
				instance.MockAwaitingLayerService.On("Dequeue", 0).Return(layout).Once()
				if _, err := BeginDequeueLayout(conn, 0); err != nil {
					t.Fatal(err)
				}

				channel := make(chan types.Layout)
				instance.MockAwaitingLayerService.On("Enqueue", mock.Anything).Return(nil).Run(func(args mock.Arguments) { channel <- args.Get(0).(types.Layout) }).Once()
				conn.Close()

				select {
				case got := <-channel:
					AssertLayout(t, got, layout)
				case <-time.After(time.Second):
					t.Fatal("timed out after one second")
				}
			})
		}
	})

	t.Run("dequeue more than one with one worker", func(t *testing.T) {
		for _, instance := range instances {
			t.Run(instance.string, func(t *testing.T) {
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+instance.string, nil)
				AssertNoError(t, err)
				defer conn.Close()

				first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}
				instance.MockAwaitingLayerService.On("Dequeue", 0).Return(first).Once()
				instance.MockAwaitingLayerService.On("Dequeue", 0).Return(second).Once()

				got, err := BeginDequeueLayout(conn, 0)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, first)
				if err := EndDequeueLayout(conn); err != nil {
					t.Fatal(err)
				}

				got, err = BeginDequeueLayout(conn, 0)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, second)
				if err := EndDequeueLayout(conn); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("workers specify priority", func(t *testing.T) {
		for _, instance := range instances {
			t.Run(instance.string, func(t *testing.T) {
				layout := types.Layout{Id: "prioritize that worker"}
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+instance.string+"?priority=123", nil)
				if err != nil {
					t.Fatal(err)
				}

				instance.MockAwaitingLayerService.On("Dequeue", 123).Return(layout).Once()

				gotChannel := make(chan struct {
					types.Layout
					error
				})
				go func() {
					got, err := BeginDequeueLayout(conn, 1)
					if err != nil {
						gotChannel <- struct {
							types.Layout
							error
						}{got, err}
						return
					}
					err = EndDequeueLayout(conn)
					gotChannel <- struct {
						types.Layout
						error
					}{got, err}
				}()
				select {
				case result := <-gotChannel:
					got, err := result.Layout, result.error
					if err != nil {
						t.Fatal(err)
					} else if got != layout {
						t.Fatalf("got %+v want %+v", got, layout)
					}
				case <-time.After(time.Second):
					t.Fatal("timed out after one second")
				}
			})
		}
	})
}
