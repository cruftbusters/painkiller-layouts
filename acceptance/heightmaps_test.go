package acceptance

import (
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/heightmap"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestHeightmaps(t *testing.T) {
	go func() {
		http.ListenAndServe(":8080", heightmap.Handler())
	}()

	client := Client{t: t, BaseUrl: "http://localhost:8080"}

	t.Run("get missing heightmap", func(t *testing.T) {
		client.GetMetadataExpectNotFound()
	})

	t.Run("create new heightmap", func(t *testing.T) {
		metadata := client.Create(Metadata{
			Size: "large",
		})

		assertMetadata(t,
			metadata,
			Metadata{
				Id:   metadata.Id,
				Size: "large",
			})

		assertMetadata(t, client.GetMetadata(), metadata)
	})
}

func assertMetadata(t testing.TB, got Metadata, want Metadata) {
	t.Helper()
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
