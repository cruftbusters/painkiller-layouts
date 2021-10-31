package heightmap

import (
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubUUIDService struct {
	id string
}

func (service StubUUIDService) NewUUID() string {
	return service.id
}

func TestService(t *testing.T) {
	uuidService := StubUUIDService{"beefdead"}
	service := &DefaultService{uuidService: uuidService}
	t.Run("get when missing", func(t *testing.T) {
		assertGetIsNil(t, service)
	})

	t.Run("get after post", func(t *testing.T) {
		up, expectedDown := Metadata{
			Id:   "ignore",
			Size: "dont ignore",
		}, Metadata{
			Id:   "beefdead",
			Size: "dont ignore",
		}

		down := service.post(up)
		assertMetadata(t, down, expectedDown)
		assertMetadata(t, *service.get(""), expectedDown)
	})
}

func assertGetIsNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.get("")
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
