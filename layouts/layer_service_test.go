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

	t.Run("put and get", func(t *testing.T) {
		id := "bhan mi"
		name, layer, contentType := "heightmap.jpg", []byte("vegan impossible burger"), "image/jpeg"
		heightmapURL := "http://baseURL/v1/layouts/bhan mi/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		err := layerService.Put(id, name, layer)
		AssertNoError(t, err)

		got, gotContentType, err := layerService.Get(id, name)
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
