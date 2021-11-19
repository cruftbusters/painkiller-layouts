package layouts

import (
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
)

func TestVersionController(t *testing.T) {
	controller := VersionController{}

	httpBaseURL, _ := TestController(controller)
	client := ClientV2{BaseURL: httpBaseURL}

	got := client.GetVersion(t).Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
