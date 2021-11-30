package acceptance

import (
	"fmt"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
)

func TestLayout(t *testing.T) {
	httpBaseURL, wsBaseURL := TestServer(v1.Handler)
	client := ClientV2{BaseURL: httpBaseURL}

	close, err := DrainAwaitingLayers(wsBaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	t.Run("get missing layout", func(t *testing.T) {
		client.GetLayoutExpectNotFound(t, "deadbeef")
	})

	t.Run("patch missing layout", func(t *testing.T) {
		client.PatchLayoutExpectNotFound(t, "garbotron")
	})

	t.Run("create and get layout", func(t *testing.T) {
		got := client.CreateLayout(t, Layout{})
		defer client.DeleteLayout(t, got.Id)
		AssertLayout(t, got, Layout{Id: got.Id})
		AssertLayout(t, client.GetLayout(t, got.Id), got)
	})

	t.Run("delete layout", func(t *testing.T) {
		layout := client.CreateLayout(t, Layout{})
		client.DeleteLayout(t, layout.Id)
		client.GetLayoutExpectNotFound(t, layout.Id)
	})

	t.Run("get all layouts", func(t *testing.T) {
		first := client.CreateLayout(t, Layout{})
		defer client.DeleteLayout(t, first.Id)
		second := client.CreateLayout(t, Layout{
			Scale:             9.87,
			Size:              Size{Width: 1, Height: 2},
			Bounds:            Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
			HeightmapURL:      "heightmap url",
			HiResHeightmapURL: "hi res heightmap url",
			HillshadeURL:      "hillshade url",
			HiResHillshadeURL: "hi res hillshade url",
		})
		defer client.DeleteLayout(t, second.Id)

		AssertLayoutsUnordered(t,
			client.GetLayouts(t),
			[]Layout{first, second},
		)
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
			layout := client.CreateLayout(t,
				Layout{
					Scale:             0.75,
					Size:              Size{Width: 1, Height: 2},
					Bounds:            Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
					HeightmapURL:      "old heightmap url",
					HiResHeightmapURL: "old hi res heightmap url",
					HillshadeURL:      "old hillshade url",
					HiResHillshadeURL: "old hi res hillshade url",
				},
			)
			defer client.DeleteLayout(t, layout.Id)

			got := client.PatchLayout(t, layout.Id, scenario.patch)
			scenario.want(&layout)
			AssertLayout(t, got, layout)

			got = client.GetLayout(t, layout.Id)
			AssertLayout(t, got, layout)
		})
	}
}
