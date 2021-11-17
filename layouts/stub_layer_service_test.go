package layouts

import (
	"bytes"
	"testing"
)

type StubLayerService struct {
	t                        testing.TB
	whenPutCalledWithId      string
	whenPutCalledWithName    string
	whenPutCalledWithLayer   []byte
	putWillReturn            error
	whenGetCalledWithId      string
	whenGetCalledWithName    string
	getWillReturnLayer       []byte
	getWillReturnContentType string
	getWillReturnError       error
	whenDeleteCalledWithId   string
	whenDeleteCalledWithName string
	deleteWillReturn         error
}

func (stub *StubLayerService) Put(gotId, gotName string, gotLayer []byte) error {
	stub.t.Helper()
	wantId := stub.whenPutCalledWithId
	if gotId != wantId {
		stub.t.Errorf("got %s want %s", gotId, wantId)
	}
	wantName := stub.whenPutCalledWithName
	if gotName != wantName {
		stub.t.Errorf("got %s want %s", gotName, wantName)
	}
	wantLayer := stub.whenPutCalledWithLayer
	if !bytes.Equal(gotLayer, wantLayer) {
		stub.t.Errorf("got %v want %v", gotLayer, wantLayer)
	}
	return stub.putWillReturn
}

func (stub *StubLayerService) Get(gotId, gotName string) ([]byte, string, error) {
	stub.t.Helper()
	wantId := stub.whenGetCalledWithId
	if gotId != wantId {
		stub.t.Errorf("got %s want %s", gotId, wantId)
	}
	wantName := stub.whenGetCalledWithName
	if gotName != wantName {
		stub.t.Errorf("got %s want %s", gotName, wantName)
	}
	return stub.getWillReturnLayer, stub.getWillReturnContentType, stub.getWillReturnError
}

func (stub *StubLayerService) Delete(gotId, gotName string) error {
	stub.t.Helper()
	wantId := stub.whenDeleteCalledWithId
	if gotId != wantId {
		stub.t.Errorf("got %s want %s", gotId, wantId)
	}
	wantName := stub.whenDeleteCalledWithName
	if gotName != wantName {
		stub.t.Errorf("got %s want %s", gotName, wantName)
	}
	return stub.deleteWillReturn
}
