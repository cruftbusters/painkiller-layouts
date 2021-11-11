package maps

import "bytes"
import "testing"

type StubHeightmapService struct {
	t                          testing.TB
	whenPutCalledWithId        string
	whenPutCalledWithHeightmap []byte
	putWillReturn              error
	whenGetCalledWith          string
	getWillReturnHeightmap     []byte
	getWillReturnContentType   string
	getWillReturnError         error
}

func (stub *StubHeightmapService) Put(gotId string, gotHeightmap []byte) error {
	stub.t.Helper()
	wantId := stub.whenPutCalledWithId
	if gotId != wantId {
		stub.t.Errorf("got %s want %s", gotId, wantId)
	}
	wantHeightmap := stub.whenPutCalledWithHeightmap
	if bytes.Compare(gotHeightmap, wantHeightmap) != 0 {
		stub.t.Errorf("got %v want %v", gotHeightmap, wantHeightmap)
	}
	return stub.putWillReturn
}

func (stub *StubHeightmapService) Get(got string) ([]byte, string, error) {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Errorf("got %s want %s", got, want)
	}
	return stub.getWillReturnHeightmap, stub.getWillReturnContentType, stub.getWillReturnError
}
