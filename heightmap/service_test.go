package heightmap

import (
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubUUIDService struct {
	idQueue []string
}

func (service *StubUUIDService) NewUUID() string {
	nextId := service.idQueue[0]
	service.idQueue = service.idQueue[1:]
	return nextId
}

func TestService(t *testing.T) {
	stubUuidService := &StubUUIDService{}
	service := NewService(stubUuidService)
	t.Run("get when missing", func(t *testing.T) {
		assertGetIsNil(t, service)
	})

	t.Run("create and get two heightmaps", func(t *testing.T) {
		stubUuidService.idQueue = []string{"first", "second"}

		gotFirst := service.post(Metadata{Size: "first size"})
		wantFirst := Metadata{Id: "first", Size: "first size"}
		assertMetadata(t, gotFirst, wantFirst)

		gotSecond := service.post(Metadata{Size: "second size"})
		wantSecond := Metadata{Id: "second", Size: "second size"}
		assertMetadata(t, gotSecond, wantSecond)

		assertMetadata(t, *service.get(gotFirst.Id), gotFirst)
		assertMetadata(t, *service.get(gotSecond.Id), gotSecond)
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
