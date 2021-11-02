package heightmap

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubService struct {
	t                    testing.TB
	whenGetCalledWith    string
	getWillReturn        *Metadata
	whenPostCalledWith   Metadata
	postWillReturn       Metadata
	whenDeleteCalledWith string
	deleteWillRaise      error
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

func (stub *StubService) Delete(got string) error {
	stub.t.Helper()
	if got != stub.whenDeleteCalledWith {
		stub.t.Errorf("got id %s want %s", got, stub.whenDeleteCalledWith)
	}
	return stub.deleteWillRaise
}

func TestController(t *testing.T) {
	listener, port := RandomPortListener()
	client := NewClient(t, fmt.Sprintf("http://localhost:%d", port))

	stubService := &StubService{t: t}
	controller := Controller{
		stubService,
	}

	go func() {
		http.Serve(listener, controller.Router())
	}()

	t.Run("get missing heightmap", func(t *testing.T) {
		stubService.whenGetCalledWith = "deadbeef"
		stubService.getWillReturn = nil

		client.GetExpectNotFound("deadbeef")
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

	t.Run("delete heightmap has error", func(t *testing.T) {
		id, want := "some id", errors.New("uh oh")
		stubService.whenDeleteCalledWith = id
		stubService.deleteWillRaise = want

		client.DeleteExpectInternalServerError(id)
	})

	t.Run("delete heightmap", func(t *testing.T) {
		id := "some id"
		stubService.whenDeleteCalledWith = id
		stubService.deleteWillRaise = nil

		client.Delete(id)
	})
}
