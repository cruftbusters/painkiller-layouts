package maps

import "testing"

type StubHeightmapService struct {
	t                 testing.TB
	whenPutCalledWith string
	putWillReturn     error
	whenGetCalledWith string
	getWillReturn     error
}

func (stub *StubHeightmapService) put(got string) error {
	stub.t.Helper()
	want := stub.whenPutCalledWith
	if got != want {
		stub.t.Errorf("got %s want %s", got, want)
	}
	return stub.putWillReturn
}

func (stub *StubHeightmapService) Get(got string) error {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Errorf("got %s want %s", got, want)
	}
	return stub.getWillReturn
}
