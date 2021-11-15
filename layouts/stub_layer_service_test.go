package layouts

import (
	"bytes"
	"testing"
)

type StubLayerService struct {
	t                        testing.TB
	whenPutCalledWithId      string
	whenPutCalledWithLayer   []byte
	putWillReturn            error
	whenGetCalledWithId      string
	getWillReturnLayer       []byte
	getWillReturnContentType string
	getWillReturnError       error
}

func (stub *StubLayerService) Put(gotId string, gotLayer []byte) error {
	stub.t.Helper()
	wantId := stub.whenPutCalledWithId
	if gotId != wantId {
		stub.t.Errorf("got %s want %s", gotId, wantId)
	}
	wantLayer := stub.whenPutCalledWithLayer
	if !bytes.Equal(gotLayer, wantLayer) {
		stub.t.Errorf("got %v want %v", gotLayer, wantLayer)
	}
	return stub.putWillReturn
}

func (stub *StubLayerService) Get(got string) ([]byte, string, error) {
	stub.t.Helper()
	want := stub.whenGetCalledWithId
	if got != want {
		stub.t.Errorf("got %s want %s", got, want)
	}
	return stub.getWillReturnLayer, stub.getWillReturnContentType, stub.getWillReturnError
}
