package acceptance

import (
	"sync"
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/gorilla/websocket"
)

func TestPingPong(t *testing.T) {
	httpBaseURL, wsBaseURL := TestServer(v1.Handler)
	instances := []string{"/v1/awaiting_heightmap", "/v1/awaiting_hillshade"}

	client := ClientV2{BaseURL: httpBaseURL}
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

	t.Run("respond to ping with pong", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				client.EnqueueLayoutExpect(path, types.Layout{}, 201)

				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer conn.Close()

				var wg sync.WaitGroup
				wg.Add(2)

				channel := make(chan error)
				go func() {
					for {
						_, _, err := conn.ReadMessage()
						channel <- err
						if err != nil {
							close(channel)
							return
						}
					}
				}()
				go func() {
					defer wg.Done()
					for {
						if err := conn.WriteMessage(websocket.BinaryMessage, nil); err != nil {
							return
						} else if err := <-channel; err != nil {
							return
						} else if err := conn.WriteMessage(websocket.BinaryMessage, nil); err != nil {
							return
						}
					}
				}()

				go func() {
					defer wg.Done()
					defer conn.Close()
					channel := make(chan string)
					conn.SetPongHandler(func(s string) error { channel <- s; return nil })
					for i := 0; i < 2; i++ {
						if err := conn.WriteControl(websocket.PingMessage, nil, time.Time{}); err != nil {
							t.Error(err)
							return
						}
						select {
						case <-channel:
							time.Sleep(1 * time.Second)
						case <-time.After(1 * time.Second):
							t.Error("waiting for pong timed out after one seconds")
							return
						}
					}
				}()
				wg.Wait()
			})
		}
	})
}
