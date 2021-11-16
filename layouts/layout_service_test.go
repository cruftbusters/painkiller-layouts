package layouts

import (
	"database/sql"
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
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	Migrate(db)
	stubUuidService := &StubUUIDService{}
	service := NewLayoutService(db, stubUuidService)
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

		withoutHeightmap := service.Create(Layout{
			Size:   Size{Width: 1, Height: 2},
			Bounds: Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
		})
		defer func() { service.Delete("first") }()
		withHeightmap := service.Create(Layout{HeightmapURL: heightmapURL})
		defer func() { service.Delete("second") }()

		got := service.GetAll(false)
		AssertLayoutsUnordered(t, got, []Layout{withoutHeightmap, withHeightmap})

		got = service.GetAll(true)
		AssertLayoutsUnordered(t, got, []Layout{withoutHeightmap})
	})

	t.Run("patch missing map", func(t *testing.T) {
		_, err := service.Patch("pragmatism", Layout{})
		AssertError(t, err, ErrLayoutNotFound)
	})

	t.Run("patch url onto layout", func(t *testing.T) {
		id := "the id"
		stubUuidService.idQueue = []string{id}
		service.Create(
			Layout{
				Size:         Size{Width: 1, Height: 2},
				Bounds:       Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
				HeightmapURL: "old heightmap URL",
			},
		)
		defer func() { service.Delete(id) }()

		got, err := service.Patch(id, Layout{HeightmapURL: "new heightmap URL"})
		AssertNoError(t, err)
		AssertLayout(t, got, Layout{
			Id:           id,
			Size:         Size{Width: 1, Height: 2},
			Bounds:       Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
			HeightmapURL: "new heightmap URL",
		})

		got, err = service.Get(id)
		AssertNoError(t, err)
		AssertLayout(t, got, Layout{
			Id:           id,
			Size:         Size{Width: 1, Height: 2},
			Bounds:       Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
			HeightmapURL: "new heightmap URL",
		})
	})

	t.Run("delete heightmap", func(t *testing.T) {
		stubUuidService.idQueue = []string{"the id"}
		service.Create(Layout{})
		service.Delete("the id")
		_, got := service.Get("the id")
		AssertError(t, got, ErrLayoutNotFound)
	})
}
