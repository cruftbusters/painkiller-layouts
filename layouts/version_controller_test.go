package layouts

import (
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
	"github.com/julienschmidt/httprouter"
)

func TestVersionController(t *testing.T) {
	controller := VersionController{}
	client, _ := NewClientV2(t, func(string) *httprouter.Router {
		router := httprouter.New()
		controller.AddRoutes(router)
		return router
	})

	got := client.GetVersion().Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
