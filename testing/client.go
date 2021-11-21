package testing

import (
	"fmt"
	"io"
	"net/http"
	"testing"

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

func (client ClientV2) CreatePendingRender(t testing.TB, layout Layout) {
	t.Helper()
	up := encode(t, layout)
	response, err := http.Post(client.baseURLF("/"), "", up)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 201)
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

func (client ClientV2) PatchLayout(t testing.TB, id string, layout Layout) Layout {
	t.Helper()

	up := encode(t, layout)
	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, up)
	AssertNoError(t, err)

	response, err := (&http.Client{}).Do(request)
	AssertNoError(t, err)
	AssertStatusCode(t, response, 200)

	return decode(t, response)
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
