package acceptance

import (
	"testing"

	"github.com/cruftbusters/painkiller-gallery/layouts"
	. "github.com/cruftbusters/painkiller-gallery/testing"
)

func TestVersion(t *testing.T) {
	client, _ := NewTestClient(t, layouts.Handler)

	got := client.GetVersion().Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
