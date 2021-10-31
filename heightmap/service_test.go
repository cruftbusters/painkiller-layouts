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
		up, expectedDown := Metadata{
			Id:   "ignore",
			Size: "dont ignore",
		}, Metadata{
			Id:   "deadbeef",
			Size: "dont ignore",
		}

		down := service.post(up)
		assertMetadata(t, down, expectedDown)
		assertMetadata(t, *service.get(), expectedDown)
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
