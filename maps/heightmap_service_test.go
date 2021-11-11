package maps

import . "github.com/cruftbusters/painkiller-gallery/testing"
import "testing"

func TestHeightmapService(t *testing.T) {
	stubMapService := &StubService{t: t}
	heightmapService := NewHeightmapService(
		stubMapService,
	)

	t.Run("put when map not found", func(t *testing.T) {
		id, err := "not found", MapNotFoundError
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = err

		got := heightmapService.Put(id)
		AssertError(t, got, err)
	})

	t.Run("get when map not found", func(t *testing.T) {
		id, err := "wimbly wombly", MapNotFoundError
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = err

		got := heightmapService.Get(id)
		AssertError(t, got, err)
	})

	t.Run("get when heightmap not found", func(t *testing.T) {
		id := "weeknights"
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		got := heightmapService.Get(id)
		AssertError(t, got, HeightmapNotFoundError)
	})

	t.Run("put and get", func(t *testing.T) {
		id := "bhan mi"
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		got := heightmapService.Put(id)
		AssertNoError(t, got)

		got = heightmapService.Get(id)
		AssertNoError(t, got)
	})
}
