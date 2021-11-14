package acceptance

import (
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/layouts"
	. "github.com/cruftbusters/painkiller-gallery/testing"
)

func TestVersion(t *testing.T) {
	listener, baseURL := RandomPortListener()
	go func() { http.Serve(listener, layouts.Handler(baseURL)) }()

	client := NewClientV2(t, baseURL)

	got := client.GetVersion().Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
