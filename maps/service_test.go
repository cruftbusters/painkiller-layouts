package maps

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
		_, got := service.Get("")
		want := MapNotFoundError
		if got != want {
			t.Fatalf("got %v want %v", got, want)
		}
	})

	t.Run("create and get heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}

		got := service.Post(Metadata{})
		want := Metadata{Id: "the id"}
		AssertMetadata(t, got, want)

		metadata, err := service.Get(got.Id)
		AssertNoError(t, err)
		AssertMetadata(t, metadata, got)

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
		id, size, url := "the id", Size{1, 2}, "new image url"
		stubUuidService.idQueue = []string{id}
		service.Post(Metadata{Size: size})

		got := service.Patch(id, Metadata{ImageURL: url})
		want := Metadata{Id: id, Size: size, ImageURL: url}
		AssertMetadata(t, got, want)

		got, err := service.Get(id)
		AssertNoError(t, err)
		AssertMetadata(t, got, want)
	})

	t.Run("delete heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}
		service.Post(Metadata{})
		service.Delete("the id")
		_, got := service.Get("the id")
		want := MapNotFoundError
		if got != want {
			t.Fatalf("got %v want %v", got, want)
		}
	})
}
