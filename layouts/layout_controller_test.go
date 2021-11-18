package layouts

import (
	"errors"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
)

func TestLayoutController(t *testing.T) {
	mockLayoutService := new(MockLayoutService)
	controller := LayoutController{
		mockLayoutService,
	}
	client, _ := NewTestClient(t, func(string, string) *httprouter.Router {
		router := httprouter.New()
		controller.AddRoutes(router)
		return router
	})

	t.Run("get missing", func(t *testing.T) {
		id := "deadbeef"
		mockLayoutService.On("Get", id).Return(Layout{}, ErrLayoutNotFound)
		client.GetLayoutExpectNotFound(id)
	})

	t.Run("patch missing", func(t *testing.T) {
		id := "william"
		mockLayoutService.On("Patch", id, mock.Anything).Return(Layout{}, ErrLayoutNotFound)
		client.PatchLayoutExpectNotFound(id)
	})

	t.Run("create", func(t *testing.T) {
		up, down := Layout{Id: "up"}, Layout{Id: "down"}
		mockLayoutService.On("Create", up).Return(down)
		got := client.CreateLayout(up)
		AssertLayout(t, got, down)
	})

	t.Run("get", func(t *testing.T) {
		id, down := "up", Layout{Id: "down"}
		mockLayoutService.On("Get", id).Return(down, nil)
		got := client.GetLayout(id)
		AssertLayout(t, got, down)
	})

	t.Run("get all", func(t *testing.T) {
		down := []Layout{{Id: "beefdead"}}
		mockLayoutService.On("GetAll").Return(down)
		got := client.GetLayouts()
		AssertLayouts(t, got, down)
	})

	t.Run("get all with no heightmap", func(t *testing.T) {
		down := []Layout{{Id: "look ma no heightmap"}}
		mockLayoutService.On("GetAllWithNoHeightmap").Return(down)
		got := client.GetLayoutsWithoutHeightmap()
		AssertLayouts(t, got, down)
	})

	t.Run("get all with no hillshade", func(t *testing.T) {
		down := []Layout{{Id: "look ma no hillshade"}}
		mockLayoutService.On("GetAllWithHeightmapWithoutHillshade").Return(down)
		got := client.GetLayoutsWithHeightmapWithoutHillshade()
		AssertLayouts(t, got, down)
	})

	t.Run("patch by id", func(t *testing.T) {
		id, up, down := "rafael", Layout{HeightmapURL: "coming through"}, Layout{Id: "rafael", HeightmapURL: "coming through for real"}
		mockLayoutService.On("Patch", id, up).Return(down, nil)
		got := client.PatchLayout(id, up)
		AssertLayout(t, got, down)
	})

	t.Run("delete has error", func(t *testing.T) {
		id := "some id"
		mockLayoutService.On("Delete", id).Return(errors.New("uh oh"))
		client.DeleteLayoutExpectInternalServerError(id)
	})

	t.Run("delete", func(t *testing.T) {
		id := "another id"
		mockLayoutService.On("Delete", id).Return(nil)
		client.DeleteLayout(id)
	})
}
