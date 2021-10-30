package main

import (
	"net/http"
	"testing"
)

func TestHeightmaps(t *testing.T) {
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	t.Run("get missing heightmap", func(t *testing.T) {
		response, err := http.Get("http://localhost:8080/v1/heightmaps/deadbeef")
		if err != nil {
			t.Fatal("got error wanted no error", err)
		}

		assertStatusCode(t, response, 404)
	})
}

func assertStatusCode(t testing.TB, response *http.Response, statusCode int) {
	t.Helper()
	if response.StatusCode != statusCode {
		t.Fatalf("got status code %d want %d", response.StatusCode, statusCode)
	}
}
