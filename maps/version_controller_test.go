package maps

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
	"github.com/julienschmidt/httprouter"
)

func TestVersionController(t *testing.T) {
	listener, port := RandomPortListener()
	baseURL := fmt.Sprintf("http://localhost:%d", port)

	controller := VersionController{}
	router := httprouter.New()
	controller.AddRoutes(router)

	go func() {
		http.Serve(listener, router)
	}()

	client := NewClientV2(t, baseURL)

	got := client.GetVersion().Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
