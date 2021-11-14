package acceptance

import (
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
)

func TestVersion(t *testing.T) {
	client, _ := NewTestClient(t, layouts.Handler)

	got := client.GetVersion().Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
