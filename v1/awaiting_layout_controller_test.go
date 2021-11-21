package v1

import (
	"fmt"
	"sync"
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
)

func TestAwaitingLayoutController(t *testing.T) {
	t.Run("gracefully close connection", func(t *testing.T) {
		controller := &AwaitingLayoutController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		for i := 0; i < 16; i++ {
			wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
			t2.AssertNoError(t, err)
			wsClient.Conn.Close()
		}

		client.EnqueueLayout(t, types.Layout{})
	})

	t.Run("ping every interval", func(t *testing.T) {
		interval := time.Second
		controller := &AwaitingLayoutController{interval}
		_, wsBaseURL := t2.TestController(controller)
		wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
		t2.AssertNoError(t, err)

		go wsClient.Conn.ReadMessage()

		signal := make(chan *struct{})
		wsClient.Conn.SetPingHandler(func(string) error { signal <- nil; return nil })

		one, five, six := time.After(time.Second), time.After(interval), time.After(interval+time.Second)
		select {
		case <-signal:
		case <-one:
			t.Fatal("expected ping in less than one second")
		}

		select {
		case <-signal:
			t.Fatalf("expected no ping in less than %s", interval)
		case <-five:
		}

		select {
		case <-signal:
		case <-six:
			t.Fatalf("expected ping in less than %s", interval+time.Second)
		}
	})

	t.Run("dispatch one awaiting layout", func(t *testing.T) {
		controller := &AwaitingLayoutController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		layouts := [2]types.Layout{}
		for i := 0; i < len(layouts); i++ {
			layouts[i] = types.Layout{Id: fmt.Sprintf("layout #%d", i)}
			client.EnqueueLayout(t, layouts[i])
		}

		var wg sync.WaitGroup
		wg.Add(len(layouts))
		for _, want := range layouts {
			wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
			t2.AssertNoError(t, err)
			defer wsClient.Conn.Close()

			go func(want types.Layout) {
				got, err := wsClient.StartDequeueAwaitingLayout()
				if err != nil {
					t.Errorf("got %s want nil", err)
				} else if got != want {
					t.Errorf("got %+v want %+v", got, want)
				}
				wg.Done()
			}(want)
		}
		wg.Wait()
	})

	t.Run("buffer awaiting layouts", func(t *testing.T) {
		controller := &AwaitingLayoutController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		layout := types.Layout{Id: "unhandled"}
		client.EnqueueLayout(t, layout)

		wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		got, err := wsClient.StartDequeueAwaitingLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, layout)
	})

	t.Run("pull multiple awaiting layouts with one worker", func(t *testing.T) {
		controller := &AwaitingLayoutController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}

		client.EnqueueLayout(t, first)
		got, err := wsClient.StartDequeueAwaitingLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, first)
		wsClient.CompleteDequeueAwaitingLayout()

		client.EnqueueLayout(t, second)
		got, err = wsClient.StartDequeueAwaitingLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
	})

	t.Run("re-dispatch abandoned awaiting layout", func(t *testing.T) {
		controller := &AwaitingLayoutController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		for i := 0; i < 16; i++ {
			wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
			t2.AssertNoError(t, err)
			wsClient.Conn.Close()
		}

		wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}

		client.EnqueueLayout(t, first)
		got, err := wsClient.StartDequeueAwaitingLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, first)
		wsClient.CompleteDequeueAwaitingLayout()

		client.EnqueueLayout(t, second)
		got, err = wsClient.StartDequeueAwaitingLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
		wsClient.Conn.Close()

		wsClient, err = t2.NewLayoutsAwaitingClient(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()
		got, err = wsClient.StartDequeueAwaitingLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
	})

	t.Run("overflow", func(t *testing.T) {
		controller := &AwaitingLayoutController{time.Second}
		httpBaseURL, _ := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		limit := 2

		for i := 0; i < limit; i++ {
			client.EnqueueLayout(t, types.Layout{})
		}
		client.EnqueueLayoutExpectInternalServerError(t, types.Layout{})
	})
}
