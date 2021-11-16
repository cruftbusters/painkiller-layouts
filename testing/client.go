package testing

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
)

type ClientV2 struct {
	t       testing.TB
	baseURL string
}

var overrideBaseURL string

func init() { flag.StringVar(&overrideBaseURL, "overrideBaseURL", "", "override base URL") }

func NewTestClient(
	t testing.TB,
	routerSupplier func(sqlite3Connection, baseURL string) *httprouter.Router,
) (ClientV2, string) {
	if overrideBaseURL == "" {
		listener, baseURL := RandomPortListener()
		router := routerSupplier("file::memory:?cache=shared", baseURL)
		go func() { http.Serve(listener, router) }()
		return ClientV2{t: t, baseURL: baseURL}, baseURL
	} else {
		return ClientV2{t: t, baseURL: overrideBaseURL}, overrideBaseURL
	}
}

func (client ClientV2) GetVersion() Version {
	response, err := http.Get(client.baseURLF("/version"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)
	versionContainer, err := DecodeVersion(response)
	AssertNoError(client.t, err)
	return versionContainer
}

func (client ClientV2) GetLayout(id string) Layout {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client ClientV2) GetLayouts() []Layout {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decodeLayouts(client.t, response)
}

func (client ClientV2) GetLayoutsWithoutHeightmap() []Layout {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts?excludeLayoutsWithHeightmap=true"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decodeLayouts(client.t, response)
}

func (client ClientV2) GetLayoutsWithoutHillshade() []Layout {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts?excludeLayoutsWithHillshade=true"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decodeLayouts(client.t, response)
}

func (client ClientV2) GetLayoutExpectNotFound(id string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) CreateLayout(layout Layout) Layout {
	client.t.Helper()
	up := encode(client.t, layout)
	response, err := http.Post(client.baseURLF("/v1/layouts"), "", up)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 201)

	return decode(client.t, response)
}

func (client ClientV2) PatchLayoutExpectNotFound(id string) {
	client.t.Helper()

	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	AssertNoError(client.t, err)

	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) PatchLayout(id string, layout Layout) Layout {
	client.t.Helper()

	up := encode(client.t, layout)
	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, up)
	AssertNoError(client.t, err)

	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client ClientV2) DeleteLayout(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 204)
}

func (client ClientV2) DeleteLayoutExpectInternalServerError(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 500)
}

func (client ClientV2) PutLayerExpectNotFound(id, name string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodPut, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) PutLayerExpectBadRequest(id, name string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodPut, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 400)
}

func (client ClientV2) PutLayer(id, name string, reader io.Reader) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/layouts/%s/%s", id, name)
	request, err := http.NewRequest(http.MethodPut, requestURL, reader)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)
}

func (client ClientV2) GetLayerExpectNotFound(id, name string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s/%s", id, name))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) GetLayer(id, name string) (io.ReadCloser, string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/layouts/%s/%s", id, name))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)
	return response.Body, response.Header.Get("Content-Type")
}

func (client ClientV2) baseURLF(path string, a ...interface{}) string {
	return client.baseURL + fmt.Sprintf(path, a...)
}
