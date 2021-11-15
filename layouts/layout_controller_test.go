package layouts

import (
	"bytes"
	"errors"
	"io"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
)

func TestLayoutController(t *testing.T) {
	stubLayoutService := &StubLayoutService{t: t}
	stubHeightmapService := &StubHeightmapService{t: t}
	controller := LayoutController{
		stubLayoutService,
		stubHeightmapService,
	}
	client, _ := NewTestClient(t, func(string, string) *httprouter.Router {
		router := httprouter.New()
		controller.AddRoutes(router)
		return router
	})

	t.Run("get missing map", func(t *testing.T) {
		stubLayoutService.whenGetCalledWith = "deadbeef"
		stubLayoutService.getWillReturnError = ErrLayoutNotFound

		client.GetLayoutExpectNotFound("deadbeef")
	})

	t.Run("create map", func(t *testing.T) {
		up, down := Layout{Id: "up"}, Layout{Id: "down"}
		stubLayoutService.whenPostCalledWith = up
		stubLayoutService.postWillReturn = down

		got := client.CreateLayout(up)
		AssertLayout(t, got, down)
	})

	t.Run("get map", func(t *testing.T) {
		stubLayoutService.whenGetCalledWith = "path-id"
		stubLayoutService.getWillReturnLayout = Layout{Id: "beefdead"}
		stubLayoutService.getWillReturnError = nil

		got := client.GetLayout("path-id")
		want := Layout{Id: "beefdead"}
		AssertLayout(t, got, want)
	})

	t.Run("get all maps", func(t *testing.T) {
		stubLayoutService.whenGetAllCalledWith = false
		stubLayoutService.getAllWillReturn = []Layout{{Id: "beefdead"}}

		got := client.GetLayouts()
		want := []Layout{{Id: "beefdead"}}
		AssertLayouts(t, got, want)
	})

	t.Run("get all maps with heightmap URL filter", func(t *testing.T) {
		stubLayoutService.whenGetAllCalledWith = true
		stubLayoutService.getAllWillReturn = []Layout{{Id: "look ma no heightmap"}}

		got := client.GetLayoutsWithoutHeightmap()
		want := []Layout{{Id: "look ma no heightmap"}}
		AssertLayouts(t, got, want)
	})

	t.Run("patch missing map", func(t *testing.T) {
		id := "william"
		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{}
		stubLayoutService.patchWillReturnError = ErrLayoutNotFound

		client.PatchLayoutExpectNotFound(id)
	})

	t.Run("patch map by id", func(t *testing.T) {
		id, up, down := "rafael", Layout{HeightmapURL: "coming through"}, Layout{Id: "rafael", HeightmapURL: "coming through for real"}
		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = up
		stubLayoutService.patchWillReturnLayout = down
		stubLayoutService.patchWillReturnError = nil

		got := client.PatchLayout(id, up)
		want := down
		AssertLayout(t, got, want)
	})

	t.Run("delete map has error", func(t *testing.T) {
		id, want := "some id", errors.New("uh oh")
		stubLayoutService.whenDeleteCalledWith = id
		stubLayoutService.deleteWillReturn = want

		client.DeleteExpectInternalServerError(id)
	})

	t.Run("delete map", func(t *testing.T) {
		id := "some id"
		stubLayoutService.whenDeleteCalledWith = id
		stubLayoutService.deleteWillReturn = nil

		client.DeleteLayout(id)
	})

	t.Run("put heightmap on missing map is not found", func(t *testing.T) {
		id := "there is no creativity"
		stubHeightmapService.whenPutCalledWithId = id
		stubHeightmapService.putWillReturn = ErrLayoutNotFound

		client.PutLayoutHeightmapExpectNotFound(id)
	})

	t.Run("get heightmap on missing map is not found", func(t *testing.T) {
		id := "walrus"
		stubHeightmapService.whenGetCalledWith = id
		stubHeightmapService.getWillReturnError = ErrLayoutNotFound

		client.GetLayoutHeightmapExpectNotFound(id)
	})

	t.Run("get heightmap is not found", func(t *testing.T) {
		id := "serendipity"
		stubHeightmapService.whenGetCalledWith = id
		stubHeightmapService.getWillReturnError = ErrHeightmapNotFound

		client.GetLayoutHeightmapExpectNotFound(id)
	})

	t.Run("put heightmap", func(t *testing.T) {
		id, up := "john denver", []byte("was a bear")
		stubHeightmapService.whenPutCalledWithId = id
		stubHeightmapService.whenPutCalledWithHeightmap = up
		stubHeightmapService.putWillReturn = nil

		client.PutLayoutHeightmap(id, bytes.NewBuffer(up))
	})

	t.Run("get heightmap", func(t *testing.T) {
		id, heightmap, contentType := "inwards", []byte("buncha bytes"), "image/png"
		stubHeightmapService.whenGetCalledWith = id
		stubHeightmapService.getWillReturnHeightmap = heightmap
		stubHeightmapService.getWillReturnContentType = contentType
		stubHeightmapService.getWillReturnError = nil

		gotReadCloser, gotContentType := client.GetLayoutHeightmap(id)
		got, err := io.ReadAll(gotReadCloser)
		AssertNoError(t, err)
		want := heightmap
		if !bytes.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
		if gotContentType != contentType {
			t.Errorf("got %s want %s", gotContentType, contentType)
		}
	})
}
