package v1

import "testing"

func TestUUIDService(t *testing.T) {
	got1 := DefaultUUIDService{}.NewUUID()
	got2 := DefaultUUIDService{}.NewUUID()

	if got1 == got2 {
		t.Errorf("want [%s, %s] to be distinct", got1, got2)
	}
}
