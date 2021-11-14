package acceptance

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/maps"
	. "github.com/cruftbusters/painkiller-gallery/testing"
)

func TestVersion(t *testing.T) {
	listener, port := RandomPortListener()
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	go func() {
		http.Serve(listener, maps.Handler(baseURL))
	}()

	client := NewClientV2(t, baseURL)

	got := client.GetVersion().Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}