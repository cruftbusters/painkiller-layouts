package maps

import "testing"

type StubHeightmapService struct {
	t                          testing.TB
	whenPutCalledWithId        string
	whenPutCalledWithHeightmap string
	putWillReturn              error
	whenGetCalledWith          string
	getWillReturnHeightmap     string
	getWillReturnError         error
}

func (stub *StubHeightmapService) Put(gotId, gotHeightmap string) error {
	stub.t.Helper()
	wantId := stub.whenPutCalledWithId
	if gotId != wantId {
		stub.t.Errorf("got %s want %s", gotId, wantId)
	}
	wantHeightmap := stub.whenPutCalledWithHeightmap
	if gotHeightmap != wantHeightmap {
		stub.t.Errorf("got %s want %s", gotHeightmap, wantHeightmap)
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
