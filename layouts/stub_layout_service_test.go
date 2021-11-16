package layouts

import (
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/types"
)

type StubLayoutService struct {
	t                                             testing.TB
	whenGetCalledWith                             string
	getWillReturnLayout                           Layout
	getWillReturnError                            error
	getAllWillReturn                              []Layout
	getAllWithNoHeightmapWillReturn               []Layout
	getAllWithHeightmapWithoutHillshadeWillReturn []Layout
	whenPostCalledWith                            Layout
	postWillReturn                                Layout
	whenPatchCalledWithId                         string
	whenPatchCalledWithLayout                     Layout
	patchWillReturnLayout                         Layout
	patchWillReturnError                          error
	whenDeleteCalledWith                          string
	deleteWillReturn                              error
}

func (stub *StubLayoutService) Create(got Layout) Layout {
	stub.t.Helper()
	want := stub.whenPostCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.postWillReturn
}

func (stub *StubLayoutService) Get(got string) (Layout, error) {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.getWillReturnLayout, stub.getWillReturnError
}

func (stub *StubLayoutService) GetAll() []Layout {
	return stub.getAllWillReturn
}

func (stub *StubLayoutService) GetAllWithNoHeightmap() []Layout {
	return stub.getAllWithNoHeightmapWillReturn
}

func (stub *StubLayoutService) GetAllWithHeightmapWithoutHillshade() []Layout {
	return stub.getAllWithHeightmapWithoutHillshadeWillReturn
}

func (stub *StubLayoutService) Patch(gotId string, gotLayout Layout) (Layout, error) {
	stub.t.Helper()
	wantId := stub.whenPatchCalledWithId
	wantLayout := stub.whenPatchCalledWithLayout
	if gotId != wantId {
		stub.t.Fatalf("got %s want %s", gotId, wantId)
	}
	if gotLayout != wantLayout {
		stub.t.Fatalf("got %#v want %#v", gotLayout, wantLayout)
	}
	return stub.patchWillReturnLayout, stub.patchWillReturnError
}

func (stub *StubLayoutService) Delete(got string) error {
	stub.t.Helper()
	if got != stub.whenDeleteCalledWith {
		stub.t.Errorf("got id %s want %s", got, stub.whenDeleteCalledWith)
	}
	return stub.deleteWillReturn
}
