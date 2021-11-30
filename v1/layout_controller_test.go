package v1

import (
	"errors"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
	"github.com/stretchr/testify/mock"
)

func TestLayoutController(t *testing.T) {
	mockLayoutService := new(MockLayoutService)
	controller := LayoutController{
		mockLayoutService,
	}

	httpBaseURL, _ := TestController(controller)
	client := ClientV2{BaseURL: httpBaseURL}

	t.Run("get missing", func(t *testing.T) {
		id := "deadbeef"
		mockLayoutService.On("Get", id).Return(Layout{}, ErrLayoutNotFound)
		client.GetLayoutExpectNotFound(t, id)
	})

	t.Run("patch missing", func(t *testing.T) {
		id := "william"
		mockLayoutService.On("Patch", id, mock.Anything).Return(Layout{}, ErrLayoutNotFound)
		client.PatchLayoutExpectNotFound(t, id)
	})

	t.Run("create", func(t *testing.T) {
		up, down := Layout{Id: "up"}, Layout{Id: "down"}
		mockLayoutService.On("Create", up).Return(down).Once()

		got := client.CreateLayout(t, up)
		AssertLayout(t, got, down)
	})

	t.Run("get", func(t *testing.T) {
		id, down := "up", Layout{Id: "down"}
		mockLayoutService.On("Get", id).Return(down, nil)
		got := client.GetLayout(t, id)
		AssertLayout(t, got, down)
	})

	t.Run("get all", func(t *testing.T) {
		down := []Layout{{Id: "beefdead"}}
		mockLayoutService.On("GetAll").Return(down)
		got := client.GetLayouts(t)
		AssertLayouts(t, got, down)
	})

	t.Run("patch", func(t *testing.T) {
		id, up, down := "rafael", Layout{HeightmapURL: "coming through"}, Layout{Id: "rafael", HeightmapURL: "coming through for real"}
		mockLayoutService.On("Patch", id, up).Return(down, nil)

		got := client.PatchLayout(t, id, up)
		AssertLayout(t, got, down)
	})

	t.Run("delete has error", func(t *testing.T) {
		id := "some id"
		mockLayoutService.On("Delete", id).Return(errors.New("uh oh"))
		client.DeleteLayoutExpectInternalServerError(t, id)
	})

	t.Run("delete", func(t *testing.T) {
		id := "another id"
		mockLayoutService.On("Delete", id).Return(nil)
		client.DeleteLayout(t, id)
	})
}
