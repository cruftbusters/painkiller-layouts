package v1

import (
	"errors"
	"reflect"
	"testing"

	"github.com/cruftbusters/painkiller-layouts/types"
)

func TestLayoutAwaitingLayerWire(t *testing.T) {
	layoutService := new(MockLayoutService)
	awaitingHeightmap := new(MockAwaitingLayerService)
	service := NewLayoutAwaitingLayerWire(
		layoutService,
		awaitingHeightmap,
	)

	t.Run("proxy create", func(t *testing.T) {
		up, down := types.Layout{Id: "hey up"}, types.Layout{Id: "hey down"}
		layoutService.On("Create", up).Return(down)
		got := service.Create(up)
		if got != down {
			t.Errorf("got %+v want %+v", got, down)
		}
	})

	t.Run("proxy delete", func(t *testing.T) {
		id, err := "delete this", errors.New("problem deleting")
		layoutService.On("Delete", id).Return(err)
		got := service.Delete(id)
		if got != err {
			t.Errorf("got %s want %s", got, err)
		}
	})

	t.Run("proxy get", func(t *testing.T) {
		id, layout, err := "get this", types.Layout{Id: "ok got it"}, errors.New("problem getting")
		layoutService.On("Get", id).Return(layout, err)
		got, gotErr := service.Get(id)
		if got != layout {
			t.Errorf("got %+v want %+v", got, layout)
		}
		if gotErr != err {
			t.Errorf("got %s want %s", gotErr, err)
		}
	})

	t.Run("proxy get all", func(t *testing.T) {
		layouts := []types.Layout{{Id: "get all these"}}
		layoutService.On("GetAll").Return(layouts)
		got := service.GetAll()
		if !reflect.DeepEqual(got, layouts) {
			t.Errorf("got %+v want %+v", got, layouts)
		}
	})

	t.Run("proxy get all with heightmap without hillshade", func(t *testing.T) {
		layouts := []types.Layout{{Id: "get all these with heightmap without hillshade"}}
		layoutService.On("GetAllWithHeightmapWithoutHillshade").Return(layouts)
		got := service.GetAllWithHeightmapWithoutHillshade()
		if !reflect.DeepEqual(got, layouts) {
			t.Errorf("got %+v want %+v", got, layouts)
		}
	})

	t.Run("proxy get all with no heightmap", func(t *testing.T) {
		layouts := []types.Layout{{Id: "get all these with no heightmap"}}
		layoutService.On("GetAllWithNoHeightmap").Return(layouts)
		got := service.GetAllWithNoHeightmap()
		if !reflect.DeepEqual(got, layouts) {
			t.Errorf("got %+v want %+v", got, layouts)
		}
	})

	t.Run("proxy patch", func(t *testing.T) {
		id, up, down, err := "patch me", types.Layout{Id: "here comes the patch"}, types.Layout{Id: "result"}, errors.New("patch broke")
		layoutService.On("Patch", id, up).Return(down, err)
		got, gotErr := service.Patch(id, up)
		if got != down {
			t.Errorf("got %+v want %+v", got, down)
		}
		if gotErr != err {
			t.Errorf("got %s want %s", gotErr, err)
		}
	})
}
