package layouts

import (
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
)

type StubUUIDService struct {
	idQueue []string
}

func (service *StubUUIDService) NewUUID() string {
	nextId := service.idQueue[0]
	service.idQueue = service.idQueue[1:]
	return nextId
}

func TestLayoutService(t *testing.T) {
	stubUuidService := &StubUUIDService{}
	service := NewLayoutService(stubUuidService)
	t.Run("get when missing", func(t *testing.T) {
		_, got := service.Get("")
		AssertError(t, got, ErrLayoutNotFound)
	})

	t.Run("create and get heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}

		got := service.Create(Layout{})
		defer func() { service.Delete("the id") }()
		want := Layout{Id: "the id"}
		AssertLayout(t, got, want)

		layout, err := service.Get(got.Id)
		AssertNoError(t, err)
		AssertLayout(t, layout, got)
	})

	t.Run("create and get all heightmaps", func(t *testing.T) {
		stubUuidService.idQueue = []string{"first", "second"}
		heightmapURL := "better not drop me"

		service.Create(Layout{})
		defer func() { service.Delete("first") }()
		service.Create(Layout{HeightmapURL: heightmapURL})
		defer func() { service.Delete("second") }()

		got := service.GetAll(false)
		want := []Layout{{Id: "first"}, {Id: "second", HeightmapURL: "better not drop me"}}
		AssertLayoutsUnordered(t, got, want)
	})

	t.Run("patch missing map", func(t *testing.T) {
		_, err := service.Patch("pragmatism", Layout{})
		AssertError(t, err, ErrLayoutNotFound)
	})

	t.Run("patch url onto layout", func(t *testing.T) {
		id, size, heightmapURL := "the id", Size{Width: 1, Height: 2}, "new heightmap url"
		stubUuidService.idQueue = []string{id}
		service.Create(Layout{Size: size})
		defer func() { service.Delete(id) }()

		got, err := service.Patch(id, Layout{HeightmapURL: heightmapURL})
		AssertNoError(t, err)
		want := Layout{Id: id, Size: size, HeightmapURL: heightmapURL}
		AssertLayout(t, got, want)

		got, err = service.Get(id)
		AssertNoError(t, err)
		AssertLayout(t, got, want)
	})

	t.Run("filter for maps with no heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"first", "second"}
		withoutHeightmap := service.Create(Layout{})
		defer func() { service.Delete(withoutHeightmap.Id) }()
		withHeightmap := service.Create(Layout{HeightmapURL: "heightmap url"})
		defer func() { service.Delete(withHeightmap.Id) }()
		AssertLayouts(t,
			service.GetAll(true),
			[]Layout{withoutHeightmap},
		)
	})

	t.Run("delete heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}
		service.Create(Layout{})
		service.Delete("the id")
		_, got := service.Get("the id")
		AssertError(t, got, ErrLayoutNotFound)
	})
}
