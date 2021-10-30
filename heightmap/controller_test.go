package heightmap

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController(t *testing.T) {
	t.Run("get missing heightmap", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/heightmaps/deadbeef", nil)
		response := httptest.NewRecorder()

		HeightmapController(response, request)

		assertStatusCode(t, response, 404)
	})

	t.Run("create heightmap", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/v1/heightmaps", nil)
		response := httptest.NewRecorder()

		HeightmapController(response, request)

		assertStatusCode(t, response, 201)
	})
}

func assertStatusCode(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Fatalf("got status code %d want %d", response.Code, want)
	}
}
