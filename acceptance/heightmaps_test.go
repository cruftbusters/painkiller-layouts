package acceptance

import (
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/heightmap"
)

func TestHeightmaps(t *testing.T) {
	go func() {
		http.ListenAndServe(":8080", heightmap.Handler())
	}()

	t.Run("get missing heightmap", func(t *testing.T) {
		response, err := http.Get("http://localhost:8080/v1/heightmaps/deadbeef")
		assertNoError(t, err)

		assertStatusCode(t, response, 404)
	})

	t.Run("create new heightmap", func(t *testing.T) {
		postResponse, err := http.Post("http://localhost:8080/v1/heightmaps", "", nil)
		assertNoError(t, err)

		assertStatusCode(t, postResponse, 201)

		response, err := http.Get("http://localhost:8080/v1/heightmaps/deadbeef")
		assertNoError(t, err)

		assertStatusCode(t, response, 200)
	})
}

func assertNoError(t testing.TB, err error) {
	if err != nil {
		t.Fatal("got error wanted no error", err)
	}
}

func assertStatusCode(t testing.TB, response *http.Response, statusCode int) {
	t.Helper()
	if response.StatusCode != statusCode {
		t.Fatalf("got status code %d want %d", response.StatusCode, statusCode)
	}
}
