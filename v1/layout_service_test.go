package v1

import (
	"database/sql"
	"fmt"
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

	t.Run("get missing layout", func(t *testing.T) {
		_, got := service.Get("")
		AssertError(t, got, ErrLayoutNotFound)
	})

	t.Run("patch missing layout", func(t *testing.T) {
		_, err := service.Patch("pragmatism", Layout{})
		AssertError(t, err, ErrLayoutNotFound)
	})

	t.Run("create get delete", func(t *testing.T) {
		id := "windows update"
		stubUuidService.idQueue = []string{id}

		got := service.Create(Layout{})
		AssertLayout(t, got, Layout{Id: id})

		layout, err := service.Get(got.Id)
		AssertNoError(t, err)
		AssertLayout(t, layout, got)

		service.Delete(id)

		_, err = service.Get(id)
		AssertError(t, err, ErrLayoutNotFound)
	})

	t.Run("get all", func(t *testing.T) {
		stubUuidService.idQueue = []string{"first", "second"}

		first := service.Create(Layout{})
		defer service.Delete("first")
		second := service.Create(Layout{
			Scale:             9.45,
			Size:              Size{Width: 1, Height: 2},
			Bounds:            Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
			HeightmapURL:      "heightmap url",
			HiResHeightmapURL: "hi res heightmap url",
			HillshadeURL:      "hillshade url",
			HiResHillshadeURL: "hi res hillshade url",
		})
		defer service.Delete("second")

		got := service.GetAll()
		AssertLayoutsUnordered(t, got, []Layout{first, second})
	})

	for _, scenario := range []struct {
		patch Layout
		want  func(*Layout)
	}{
		{
			patch: Layout{HeightmapURL: "new heightmap url"},
			want:  func(initial *Layout) { initial.HeightmapURL = "new heightmap url" },
		},
		{
			patch: Layout{HiResHeightmapURL: "new hi res heightmap url"},
			want:  func(initial *Layout) { initial.HiResHeightmapURL = "new hi res heightmap url" },
		},
		{
			patch: Layout{HillshadeURL: "new hillshade url"},
			want:  func(initial *Layout) { initial.HillshadeURL = "new hillshade url" },
		},
		{
			patch: Layout{HiResHillshadeURL: "new hi res hillshade url"},
			want:  func(initial *Layout) { initial.HiResHillshadeURL = "new hi res hillshade url" },
		},
		{
			patch: Layout{Scale: 12.34},
			want:  func(initial *Layout) { initial.Scale = 12.34 },
		},
	} {
		t.Run(fmt.Sprintf("patch layout with %+v", scenario.patch), func(t *testing.T) {
			id := "the id"
			stubUuidService.idQueue = []string{id}
			layout := service.Create(
				Layout{
					Scale:             0.12345,
					Size:              Size{Width: 1, Height: 2},
					Bounds:            Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
					HeightmapURL:      "old heightmap url",
					HiResHeightmapURL: "old hi res heightmap url",
					HillshadeURL:      "old hillshade url",
					HiResHillshadeURL: "old hi res hillshade url",
				},
			)
			defer service.Delete(id)

			got, err := service.Patch(id, scenario.patch)
			AssertNoError(t, err)
			scenario.want(&layout)
			AssertLayout(t, got, layout)

			got, err = service.Get(id)
			AssertNoError(t, err)
			AssertLayout(t, got, layout)
		})
	}
}
