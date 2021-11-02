package heightmap

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubService struct {
	t                  testing.TB
	whenGetCalledWith  string
	getWillReturn      *Metadata
	whenPostCalledWith Metadata
	postWillReturn     Metadata
}

func (stub *StubService) get(got string) *Metadata {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.getWillReturn
}

func (stub *StubService) post(got Metadata) Metadata {
	stub.t.Helper()
	want := stub.whenPostCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.postWillReturn
}

func TestController(t *testing.T) {
	listener, port := RandomPortListener()
	client := NewClient(t, fmt.Sprintf("http://localhost:%d", port))

	stubService := &StubService{t: t}
	controller := Controller{
		stubService,
	}

	go func() {
		http.Serve(listener, controller)
	}()

	t.Run("get missing heightmap", func(t *testing.T) {
		stubService.whenGetCalledWith = "deadbeef"
		stubService.getWillReturn = nil

		client.GetExpectNotFound()
	})

	t.Run("create heightmap", func(t *testing.T) {
		up, down := Metadata{Id: "up"}, Metadata{Id: "down"}
		stubService.whenPostCalledWith = up
		stubService.postWillReturn = down

		got := client.Create(up)
		AssertMetadata(t, got, down)
	})

	t.Run("get heightmap", func(t *testing.T) {
		stubService.whenGetCalledWith = "path-id"
		stubService.getWillReturn = &Metadata{Id: "beefdead"}

		got := client.Get("path-id")
		want := Metadata{Id: "beefdead"}
		AssertMetadata(t, got, want)
	})
}
