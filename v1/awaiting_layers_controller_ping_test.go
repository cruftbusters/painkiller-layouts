package v1

import (
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/gorilla/websocket"
)

func TestAwaitingControllerPing(t *testing.T) {
	awaitingHeightmap := new(MockAwaitingLayerService)
	awaitingHillshade := new(MockAwaitingLayerService)
	controller := &AwaitingLayersController{awaitingHeightmap, awaitingHillshade}
	_, wsBaseURL := TestController(controller)
	instances := []struct {
		string
		*MockAwaitingLayerService
	}{
		{"/v1/awaiting_heightmap", awaitingHeightmap},
		{"/v1/awaiting_hillshade", awaitingHillshade},
	}

	t.Run("ping every five seconds", func(t *testing.T) {
		for _, instance := range instances {
			t.Run(instance.string, func(t *testing.T) {
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+instance.string, nil)
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
}
