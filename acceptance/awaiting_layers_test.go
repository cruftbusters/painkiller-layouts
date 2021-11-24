package acceptance

import (
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/gorilla/websocket"
)

func TestAwaitingLayers(t *testing.T) {
	httpBaseURL, wsBaseURL := TestServer(v1.Handler)
	client := &ClientV2{BaseURL: httpBaseURL}

	t.Run("ping every five seconds", func(t *testing.T) {
		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
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

	t.Run("enqueue two and distribute", func(t *testing.T) {
		layout0 := types.Layout{Id: "0"}
		if err := client.EnqueueLayoutAwaitingHeightmap(layout0); err != nil {
			t.Fatal(err)
		}

		layout1 := types.Layout{Id: "1"}
		if err := client.EnqueueLayoutAwaitingHeightmap(layout1); err != nil {
			t.Fatal(err)
		}

		conn0, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn0.Close()

		conn1, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn1.Close()

		conn0.WriteMessage(websocket.BinaryMessage, nil)
		got, err := ReadLayout(conn0)
		if err != nil {
			t.Fatal(err)
		}
		AssertLayout(t, got, layout0)
		conn0.WriteMessage(websocket.BinaryMessage, nil)

		conn1.WriteMessage(websocket.BinaryMessage, nil)
		got, err = ReadLayout(conn1)
		if err != nil {
			t.Fatal(err)
		}
		AssertLayout(t, got, layout1)
		conn1.WriteMessage(websocket.BinaryMessage, nil)
	})

	t.Run("requeue work unfinished by closed workers", func(t *testing.T) {
		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
		AssertNoError(t, err)
		defer conn.Close()

		layout := types.Layout{Id: "bumpy ride"}
		if err := client.EnqueueLayoutAwaitingHeightmap(layout); err != nil {
			t.Fatal(err)
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, nil); err != nil {
			t.Fatal(err)
		} else if _, err := ReadLayout(conn); err != nil {
			t.Fatal(err)
		}
		conn.Close()

		conn, _, err = websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
		AssertNoError(t, err)
		defer conn.Close()

		if err := conn.WriteMessage(websocket.BinaryMessage, nil); err != nil {
			t.Fatal(err)
		}
		got, err := ReadLayout(conn)
		if err != nil {
			t.Fatal(err)
		}
		AssertLayout(t, got, layout)
		conn.WriteMessage(websocket.BinaryMessage, nil)
	})

	t.Run("queue is full", func(t *testing.T) {
		queueSize := 2

		for i := 0; i < queueSize; i++ {
			if err := client.EnqueueLayoutAwaitingHeightmap(types.Layout{}); err != nil {
				t.Fatal(err)
			}
		}

		if err := client.EnqueueLayoutAwaitingHeightmapExpectInternalServerError(types.Layout{Id: "not gunna fit"}); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < queueSize; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
			AssertNoError(t, err)
			defer conn.Close()

			conn.WriteMessage(websocket.BinaryMessage, nil)
			if _, err = ReadLayout(conn); err != nil {
				t.Fatal(err)
			}
			conn.WriteMessage(websocket.BinaryMessage, nil)
		}
	})
}
