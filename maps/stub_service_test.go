package maps

import (
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubService struct {
	t                           testing.TB
	whenGetCalledWith           string
	getWillReturnMetadata       Metadata
	getWillReturnError          error
	getAllWillReturn            []Metadata
	whenPostCalledWith          Metadata
	postWillReturn              Metadata
	whenPatchCalledWithId       string
	whenPatchCalledWithMetadata Metadata
	patchWillReturn             Metadata
	whenDeleteCalledWith        string
	deleteWillRaise             error
}

func (stub *StubService) Post(got Metadata) Metadata {
	stub.t.Helper()
	want := stub.whenPostCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.postWillReturn
}

func (stub *StubService) Get(got string) (Metadata, error) {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.getWillReturnMetadata, stub.getWillReturnError
}

func (stub *StubService) GetAll() []Metadata {
	return stub.getAllWillReturn
}

func (stub *StubService) Patch(gotId string, gotMetadata Metadata) Metadata {
	stub.t.Helper()
	wantId := stub.whenPatchCalledWithId
	wantMetadata := stub.whenPatchCalledWithMetadata
	if gotId != wantId {
		stub.t.Fatalf("got %s want %s", gotId, wantId)
	}
	if gotMetadata != wantMetadata {
		stub.t.Fatalf("got %#v want %#v", gotMetadata, wantMetadata)
	}
	return stub.patchWillReturn
}

func (stub *StubService) Delete(got string) error {
	stub.t.Helper()
	if got != stub.whenDeleteCalledWith {
		stub.t.Errorf("got id %s want %s", got, stub.whenDeleteCalledWith)
	}
	return stub.deleteWillRaise
}
