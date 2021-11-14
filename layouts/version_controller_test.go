package layouts

import (
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
	"github.com/julienschmidt/httprouter"
)

func TestVersionController(t *testing.T) {
	listener, baseURL := RandomPortListener()
	client := NewClientV2(t, baseURL)

	controller := VersionController{}
	router := httprouter.New()
	controller.AddRoutes(router)

	go func() { http.Serve(listener, router) }()

	got := client.GetVersion().Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
