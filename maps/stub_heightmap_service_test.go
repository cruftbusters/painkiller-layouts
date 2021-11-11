package maps

import "testing"

type StubHeightmapService struct {
	t                      testing.TB
	whenPutCalledWith      string
	putWillReturn          error
	whenGetCalledWith      string
	getWillReturnHeightmap string
	getWillReturnError     error
}

func (stub *StubHeightmapService) Put(got string) error {
	stub.t.Helper()
	want := stub.whenPutCalledWith
	if got != want {
		stub.t.Errorf("got %s want %s", got, want)
	}
	return stub.putWillReturn
}

func (stub *StubHeightmapService) Get(got string) (string, error) {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Errorf("got %s want %s", got, want)
	}
	return stub.getWillReturnHeightmap, stub.getWillReturnError
}
