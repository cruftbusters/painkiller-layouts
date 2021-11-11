package maps

import (
	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
	"reflect"
	"testing"
)

func TestHeightmapService(t *testing.T) {
	stubMapService := &StubService{t: t}
	heightmapService := NewHeightmapService(
		"http://baseURL",
		stubMapService,
	)

	t.Run("put when map not found", func(t *testing.T) {
		id, err := "not found", MapNotFoundError
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = err

		got := heightmapService.Put(id, nil)
		AssertError(t, got, err)
	})

	t.Run("get when map not found", func(t *testing.T) {
		id, err := "wimbly wombly", MapNotFoundError
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = err

		_, _, got := heightmapService.Get(id)
		AssertError(t, got, err)
	})

	t.Run("get when heightmap not found", func(t *testing.T) {
		id := "weeknights"
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		_, _, got := heightmapService.Get(id)
		AssertError(t, got, HeightmapNotFoundError)
	})

	t.Run("put and get", func(t *testing.T) {
		id, heightmap, contentType := "bhan mi", []byte("vegan impossible burger"), "image/jpeg"

		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		stubMapService.whenPatchCalledWithId = id
		stubMapService.whenPatchCalledWithMetadata = Metadata{ImageURL: "http://baseURL/v1/maps/bhan mi/heightmap.jpg"}
		stubMapService.patchWillReturnMetadata = Metadata{}
		stubMapService.patchWillReturnError = nil

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

		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		stubMapService.whenPatchCalledWithId = id
		stubMapService.whenPatchCalledWithMetadata = Metadata{ImageURL: heightmapURL}
		stubMapService.patchWillReturnMetadata = Metadata{}
		stubMapService.patchWillReturnError = nil

		heightmapService.Put(id, nil)
	})

	t.Run("put heightmap has error upon updating heightmap URL", func(t *testing.T) {
		id, heightmapURL := "stitchy", "http://baseURL/v1/maps/stitchy/heightmap.jpg"

		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		stubMapService.whenPatchCalledWithId = id
		stubMapService.whenPatchCalledWithMetadata = Metadata{ImageURL: heightmapURL}
		stubMapService.patchWillReturnMetadata = Metadata{}
		stubMapService.patchWillReturnError = MapNotFoundError

		got := heightmapService.Put(id, nil)
		AssertError(t, got, MapNotFoundError)
	})
}
