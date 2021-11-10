package maps

import "testing"

type StubHeightmapService struct {
	t                 testing.TB
	whenPutCalledWith string
	putWillReturn     error
}

func (stub *StubHeightmapService) put(got string) error {
	stub.t.Helper()
	want := stub.whenPutCalledWith
	if got != want {
		stub.t.Errorf("got %s want %s", got, want)
	}
	return stub.putWillReturn
}
