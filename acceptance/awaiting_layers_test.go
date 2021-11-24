package acceptance

import (
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/gorilla/websocket"
)

func TestAwaitingLayers(t *testing.T) {
	httpBaseURL, wsBaseURL := TestServer(v1.Handler)
	client := &ClientV2{BaseURL: httpBaseURL}
	instances := []string{"/v1/awaiting_heightmap", "/v1/awaiting_hillshade"}

	t.Run("enqueue two and distribute", func(t *testing.T) {
		for _, path := range instances {
			t.Run("enqueue "+path, func(t *testing.T) {
				if err := client.EnqueueLayoutExpect(path, types.Layout{Id: path + "0"}, 201); err != nil {
					t.Fatal(err)
				}
				if err := client.EnqueueLayoutExpect(path, types.Layout{Id: path + "1"}, 201); err != nil {
					t.Fatal(err)
				}
			})
		}

		for i := len(instances) - 1; i >= 0; i-- {
			path := instances[i]
			t.Run("dequeue "+path, func(t *testing.T) {
				conn0, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer conn0.Close()

				conn1, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer conn1.Close()

				got, err := BeginDequeueLayout(conn0)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, types.Layout{Id: path + "0"})
				if err := EndDequeueLayout(conn0); err != nil {
					t.Fatal(err)
				}

				got, err = BeginDequeueLayout(conn1)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, types.Layout{Id: path + "1"})
				if err := EndDequeueLayout(conn1); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("requeue work unfinished by closed workers", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				AssertNoError(t, err)
				defer conn.Close()

				layout := types.Layout{Id: "bumpy ride"}
				if err := client.EnqueueLayoutExpect(path, layout, 201); err != nil {
					t.Fatal(err)
				}

				if _, err := BeginDequeueLayout(conn); err != nil {
					t.Fatal(err)
				}
				conn.Close()

				conn, _, err = websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				AssertNoError(t, err)
				defer conn.Close()

				got, err := BeginDequeueLayout(conn)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, layout)
				if err := EndDequeueLayout(conn); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("queue is full", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				queueSize := 8

				for i := 0; i < queueSize; i++ {
					if err := client.EnqueueLayoutExpect(path, types.Layout{}, 201); err != nil {
						t.Fatal(err)
					}
				}

				if err := client.EnqueueLayoutExpect(path, types.Layout{Id: "not gunna fit"}, 500); err != nil {
					t.Fatal(err)
				}

				for i := 0; i < queueSize; i++ {
					conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
					AssertNoError(t, err)
					defer conn.Close()

					if _, err = BeginDequeueLayout(conn); err != nil {
						t.Fatal(err)
					}
					EndDequeueLayout(conn)
				}
			})
		}
	})

	t.Run("dequeue more than one with one worker", func(t *testing.T) {
		for _, path := range instances {
			t.Run(path, func(t *testing.T) {
				conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+path, nil)
				AssertNoError(t, err)
				defer conn.Close()

				first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}
				if err := client.EnqueueLayoutExpect(path, first, 201); err != nil {
					t.Fatal(err)
				} else if err := client.EnqueueLayoutExpect(path, second, 201); err != nil {
					t.Fatal(err)
				}

				got, err := BeginDequeueLayout(conn)
				if err != nil {
					t.Fatal(err)
				}
				AssertLayout(t, got, first)
				if err := EndDequeueLayout(conn); err != nil {
					t.Fatal(err)
				}

				got, err = BeginDequeueLayout(conn)
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

	t.Run("new layout enqueues awaiting heightmap", func(t *testing.T) {
		created, err := client.CreateLayoutExpect(types.Layout{}, 201)
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteLayout(t, created.Id)

		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
		AssertNoError(t, err)
		defer conn.Close()

		got, err := BeginDequeueLayout(conn)
		if err != nil {
			t.Fatal(err)
		}
		AssertLayout(t, got, created)
		if err := EndDequeueLayout(conn); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("patch of empty heightmap url enqueues awaiting hillshade", func(t *testing.T) {
		created, err := client.CreateLayoutExpect(types.Layout{}, 201)
		if err != nil {
			t.Fatal(err)
		}
		defer client.DeleteLayout(t, created.Id)

		client.PatchLayout(t, created.Id, types.Layout{HillshadeURL: "irrelevant"})
		client.PatchLayout(t, created.Id, types.Layout{HeightmapURL: "time for hillshade"})
		client.PatchLayout(t, created.Id, types.Layout{HillshadeURL: "irrelevant2"})

		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_hillshade", nil)
		AssertNoError(t, err)
		defer conn.Close()

		got, err := BeginDequeueLayout(conn)
		if err != nil {
			t.Fatal(err)
		}
		AssertLayout(t, got, types.Layout{Id: created.Id, HeightmapURL: "time for hillshade", HillshadeURL: "irrelevant"})
		if err := EndDequeueLayout(conn); err != nil {
			t.Fatal(err)
		}
	})
}
