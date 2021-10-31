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
		assertPost(t, service, Metadata{Id: "deadbeef"})
		assertGet(t, service, Metadata{Id: "deadbeef"})
	})
}

func assertGetIsNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.get()
	if metadata != nil {
		t.Fatal("got metadata but want nil")
	}
}

func assertGet(t testing.TB, service Service, want Metadata) {
	t.Helper()
	got := *service.get()
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func assertPost(t testing.TB, service Service, want Metadata) {
	t.Helper()
	got := service.post()
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
