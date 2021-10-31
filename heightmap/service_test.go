package heightmap

import (
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestService(t *testing.T) {
	t.Run("get when missing", func(t *testing.T) {
		assertGetIsNil(t, &DefaultService{})
	})

	t.Run("get after post", func(t *testing.T) {
		service := &DefaultService{}
		assertMetadata(t, service.post(Metadata{}), Metadata{Id: "deadbeef"})
		assertMetadata(t, *service.get(), Metadata{Id: "deadbeef"})
	})
}

func assertGetIsNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.get()
	if metadata != nil {
		t.Fatal("got metadata but want nil")
	}
}

func assertMetadata(t testing.TB, got Metadata, want Metadata) {
	t.Helper()
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
