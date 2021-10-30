package main

import (
	"net/http"
	"testing"
)

func TestHeightmaps(t *testing.T) {
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	t.Run("connect", func(t *testing.T) {
		_, err := http.Get("http://localhost:8080/v1/heightmaps/deadbeef")
		if err != nil {
			t.Fatal("got error wanted no error", err)
		}
	})
}
