package layouts

import (
	"database/sql"
	"reflect"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func TestLayerService(t *testing.T) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	Migrate(db)
	stubLayoutService := &StubLayoutService{t: t}
	layerService := NewLayerService(
		"http://baseURL",
		db,
		stubLayoutService,
	)

	t.Run("put when layout not found", func(t *testing.T) {
		id, err := "not found", ErrLayoutNotFound
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = err

		got := layerService.Put(id, "heightmap.jpg", nil)
		AssertError(t, got, err)
	})

	t.Run("get when layout not found", func(t *testing.T) {
		id, err := "wimbly wombly", ErrLayoutNotFound
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = err

		_, _, got := layerService.Get(id, "heightmap.jpg")
		AssertError(t, got, err)
	})

	t.Run("get when layer not found", func(t *testing.T) {
		id := "weeknights"
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		_, _, got := layerService.Get(id, "heightmap.jpg")
		AssertError(t, got, ErrLayerNotFound)
	})

	t.Run("put and get layers", func(t *testing.T) {
		id := "bhan mi"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		steps := []struct {
			name, contentType string
			layer             []byte
			patch             Layout
		}{
			{
				name:        "heightmap.jpg",
				contentType: "image/jpeg",
				layer:       []byte("vegan impossible burger"),
				patch: Layout{
					HeightmapURL: "http://baseURL/v1/layouts/bhan mi/heightmap.jpg",
				},
			},
			{
				name:        "hillshade.jpg",
				contentType: "image/jpeg",
				layer:       []byte("mega wompus"),
				patch: Layout{
					HeightmapURL: "http://baseURL/v1/layouts/bhan mi/hillshade.jpg",
				},
			},
		}

		for _, step := range steps {
			stubLayoutService.whenPatchCalledWithId = id
			stubLayoutService.whenPatchCalledWithLayout = step.patch
			stubLayoutService.patchWillReturnLayout = Layout{}
			stubLayoutService.patchWillReturnError = nil

			err := layerService.Put(id, step.name, step.layer)
			AssertNoError(t, err)
		}

		for _, step := range steps {
			gotLayer, gotContentType, err := layerService.Get(id, step.name)
			AssertNoError(t, err)
			wantLayer := step.layer
			if !reflect.DeepEqual(gotLayer, wantLayer) {
				t.Errorf("got %v want %v", gotLayer, wantLayer)
			}
			if gotContentType != step.contentType {
				t.Errorf("got %s want %s", gotContentType, step.contentType)
			}
		}
	})

	t.Run("put heightmap updates heightmap URL", func(t *testing.T) {
		id, name, heightmapURL := "itchy", "heightmap.jpg", "http://baseURL/v1/layouts/itchy/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		layerService.Put(id, name, nil)
	})

	t.Run("put heightmap has error upon updating heightmap URL", func(t *testing.T) {
		id, name, heightmapURL := "stitchy", "heightmap.jpg", "http://baseURL/v1/layouts/stitchy/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = ErrLayoutNotFound

		got := layerService.Put(id, name, nil)
		AssertError(t, got, ErrLayoutNotFound)
	})

	t.Run("update heightmap", func(t *testing.T) {
		id, name, heightmapURL := "why am i specifying this again", "heightmap.jpg", "http://baseURL/v1/layouts/why am i specifying this again/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		layerService.Put(id, name, []byte("deja vu"))
		layerService.Put(id, name, []byte("deja vu"))
	})
}
