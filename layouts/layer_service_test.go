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
					HillshadeURL: "http://baseURL/v1/layouts/bhan mi/hillshade.jpg",
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

	t.Run("update layer", func(t *testing.T) {
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

	t.Run("put layer updates corresponding layout URL", func(t *testing.T) {
		for _, scenario := range []struct {
			id, name string
			patch    Layout
		}{
			{
				id:    "not so unique",
				name:  "heightmap.jpg",
				patch: Layout{HeightmapURL: "http://baseURL/v1/layouts/not so unique/heightmap.jpg"},
			},
			{
				id:    "not so distinct",
				name:  "hillshade.jpg",
				patch: Layout{HillshadeURL: "http://baseURL/v1/layouts/not so distinct/hillshade.jpg"},
			},
		} {
			t.Run("put heightmap updates heightmap URL", func(t *testing.T) {
				stubLayoutService.whenGetCalledWith = scenario.id
				stubLayoutService.getWillReturnError = nil

				stubLayoutService.whenPatchCalledWithId = scenario.id
				stubLayoutService.whenPatchCalledWithLayout = scenario.patch
				stubLayoutService.patchWillReturnLayout = Layout{}
				stubLayoutService.patchWillReturnError = nil

				layerService.Put(scenario.id, scenario.name, nil)
			})

			t.Run("put heightmap has error upon updating heightmap URL", func(t *testing.T) {
				stubLayoutService.whenGetCalledWith = scenario.id
				stubLayoutService.getWillReturnError = nil

				stubLayoutService.whenPatchCalledWithId = scenario.id
				stubLayoutService.whenPatchCalledWithLayout = scenario.patch
				stubLayoutService.patchWillReturnLayout = Layout{}
				stubLayoutService.patchWillReturnError = ErrLayoutNotFound

				got := layerService.Put(scenario.id, scenario.name, nil)
				AssertError(t, got, ErrLayoutNotFound)
			})
		}
	})
}
