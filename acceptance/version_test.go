package acceptance

import (
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
)

func TestVersion(t *testing.T) {
	httpBaseURL, _ := TestServer(layouts.Handler)
	client := ClientV2{BaseURL: httpBaseURL}

	got := client.GetVersion(t).Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
