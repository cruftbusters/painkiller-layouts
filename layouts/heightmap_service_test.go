package layouts

import (
	"reflect"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func TestHeightmapService(t *testing.T) {
	stubLayoutService := &StubLayoutService{t: t}
	heightmapService := NewHeightmapService(
		"http://baseURL",
		stubLayoutService,
	)

	t.Run("put when map not found", func(t *testing.T) {
		id, err := "not found", ErrLayoutNotFound
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = err

		got := heightmapService.Put(id, nil)
		AssertError(t, got, err)
	})

	t.Run("get when map not found", func(t *testing.T) {
		id, err := "wimbly wombly", ErrLayoutNotFound
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = err

		_, _, got := heightmapService.Get(id)
		AssertError(t, got, err)
	})

	t.Run("get when heightmap not found", func(t *testing.T) {
		id := "weeknights"
		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		_, _, got := heightmapService.Get(id)
		AssertError(t, got, ErrHeightmapNotFound)
	})

	t.Run("put and get", func(t *testing.T) {
		id, heightmap, contentType := "bhan mi", []byte("vegan impossible burger"), "image/jpeg"
		heightmapURL := "http://baseURL/v1/maps/bhan mi/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		err := heightmapService.Put(id, heightmap)
		AssertNoError(t, err)

		got, gotContentType, err := heightmapService.Get(id)
		AssertNoError(t, err)
		want := heightmap
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
		if gotContentType != contentType {
			t.Errorf("got %s want %s", gotContentType, contentType)
		}
	})

	t.Run("put heightmap updates heightmap URL", func(t *testing.T) {
		id, heightmapURL := "itchy", "http://baseURL/v1/maps/itchy/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = nil

		heightmapService.Put(id, nil)
	})

	t.Run("put heightmap has error upon updating heightmap URL", func(t *testing.T) {
		id, heightmapURL := "stitchy", "http://baseURL/v1/maps/stitchy/heightmap.jpg"

		stubLayoutService.whenGetCalledWith = id
		stubLayoutService.getWillReturnError = nil

		stubLayoutService.whenPatchCalledWithId = id
		stubLayoutService.whenPatchCalledWithLayout = Layout{HeightmapURL: heightmapURL}
		stubLayoutService.patchWillReturnLayout = Layout{}
		stubLayoutService.patchWillReturnError = ErrLayoutNotFound

		got := heightmapService.Put(id, nil)
		AssertError(t, got, ErrLayoutNotFound)
	})
}
