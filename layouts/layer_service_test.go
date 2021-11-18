package layouts

import (
	"bytes"
	"database/sql"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
	"github.com/stretchr/testify/mock"
)

func TestLayerService(t *testing.T) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	Migrate(db)
	mockLayoutService := new(MockLayoutService)
	layerService := NewLayerService(
		"http://baseURL",
		db,
		mockLayoutService,
	)

	t.Run("put when layout not found", func(t *testing.T) {
		id, err := "not found", ErrLayoutNotFound
		mockLayoutService.On("Get", id).Return(Layout{}, err)
		got := layerService.Put(id, "heightmap.jpg", nil)
		AssertError(t, got, err)
	})

	t.Run("get when layer not found", func(t *testing.T) {
		_, _, got := layerService.Get("weeknights", "heightmap.jpg")
		AssertError(t, got, ErrLayerNotFound)
	})

	t.Run("put get delete", func(t *testing.T) {
		id := "bhan mi"
		mockLayoutService.On("Get", id).Return(Layout{}, nil)

		steps := []struct {
			name, contentType string
			layer             []byte
		}{
			{
				name:        "heightmap.jpg",
				contentType: "image/jpeg",
				layer:       []byte("vegan impossible burger"),
			},
			{
				name:        "hillshade.jpg",
				contentType: "image/jpeg",
				layer:       []byte("mega wompus"),
			},
		}

		for _, step := range steps {
			mockLayoutService.On("Patch", id, mock.Anything).Return(Layout{}, nil)
			err := layerService.Put(id, step.name, step.layer)
			AssertNoError(t, err)
		}

		for _, step := range steps {
			gotLayer, gotContentType, err := layerService.Get(id, step.name)
			AssertNoError(t, err)
			wantLayer := step.layer
			if !bytes.Equal(gotLayer, wantLayer) {
				t.Errorf("got %v want %v", gotLayer, wantLayer)
			}
			if gotContentType != step.contentType {
				t.Errorf("got %s want %s", gotContentType, step.contentType)
			}
		}

		t.Run("put again", func(t *testing.T) {
			for _, step := range steps {
				mockLayoutService.On("Patch", id, mock.Anything).Return(Layout{}, nil)
				err := layerService.Put(id, step.name, step.layer)
				AssertNoError(t, err)
			}
		})

		t.Run("delete", func(t *testing.T) {
			for _, step := range steps {
				layerService.Delete(id, step.name)
				_, _, err := layerService.Get(id, "heightmap.jpg")
				AssertError(t, err, ErrLayerNotFound)
			}
		})
	})

	t.Run("put layer updates corresponding layout URL", func(t *testing.T) {
		id := "not so unique"
		mockLayoutService.On("Get", id).Return(Layout{}, nil)

		for _, scenario := range []struct {
			name  string
			patch Layout
		}{
			{
				name:  "heightmap.jpg",
				patch: Layout{HeightmapURL: "http://baseURL/v1/layouts/not so unique/heightmap.jpg"},
			},
			{
				name:  "hillshade.jpg",
				patch: Layout{HillshadeURL: "http://baseURL/v1/layouts/not so unique/hillshade.jpg"},
			},
		} {
			mockLayoutService.On("Patch", id, scenario.patch).Return(Layout{}, nil).Once()
			err := layerService.Put(id, scenario.name, nil)
			AssertNoError(t, err)

			t.Run("patch layout has error", func(t *testing.T) {
				mockLayoutService.On("Patch", id, scenario.patch).Return(Layout{}, ErrLayoutNotFound)
				err := layerService.Put(id, scenario.name, nil)
				AssertError(t, err, ErrLayoutNotFound)
			})
		}
	})

	t.Run("layers are present after deleting layout", func(t *testing.T) {
		id := "wimbly wombly"

		mockLayoutService.On("Get", id).Return(Layout{}, nil)
		mockLayoutService.On("Patch", id, mock.Anything).Return(Layout{}, nil)

		layerService.Put(id, "heightmap.jpg", nil)

		mockLayoutService.On("Get", id).Return(ErrLayerNotFound)

		_, _, err := layerService.Get(id, "heightmap.jpg")
		AssertNoError(t, err)
	})
}
