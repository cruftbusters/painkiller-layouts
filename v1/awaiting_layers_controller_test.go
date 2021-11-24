package v1

import (
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

type MockAwaitingLayerService struct {
	mock.Mock
}

func (m *MockAwaitingLayerService) Enqueue(got types.Layout) error {
	args := m.Called(got)
	return args.Error(0)
}

func (m *MockAwaitingLayerService) Dequeue() types.Layout {
	args := m.Called()
	return args.Get(0).(types.Layout)
}

func TestAwaitingLayers(t *testing.T) {
	awaitingHeightmap := new(MockAwaitingLayerService)
	controller := &AwaitingLayersController{awaitingHeightmap}
	httpBaseURL, wsBaseURL := TestController(controller)
	client := ClientV2{BaseURL: httpBaseURL}
	instances := []string{"/v1/awaiting_heightmap", "/v1/awaiting_hillshade"}

	t.Run("ping every five seconds", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				AssertNoError(t, err)
				defer conn.Close()
				go conn.ReadMessage()

				ping := make(chan *struct{})
				conn.SetPingHandler(func(string) error { ping <- nil; return nil })

				one, five, six := time.After(time.Second), time.After(5*time.Second), time.After(6*time.Second)
				select {
				case <-ping:
				case <-one:
					t.Fatal("timed out waiting for first ping")
				}

				select {
				case <-ping:
					t.Fatal("second ping too early")
				case <-five:
				}

				select {
				case <-ping:
				case <-six:
					t.Fatal("second ping too late")
				}
			})
		}
	})

	t.Run("enqueue one", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				layout := types.Layout{Id: "enqueue me"}
				awaitingHeightmap.On("Enqueue", layout).Return(nil)

				if err := client.EnqueueLayoutExpect(path, layout, 201); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("enqueue one when queue is full", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				layout := types.Layout{Id: "im not gunna fit"}
				awaitingHeightmap.On("Enqueue", layout).Return(ErrQueueFull)

				if err := client.EnqueueLayoutExpect(path, layout, 500); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("dequeue one", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				layout := types.Layout{Id: "rabid dequeueing"}
				awaitingHeightmap.On("Dequeue").Return(layout).Once()

				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer conn.Close()

				got, err := BeginDequeueLayout(conn)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, layout)
				conn.WriteMessage(websocket.BinaryMessage, nil)
			})
		}
	})

	t.Run("requeue work unfinished by closed workers", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer conn.Close()

				layout := types.Layout{Id: "requeue me"}
				awaitingHeightmap.On("Dequeue").Return(layout).Once()
				if _, err := BeginDequeueLayout(conn); err != nil {
					t.Fatal(err)
				}

				channel := make(chan types.Layout)
				awaitingHeightmap.On("Enqueue", mock.Anything).Return(nil).Run(func(args mock.Arguments) { channel <- args.Get(0).(types.Layout) }).Once()
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
}
