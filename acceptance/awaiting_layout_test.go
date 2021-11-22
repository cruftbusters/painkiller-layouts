package acceptance

import (
	"fmt"
	"sync"
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
)

func TestPendingRenders(t *testing.T) {
	httpBaseURL, wsBaseURL := t2.TestServer(v1.Handler)
	client := t2.ClientV2{BaseURL: httpBaseURL}

	t.Run("enqueue layout awaiting heightmap when queue full", func(t *testing.T) {
		queueSize := 2

		for i := 0; i < queueSize; i++ {
			client.EnqueueLayout(t, types.Layout{})
		}
		client.EnqueueLayoutExpectInternalServerError(t, types.Layout{})

		wsClient, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()
		for i := 0; i < queueSize; i++ {
			_, err := wsClient.StartDequeue()
			t2.AssertNoError(t, err)
			wsClient.EndDequeue()
		}
	})

	t.Run("ping every interval", func(t *testing.T) {
		wsClient, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		go wsClient.Conn.ReadMessage()

		signal := make(chan *struct{})
		wsClient.Conn.SetPingHandler(func(string) error { signal <- nil; return nil })

		one, five, six := time.After(time.Second), time.After(5*time.Second), time.After(6*time.Second)
		select {
		case <-signal:
		case <-one:
			t.Fatal("expected ping in less than one second")
		}

		select {
		case <-signal:
			t.Fatal("expected no ping in less than five seconds")
		case <-five:
		}

		select {
		case <-signal:
		case <-six:
			t.Fatal("expected ping in less than six seconds")
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
			wsClient, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
			t2.AssertNoError(t, err)
			defer wsClient.Conn.Close()

			go func(want types.Layout) {
				t.Run(want.Id, func(t *testing.T) {
					got, err := wsClient.StartDequeue()
					if err != nil {
						t.Errorf("got %s want nil", err)
					} else if got != want {
						t.Errorf("got %+v want %+v", got, want)
					}
					wsClient.EndDequeue()
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

		wsClient, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		got, err := wsClient.StartDequeue()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, layout)
		wsClient.EndDequeue()
	})

	t.Run("pull multiple awaiting layouts with one worker", func(t *testing.T) {
		wsClient, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}

		client.EnqueueLayout(t, first)
		got, err := wsClient.StartDequeue()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, first)
		wsClient.EndDequeue()

		client.EnqueueLayout(t, second)
		got, err = wsClient.StartDequeue()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
		wsClient.EndDequeue()
	})

	t.Run("re-dispatch abandoned awaiting layout", func(t *testing.T) {
		for i := 0; i < 16; i++ {
			wsClient, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
			t2.AssertNoError(t, err)
			wsClient.Conn.Close()
		}

		wsClient, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()

		first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}

		client.EnqueueLayout(t, first)
		got, err := wsClient.StartDequeue()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, first)
		wsClient.EndDequeue()

		client.EnqueueLayout(t, second)
		got, err = wsClient.StartDequeue()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
		wsClient.Conn.Close()

		wsClient, err = t2.LayoutsAwaitingHeightmap(wsBaseURL)
		t2.AssertNoError(t, err)
		defer wsClient.Conn.Close()
		got, err = wsClient.StartDequeue()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
		wsClient.EndDequeue()
	})

	t.Run("enqueue new layouts", func(t *testing.T) {
		created := client.CreateLayout(t, types.Layout{})
		defer client.DeleteLayout(t, created.Id)

		awaiting, err := t2.LayoutsAwaitingHeightmap(wsBaseURL)
		t2.AssertNoError(t, err)
		got, err := awaiting.StartDequeue()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, created)
		awaiting.EndDequeue()
		awaiting.Conn.Close()

		t.Run("enqueue and dequeue layouts awaiting hillshade", func(t *testing.T) {
			patched := client.PatchLayout(t, created.Id, types.Layout{HeightmapURL: "heightmap URL"})

			awaiting, err := t2.LayoutsAwaitingHillshade(wsBaseURL)
			t2.AssertNoError(t, err)
			got, err := awaiting.StartDequeue()
			t2.AssertNoError(t, err)
			t2.AssertLayout(t, got, patched)
			awaiting.EndDequeue()
		})
	})
}
