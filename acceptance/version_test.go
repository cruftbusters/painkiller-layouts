package acceptance

import (
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
)

func TestVersion(t *testing.T) {
	client, _ := NewTestClient(layouts.Handler)

	got := client.GetVersion(t).Version
	want := "1"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
