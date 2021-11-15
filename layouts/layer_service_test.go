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
	stubLayoutService := &StubLayoutService{t: t}
	layerService := NewLayerService(
		"http://baseURL",
		db,
		stubLayoutService,
	)

	t.Run("put when map not found", func(t *testing.T) {
		id, err := "not found", ErrLayoutNotFound
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = err

		got := layerService.Put(id, nil)
		AssertError(t, got, err)
	})

	t.Run("get when map not found", func(t *testing.T) {
		id, err := "wimbly wombly", ErrLayoutNotFound
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = err

		_, _, got := layerService.Get(id)
		AssertError(t, got, err)
	})

	t.Run("get when heightmap not found", func(t *testing.T) {
		id := "weeknights"
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		_, _, got := layerService.Get(id)
		AssertError(t, got, ErrLayerNotFound)
	})

	t.Run("put and get", func(t *testing.T) {
		id, layer, contentType := "bhan mi", []byte("vegan impossible burger"), "image/jpeg"
		heightmapURL := "http://baseURL/v1/layouts/bhan mi/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		err := layerService.Put(id, layer)
		AssertNoError(t, err)

		got, gotContentType, err := layerService.Get(id)
		AssertNoError(t, err)
		want := layer
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
		if gotContentType != contentType {
			t.Errorf("got %s want %s", gotContentType, contentType)
		}
	})

	t.Run("put heightmap updates heightmap URL", func(t *testing.T) {
		id, heightmapURL := "itchy", "http://baseURL/v1/layouts/itchy/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		layerService.Put(id, nil)
	})

	t.Run("put heightmap has error upon updating heightmap URL", func(t *testing.T) {
		id, heightmapURL := "stitchy", "http://baseURL/v1/layouts/stitchy/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = ErrLayoutNotFound

		got := layerService.Put(id, nil)
		AssertError(t, got, ErrLayoutNotFound)
	})

	t.Run("update heightmap", func(t *testing.T) {
		id, heightmapURL := "why am i specifying this again", "http://baseURL/v1/layouts/why am i specifying this again/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		layerService.Put(id, []byte("deja vu"))
		layerService.Put(id, []byte("deja vu"))
	})
}
