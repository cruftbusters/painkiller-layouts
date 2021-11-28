package v1

import (
	"errors"
	"reflect"
	"testing"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/stretchr/testify/mock"
)

func TestLayoutAwaitingLayerWire(t *testing.T) {
	layoutService := new(MockLayoutService)
	awaitingHeightmap := new(MockAwaitingLayerService)
	awaitingHillshade := new(MockAwaitingLayerService)
	service := NewLayoutAwaitingLayerWire(
		layoutService,
		awaitingHeightmap,
		awaitingHillshade,
	)

	t.Run("proxy create and notify awaiting heightmap", func(t *testing.T) {
		up, down := types.Layout{Id: "hey up"}, types.Layout{Id: "hey down"}
		awaitingHeightmap.On("Enqueue", mock.Anything).Return(nil).Once()
		layoutService.On("Create", up).Return(down).Once()
		got := service.Create(up)
		if got != down {
			t.Errorf("got %+v want %+v", got, down)
		}
		awaitingHeightmap.AssertCalled(t, "Enqueue", down)
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
		id, up, down, err := "patch me", types.Layout{Id: "up up and away"}, types.Layout{Id: "down down and away"}, errors.New("patch broke")
		layoutService.On("Patch", id, up).Return(down, err).Once()
		got, gotErr := service.Patch(id, up)
		if got != down {
			t.Errorf("got %+v want %+v", got, down)
		}
		if gotErr != err {
			t.Errorf("got %s want %s", gotErr, err)
		}

		t.Run("notify awaiting hillshades when heightmap URL is populated", func(t *testing.T) {
			up := types.Layout{HeightmapURL: "not blank"}
			layoutService.On("Patch", id, up).Return(down, err).Once()
			awaitingHillshade.On("Enqueue", down).Return(nil).Once()
			service.Patch(id, up)
			awaitingHillshade.AssertCalled(t, "Enqueue", down)
		})
	})
}