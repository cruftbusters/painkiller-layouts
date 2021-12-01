package acceptance

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	. "github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
)

func TestLayers(t *testing.T) {
	httpBaseURL, _ := TestServer(v1.Handler)
	client := ClientV2{BaseURL: httpBaseURL}

	id := client.CreateLayout(t, Layout{}).Id
	defer client.DeleteLayout(t, id)

	t.Run("put layers on missing layout", func(t *testing.T) {
		client.PutLayerExpectNotFound(t, "deadbeef", "heightmap.jpg")
		client.PutLayerExpectNotFound(t, "deadbeef", "heightmap.tif")
		client.PutLayerExpectNotFound(t, "deadbeef", "hillshade.jpg")
		client.PutLayerExpectNotFound(t, "deadbeef", "hillshade.tif")
	})

	t.Run("get missing layers on layout", func(t *testing.T) {
		client.GetLayerExpectNotFound(t, id, "heightmap.jpg")
		client.GetLayerExpectNotFound(t, id, "heightmap.tif")
		client.GetLayerExpectNotFound(t, id, "hillshade.jpg")
		client.GetLayerExpectNotFound(t, id, "hillshade.tif")
	})

	for _, name := range []string{
		"heightmap.txt",
		"hillshade.png",
		"anything.jpg",
	} {
		t.Run(fmt.Sprintf("put '%s' bad request", name), func(t *testing.T) {
			client.PutLayerExpectBadRequest(t, id, name)
		})

		t.Run(fmt.Sprintf("get '%s' not found", name), func(t *testing.T) {
			client.GetLayerExpectNotFound(t, id, name)
		})
	}

	t.Run("put get delete", func(t *testing.T) {
		scenarios := []struct {
			name        string
			layer       []byte
			contentType string
		}{
			{
				name:        "heightmap.jpg",
				layer:       []byte("heightmap bytes"),
				contentType: "image/jpeg",
			},
			{
				name:        "heightmap.tif",
				layer:       []byte("hi res heightmap bytes"),
				contentType: "image/tiff",
			},
			{
				name:        "hillshade.jpg",
				layer:       []byte("hillshade bytes"),
				contentType: "image/jpeg",
			},
			{
				name:        "hillshade.tif",
				layer:       []byte("hi res hillshade bytes"),
				contentType: "image/tiff",
			},
		}

		for _, scenario := range scenarios {
			t.Run("put "+scenario.name, func(t *testing.T) {
				client.PutLayer(t, id, scenario.name, scenario.contentType, bytes.NewBuffer(scenario.layer))
			})
		}

		for _, scenario := range scenarios {
			t.Run("get "+scenario.name, func(t *testing.T) {
				gotReadCloser, gotContentType := client.GetLayer(t, id, scenario.name)
				got, err := io.ReadAll(gotReadCloser)
				AssertNoError(t, err)
				if !bytes.Equal(got, scenario.layer) {
					t.Errorf("got %v want %v", got, scenario.layer)
				}
				if gotContentType != scenario.contentType {
					t.Errorf("got %s want %s", gotContentType, scenario.contentType)
				}
			})
		}

		for _, scenario := range scenarios {
			t.Run("delete "+scenario.name, func(t *testing.T) {
				client.DeleteLayer(t, id, scenario.name)
				client.GetLayerExpectNotFound(t, id, scenario.name)
			})
		}
	})

	type URLSelector func(types.Layout) string
	for _, instance := range []struct {
		string
		URLSelector
	}{{"heightmap", func(layout types.Layout) string { return layout.HeightmapURL }}, {"hillshade", func(layout types.Layout) string { return layout.HillshadeURL }}} {
		t.Run(fmt.Sprintf("put /v1/layouts/:id/%s.jpg updates layout url", instance.string), func(t *testing.T) {
			want := "hello werald " + instance.string
			client.PutLayer(t, id, instance.string+".jpg", "image/jpeg", strings.NewReader(want))

			url := instance.URLSelector(client.GetLayout(t, id))
			response, err := http.Get(url)
			if err != nil {
				t.Fatal(err)
			} else if response.StatusCode != 200 {
				t.Fatalf("got status code %d want 200", response.StatusCode)
			}

			builder := new(strings.Builder)
			if _, err := io.Copy(builder, response.Body); err != nil {
				t.Fatal(err)
			}

			got := builder.String()
			if got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}

	t.Run("layers are present after deleting layout", func(t *testing.T) {
		id := client.CreateLayout(t, Layout{}).Id
		client.PutLayer(t, id, "heightmap.jpg", "", nil)
		client.DeleteLayout(t, id)
		client.GetLayer(t, id, "heightmap.jpg")
	})
}
