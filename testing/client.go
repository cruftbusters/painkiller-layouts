package testing

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/types"
)

type ClientV2 struct {
	BaseURL string
}

func (client ClientV2) GetVersion(t testing.TB) Version {
	t.Helper()
	response, err := http.Get(client.baseURLF("/version"))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)
	versionContainer, err := DecodeVersion(response)
	AssertNoError(t, err)
	return versionContainer
}

func (client ClientV2) GetLayout(t testing.TB, id string) Layout {
	t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s", id))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)

	return decode(t, response)
}

func (client ClientV2) GetLayouts(t testing.TB) []Layout {
	t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts"))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)

	return decodeLayouts(t, response)
}

func (client ClientV2) GetLayoutsWithoutHeightmap(t testing.TB) []Layout {
	t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts?excludeLayoutsWithHeightmap=true"))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)

	return decodeLayouts(t, response)
}

func (client ClientV2) GetLayoutsWithHeightmapWithoutHillshade(t testing.TB) []Layout {
	t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts?withHeightmapWithoutHillshade=true"))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)

	return decodeLayouts(t, response)
}

func (client ClientV2) GetLayoutExpectNotFound(t testing.TB, id string) {
	t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s", id))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 404)
}

func (client ClientV2) CreateLayout(t testing.TB, layout Layout) Layout {
	t.Helper()
	up := encode(t, layout)
	response, err := http.Post(client.baseURLF("/v1/layouts"), "", up)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 201)

	return decode(t, response)
}

func (client ClientV2) CreateLayoutExpect(layout Layout, statusCode int) (Layout, error) {
	up := &bytes.Buffer{}
	if err := json.NewEncoder(up).Encode(layout); err != nil {
		return layout, err
	}
	response, err := http.Post(client.baseURLF("/v1/layouts"), "", up)
	if err != nil {
		return layout, err
	} else if response.StatusCode != statusCode {
		return layout, fmt.Errorf("got status code %d want %d", response.StatusCode, statusCode)
	}
	if err := json.NewDecoder(response.Body).Decode(&layout); err != nil {
		return layout, err
	}
	return layout, nil
}

func (client ClientV2) CreateLayoutExpectInternalServerError(t testing.TB, layout Layout) {
	t.Helper()
	up := encode(t, layout)
	response, err := http.Post(client.baseURLF("/v1/layouts"), "", up)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 500)
}

func (client ClientV2) EnqueueLayoutAwaitingHeightmap(layout Layout) error {
	return client.EnqueueLayoutExpect("/v1/awaiting_heightmap", layout, 201)
}

func (client ClientV2) EnqueueLayoutAwaitingHeightmapExpectInternalServerError(layout Layout) error {
	channel := make(chan error)
	go func() { channel <- client.EnqueueLayoutExpect("/v1/awaiting_heightmap", layout, 500) }()
	select {
	case err := <-channel:
		return err
	case <-time.After(time.Second):
		return errors.New("timed out after one second")
	}
}

func (client ClientV2) EnqueueLayoutExpect(path string, layout Layout, statusCode int) error {
	up := &bytes.Buffer{}
	if err := json.NewEncoder(up).Encode(layout); err != nil {
		return err
	}
	response, err := http.Post(client.baseURLF(path), "", up)
	if err != nil {
		return err
	} else if response.StatusCode != statusCode {
		return fmt.Errorf("got status code %d want %d", response.StatusCode, statusCode)
	}
	return nil
}

func (client ClientV2) PatchLayoutExpectNotFound(t testing.TB, id string) {
	t.Helper()

	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	AssertNoError(t, err)

	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 404)
}

func (client ClientV2) PatchLayoutExpectInternalServerError(t testing.TB, id string) {
	t.Helper()

	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	AssertNoError(t, err)

	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 500)
}

func (client ClientV2) PatchLayout(t testing.TB, id string, layout Layout) Layout {
	t.Helper()

	up := encode(t, layout)
	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, up)
	AssertNoError(t, err)

	channel := make(chan struct {
		*http.Response
		error
	})
	go func() {
		response, err := (&http.Client{}).Do(request)
		channel <- struct {
			*http.Response
			error
		}{response, err}
	}()

	select {
	case result := <-channel:
		AssertNoError(t, result.error)
		AssertStatusCode(t, result.Response, 200)
		return decode(t, result.Response)
	case <-time.After(time.Second):
		t.Error("timed out after one second")
		return Layout{}
	}
}

func (client ClientV2) DeleteLayout(t testing.TB, id string) {
	t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 204)
}

func (client ClientV2) DeleteLayoutExpectInternalServerError(t testing.TB, id string) {
	t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 500)
}

func (client ClientV2) PutLayerExpectNotFound(t testing.TB, id, name string) {
	t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodPut, requestURL, nil)
	AssertNoError(t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 404)
}

func (client ClientV2) PutLayerExpectBadRequest(t testing.TB, id, name string) {
	t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodPut, requestURL, nil)
	AssertNoError(t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 400)
}

func (client ClientV2) PutLayer(t testing.TB, id, name string, reader io.Reader) {
	t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodPut, requestURL, reader)
	AssertNoError(t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)
}

func (client ClientV2) GetLayerExpectNotFound(t testing.TB, id, name string) {
	t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s/%s", id, name))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 404)
}

func (client ClientV2) GetLayer(t testing.TB, id, name string) (io.ReadCloser, string) {
	t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s/%s", id, name))
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)
	return response.Body, response.Header.Get("Content-Type")
}

func (client ClientV2) DeleteLayer(t testing.TB, id, name string) {
	t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 204)
}

func (client ClientV2) DeleteLayerExpectInternalServerError(t testing.TB, id, name string) {
	t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 500)
}

func (client ClientV2) baseURLF(path string, a ...interface{}) string {
	return client.BaseURL + fmt.Sprintf(path, a...)
}
