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

	t.Run("create and get heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}

		got := service.Post(Metadata{})
		want := Metadata{Id: "the id"}
		AssertMetadata(t, got, want)

		AssertMetadata(t, *service.Get(got.Id), got)

		service.Delete("the id")
	})

	t.Run("create and get all heightmaps", func(t *testing.T) {
		stubUuidService.idQueue = []string{"first", "second"}

		service.Post(Metadata{})
		service.Post(Metadata{})

		got := service.GetAll()
		want := []Metadata{
			Metadata{Id: "first"},
			Metadata{Id: "second"},
		}
		AssertAllMetadataUnordered(t, got, want)
	})

	t.Run("patch url onto metadata", func(t *testing.T) {
		id, size, url := "the id", "old size", "new image url"
		stubUuidService.idQueue = []string{id}
		service.Post(Metadata{Size: size})

		got := service.Patch(id, Metadata{ImageURL: url})
		want := Metadata{Id: id, Size: size, ImageURL: url}
		AssertMetadata(t, got, want)

		got = *service.Get(id)
		AssertMetadata(t, got, want)
	})

	t.Run("delete heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}
		service.Post(Metadata{Size: ""})
		service.Delete("the id")
		got := service.Get("the id")
		if got != nil {
			t.Fatalf("got %v want nil", got)
		}
	})
}

func assertGetIsNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.Get("")
	if metadata != nil {
		t.Fatal("got metadata but want nil")
	}
}
