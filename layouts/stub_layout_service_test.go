package layouts

import (
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubLayoutService struct {
	t                           testing.TB
	whenGetCalledWith           string
	getWillReturnMetadata       Metadata
	getWillReturnError          error
	whenGetAllCalledWith        bool
	getAllWillReturn            []Metadata
	whenPostCalledWith          Metadata
	postWillReturn              Metadata
	whenPatchCalledWithId       string
	whenPatchCalledWithMetadata Metadata
	patchWillReturnMetadata     Metadata
	patchWillReturnError        error
	whenDeleteCalledWith        string
	deleteWillReturn            error
}

func (stub *StubLayoutService) Create(got Metadata) Metadata {
	stub.t.Helper()
	want := stub.whenPostCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.postWillReturn
}

func (stub *StubLayoutService) Get(got string) (Metadata, error) {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.getWillReturnMetadata, stub.getWillReturnError
}

func (stub *StubLayoutService) GetAll(got bool) []Metadata {
	stub.t.Helper()
	want := stub.whenGetAllCalledWith
	if got != want {
		stub.t.Fatalf("got %t want %t", got, want)
	}
	return stub.getAllWillReturn
}

func (stub *StubLayoutService) Patch(gotId string, gotMetadata Metadata) (Metadata, error) {
	stub.t.Helper()
	wantId := stub.whenPatchCalledWithId
	wantMetadata := stub.whenPatchCalledWithMetadata
	if gotId != wantId {
		stub.t.Fatalf("got %s want %s", gotId, wantId)
	}
	if gotMetadata != wantMetadata {
		stub.t.Fatalf("got %#v want %#v", gotMetadata, wantMetadata)
	}
	return stub.patchWillReturnMetadata, stub.patchWillReturnError
}

func (stub *StubLayoutService) Delete(got string) error {
	stub.t.Helper()
	if got != stub.whenDeleteCalledWith {
		stub.t.Errorf("got id %s want %s", got, stub.whenDeleteCalledWith)
	}
	return stub.deleteWillReturn
}