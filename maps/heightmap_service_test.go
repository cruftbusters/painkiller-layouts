package maps

import "testing"

func TestHeightmapService(t *testing.T) {
	stubMapService := &StubService{t: t}
	heightmapService := DefaultHeightmapService{
		stubMapService,
	}

	t.Run("put returns error when map not found", func(t *testing.T) {
		id, err := "not found", MapNotFoundError
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = err

		got := heightmapService.put(id)
		want := err
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("put", func(t *testing.T) {
		id := "found"
		stubMapService.whenGetCalledWith = id
		stubMapService.getWillReturnError = nil

		got := heightmapService.put(id)
		if got != nil {
			t.Errorf("got %v want nil", got)
		}
	})
}
