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

		got := heightmapService.put(id)
		want := err
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("get when map not found", func(t *testing.T) {
		id, err := "wimbly wombly", MapNotFoundError
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = err

		got := heightmapService.Get(id)
		want := err
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("get when heightmap not found", func(t *testing.T) {
		id, err := "weeknights", HeightmapNotFoundError
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		got := heightmapService.Get(id)
		want := err
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("put and get", func(t *testing.T) {
		id := "bhan mi"
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		got := heightmapService.put(id)
		AssertNoError(t, got)

		got = heightmapService.Get(id)
		AssertNoError(t, got)
	})
}
