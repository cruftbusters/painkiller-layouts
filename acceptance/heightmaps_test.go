package acceptance

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/heightmap"
	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestHeightmaps(t *testing.T) {
	listener, port := RandomPortListener()
	go func() {
		http.Serve(listener, heightmap.Handler())
	}()

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	client := NewClient(t, baseURL)

	t.Run("get missing heightmap", func(t *testing.T) {
		client.GetExpectNotFound()
	})

	t.Run("create and get two heightmaps", func(t *testing.T) {
		gotFirst := client.Create(Metadata{Size: "first"})
		wantFirst := Metadata{Id: gotFirst.Id, Size: "first"}
		AssertMetadata(t, gotFirst, wantFirst)

		gotSecond := client.Create(Metadata{Size: "second"})
		wantSecond := Metadata{Id: gotSecond.Id, Size: "second"}
		AssertMetadata(t, gotSecond, wantSecond)

		AssertMetadata(t, client.Get(gotFirst.Id), gotFirst)
		AssertMetadata(t, client.Get(gotSecond.Id), gotSecond)
	})
}
