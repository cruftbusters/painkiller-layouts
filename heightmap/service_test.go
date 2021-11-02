package heightmap

import (
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/testing"
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
		AssertMetadata(t, gotFirst, wantFirst)

		gotSecond := service.post(Metadata{Size: "second size"})
		wantSecond := Metadata{Id: "second", Size: "second size"}
		AssertMetadata(t, gotSecond, wantSecond)

		AssertMetadata(t, *service.get(gotFirst.Id), gotFirst)
		AssertMetadata(t, *service.get(gotSecond.Id), gotSecond)
	})

	t.Run("delete heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}
		service.post(Metadata{Size: ""})
		service.Delete("the id")
		got := service.get("the id")
		if got != nil {
			t.Fatalf("got %v want nil", got)
		}
	})
}

func assertGetIsNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.get("")
	if metadata != nil {
		t.Fatal("got metadata but want nil")
	}
}
