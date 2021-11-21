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
	interval := time.Second
	awaitingLayouts := make(chan types.Layout, 2)
	controller := &AwaitingLayoutController{interval, awaitingLayouts}
	httpBaseURL, wsBaseURL := t2.TestController(controller)
	client := t2.ClientV2{BaseURL: httpBaseURL}

	t.Run("ping every interval", func(t *testing.T) {
		wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

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
				t.Run(want.Id, func(t *testing.T) {
					got, err := wsClient.StartDequeueAwaitingLayout()
					if err != nil {
						t.Errorf("got %s want nil", err)
					} else if got != want {
						t.Errorf("got %+v want %+v", got, want)
					}
					wsClient.CompleteDequeueAwaitingLayout()
					wsClient.Conn.Close()
					wg.Done()
				})
			}(want)
		}
		wg.Wait()
	})

	t.Run("buffer awaiting layouts", func(t *testing.T) {
		layout := types.Layout{Id: "unhandled"}
		client.EnqueueLayout(t, layout)

		wsClient, err := t2.NewLayoutsAwaitingClient(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		got, err := wsClient.StartDequeueAwaitingLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, layout)
		wsClient.CompleteDequeueAwaitingLayout()
	})

	t.Run("pull multiple awaiting layouts with one worker", func(t *testing.T) {
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
		wsClient.CompleteDequeueAwaitingLayout()
	})

	t.Run("re-dispatch abandoned awaiting layout", func(t *testing.T) {
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
		wsClient.CompleteDequeueAwaitingLayout()
	})

	t.Run("overflow", func(t *testing.T) {
		limit := 2

		for i := 0; i < limit; i++ {
			client.EnqueueLayout(t, types.Layout{})
		}
		client.EnqueueLayoutExpectInternalServerError(t, types.Layout{})
	})
}
