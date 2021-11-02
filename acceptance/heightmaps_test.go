package acceptance

import (
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/assertions"
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

	t.Run("create and get two heightmaps", func(t *testing.T) {
		gotFirst := client.Create(Metadata{Size: "first"})
		wantFirst := Metadata{Id: gotFirst.Id, Size: "first"}
		AssertMetadata(t, gotFirst, wantFirst)

		gotSecond := client.Create(Metadata{Size: "second"})
		wantSecond := Metadata{Id: gotSecond.Id, Size: "second"}
		AssertMetadata(t, gotSecond, wantSecond)

		AssertMetadata(t, client.GetMetadata(gotFirst.Id), gotFirst)
		AssertMetadata(t, client.GetMetadata(gotSecond.Id), gotSecond)
	})
}
