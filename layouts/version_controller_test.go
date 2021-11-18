package layouts

import (
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/julienschmidt/httprouter"
)

func TestVersionController(t *testing.T) {
	controller := VersionController{}
	client, _ := NewTestClient(func(string, string) *httprouter.Router {
		router := httprouter.New()
		controller.AddRoutes(router)
		return router
	})

	got := client.GetVersion(t).Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
