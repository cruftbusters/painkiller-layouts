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
		AssertError(t, got, ErrMapNotFound)
	})

	t.Run("create and get heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}

		got := service.Create(Metadata{})
		want := Metadata{Id: "the id"}
		AssertMetadata(t, got, want)

		metadata, err := service.Get(got.Id)
		AssertNoError(t, err)
		AssertMetadata(t, metadata, got)

		service.Delete("the id")
	})

	t.Run("create and get all heightmaps", func(t *testing.T) {
		stubUuidService.idQueue = []string{"first", "second"}
		heightmapURL := "better not drop me"

		service.Create(Metadata{})
		service.Create(Metadata{HeightmapURL: heightmapURL})

		got := service.GetAll(false)
		want := []Metadata{{Id: "first"}, {Id: "second", HeightmapURL: "better not drop me"}}
		AssertAllMetadataUnordered(t, got, want)

		service.Delete("first")
		service.Delete("second")
	})

	t.Run("patch missing map", func(t *testing.T) {
		_, err := service.Patch("pragmatism", Metadata{})
		AssertError(t, err, ErrMapNotFound)
	})

	t.Run("patch url onto metadata", func(t *testing.T) {
		id, size, heightmapURL := "the id", Size{Width: 1, Height: 2}, "new heightmap url"
		stubUuidService.idQueue = []string{id}
		service.Create(Metadata{Size: size})

		got, err := service.Patch(id, Metadata{HeightmapURL: heightmapURL})
		AssertNoError(t, err)
		want := Metadata{Id: id, Size: size, HeightmapURL: heightmapURL}
		AssertMetadata(t, got, want)

		got, err = service.Get(id)
		AssertNoError(t, err)
		AssertMetadata(t, got, want)

		service.Delete(id)
	})

	t.Run("filter for maps with no heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"first", "second"}
		withoutHeightmap := service.Create(Metadata{})
		withHeightmap := service.Create(Metadata{HeightmapURL: "heightmap url"})
		AssertAllMetadata(t,
			service.GetAll(true),
			[]Metadata{withoutHeightmap},
		)

		service.Delete(withoutHeightmap.Id)
		service.Delete(withHeightmap.Id)
	})

	t.Run("delete heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}
		service.Create(Metadata{})
		service.Delete("the id")
		_, got := service.Get("the id")
		AssertError(t, got, ErrMapNotFound)
	})
}
