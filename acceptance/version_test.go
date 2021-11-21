package acceptance

import (
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
)

func TestVersion(t *testing.T) {
	httpBaseURL, _ := TestServer(v1.Handler)
	client := ClientV2{BaseURL: httpBaseURL}

	got := client.GetVersion(t).Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
