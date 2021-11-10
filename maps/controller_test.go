package maps

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
	"github.com/julienschmidt/httprouter"
)

func TestController(t *testing.T) {
	listener, port := RandomPortListener()
	client := NewClient(t, fmt.Sprintf("http://localhost:%d", port))

	stubService := &StubService{t: t}
	controller := &Controller{stubService}
	router := httprouter.New()
	controller.AddRoutes(router)

	go func() {
		http.Serve(listener, router)
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

	t.Run("get all heightmaps", func(t *testing.T) {
		stubService.getAllWillReturn = []Metadata{{Id: "beefdead"}}

		got := client.GetAll()
		want := []Metadata{{Id: "beefdead"}}
		AssertAllMetadata(t, got, want)
	})

	t.Run("patch heightmap by id", func(t *testing.T) {
		id, up, down := "rafael", Metadata{ImageURL: "coming through"}, Metadata{Id: "rafael", ImageURL: "coming through for real"}
		stubService.whenPatchCalledWithId = id
		stubService.whenPatchCalledWithMetadata = up
		stubService.patchWillReturn = down

		got := client.Patch(id, up)
		want := down
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
