package heightmap

import "testing"

func TestUUIDService(t *testing.T) {
	got := DefaultUUIDService{}.NewUUID()
	want := "deadbeef"

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
