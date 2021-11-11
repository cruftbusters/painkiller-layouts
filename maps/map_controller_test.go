package maps

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
	"github.com/julienschmidt/httprouter"
)

func TestMapController(t *testing.T) {
	listener, port := RandomPortListener()
	client := NewClientV2(t, fmt.Sprintf("http://localhost:%d", port))

	stubService := &StubService{t: t}
	stubHeightmapService := &StubHeightmapService{t: t}
	controller := MapController{
		stubService,
		stubHeightmapService,
	}
	router := httprouter.New()
	controller.AddRoutes(router)

	go func() {
		http.Serve(listener, router)
	}()

	t.Run("get missing map", func(t *testing.T) {
		stubService.whenGetCalledWith = "deadbeef"
		stubService.getWillReturnError = MapNotFoundError

		client.GetExpectNotFound("deadbeef")
	})

	t.Run("create map", func(t *testing.T) {
		up, down := Metadata{Id: "up"}, Metadata{Id: "down"}
		stubService.whenPostCalledWith = up
		stubService.postWillReturn = down

		got := client.Create(up)
		AssertMetadata(t, got, down)
	})

	t.Run("get map", func(t *testing.T) {
		stubService.whenGetCalledWith = "path-id"
		stubService.getWillReturnMetadata = Metadata{Id: "beefdead"}
		stubService.getWillReturnError = nil

		got := client.Get("path-id")
		want := Metadata{Id: "beefdead"}
		AssertMetadata(t, got, want)
	})

	t.Run("get all maps", func(t *testing.T) {
		stubService.getAllWillReturn = []Metadata{{Id: "beefdead"}}

		got := client.GetAll()
		want := []Metadata{{Id: "beefdead"}}
		AssertAllMetadata(t, got, want)
	})

	t.Run("patch missing map", func(t *testing.T) {
		id := "william"
		stubService.whenPatchCalledWithId = id
		stubService.whenPatchCalledWithMetadata = Metadata{}
		stubService.patchWillReturnError = MapNotFoundError

		client.PatchExpectNotFound(id)
	})

	t.Run("patch map by id", func(t *testing.T) {
		id, up, down := "rafael", Metadata{ImageURL: "coming through"}, Metadata{Id: "rafael", ImageURL: "coming through for real"}
		stubService.whenPatchCalledWithId = id
		stubService.whenPatchCalledWithMetadata = up
		stubService.patchWillReturnMetadata = down
		stubService.patchWillReturnError = nil

		got := client.Patch(id, up)
		want := down
		AssertMetadata(t, got, want)
	})

	t.Run("delete map has error", func(t *testing.T) {
		id, want := "some id", errors.New("uh oh")
		stubService.whenDeleteCalledWith = id
		stubService.deleteWillRaise = want

		client.DeleteExpectInternalServerError(id)
	})

	t.Run("delete map", func(t *testing.T) {
		id := "some id"
		stubService.whenDeleteCalledWith = id
		stubService.deleteWillRaise = nil

		client.Delete(id)
	})

	t.Run("put heightmap on missing map is not found", func(t *testing.T) {
		id := "there is no creativity"
		stubHeightmapService.whenPutCalledWithId = id
		stubHeightmapService.putWillReturn = MapNotFoundError

		client.PutHeightmapExpectNotFound(id)
	})

	t.Run("get heightmap on missing map is not found", func(t *testing.T) {
		id := "walrus"
		stubHeightmapService.whenGetCalledWith = id
		stubHeightmapService.getWillReturnError = MapNotFoundError

		client.GetHeightmapExpectNotFound(id)
	})

	t.Run("get heightmap is not found", func(t *testing.T) {
		id := "serendipity"
		stubHeightmapService.whenGetCalledWith = id
		stubHeightmapService.getWillReturnError = HeightmapNotFoundError

		client.GetHeightmapExpectNotFound(id)
	})

	t.Run("put heightmap", func(t *testing.T) {
		id, up := "john denver", []byte("was a bear")
		stubHeightmapService.whenPutCalledWithId = id
		stubHeightmapService.whenPutCalledWithHeightmap = up
		stubHeightmapService.putWillReturn = nil

		client.PutHeightmap(id, bytes.NewBuffer(up))
	})

	t.Run("get heightmap", func(t *testing.T) {
		id, heightmap, contentType := "inwards", []byte("buncha bytes"), "image/png"
		stubHeightmapService.whenGetCalledWith = id
		stubHeightmapService.getWillReturnHeightmap = heightmap
		stubHeightmapService.getWillReturnContentType = contentType
		stubHeightmapService.getWillReturnError = nil

		gotReadCloser, gotContentType := client.GetHeightmap(id)
		got, err := io.ReadAll(gotReadCloser)
		AssertNoError(t, err)
		want := heightmap
		if bytes.Compare(got, want) != 0 {
			t.Errorf("got %v want %v", got, want)
		}
		if gotContentType != contentType {
			t.Errorf("got %s want %s", gotContentType, contentType)
		}
	})
}
