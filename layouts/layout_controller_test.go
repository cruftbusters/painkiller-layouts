package layouts

import (
	"errors"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
)

func TestLayoutController(t *testing.T) {
	stubLayoutService := &StubLayoutService{t: t}
	controller := LayoutController{
		stubLayoutService,
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

		client.DeleteLayoutExpectInternalServerError(id)
	})

	t.Run("delete map", func(t *testing.T) {
		id := "some id"
		stubLayoutService.whenDeleteCalledWith = id
		stubLayoutService.deleteWillReturn = nil

		client.DeleteLayout(id)
	})
}
