package testing

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
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
	routerSupplier func(baseURL string) *httprouter.Router,
) (ClientV2, string) {
	if overrideBaseURL == "" {
		listener, baseURL := RandomPortListener()
		go func() { http.Serve(listener, routerSupplier(baseURL)) }()
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

func (client ClientV2) Get(id string) Layout {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client ClientV2) GetAll() []Layout {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decodeLayouts(client.t, response)
}

func (client ClientV2) GetAllWithoutHeightmap() []Layout {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps?excludeMapsWithHeightmap=true"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decodeLayouts(client.t, response)
}

func (client ClientV2) GetExpectNotFound(id string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) Create(layout Layout) Layout {
	client.t.Helper()
	up := encode(client.t, layout)
	response, err := http.Post(client.baseURLF("/v1/maps"), "", up)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 201)

	return decode(client.t, response)
}

func (client ClientV2) PatchExpectNotFound(id string) {
	client.t.Helper()

	requestURL := client.baseURLF("/v1/maps/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	AssertNoError(client.t, err)

	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) Patch(id string, layout Layout) Layout {
	client.t.Helper()

	up := encode(client.t, layout)
	requestURL := client.baseURLF("/v1/maps/%s", id)
	request, err := http.NewRequest(http.MethodPatch, requestURL, up)
	AssertNoError(client.t, err)

	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client ClientV2) Delete(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/maps/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 204)
}

func (client ClientV2) DeleteExpectInternalServerError(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/maps/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 500)
}

func (client ClientV2) PutHeightmapExpectNotFound(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/maps/%s/heightmap.jpg", id)
	request, err := http.NewRequest(http.MethodPut, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) PutHeightmap(id string, heightmap io.Reader) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/maps/%s/heightmap.jpg", id)
	request, err := http.NewRequest(http.MethodPut, requestURL, heightmap)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)
}

func (client ClientV2) GetHeightmapExpectNotFound(id string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps/%s/heightmap.jpg", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) GetHeightmap(id string) (io.ReadCloser, string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps/%s/heightmap.jpg", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)
	return response.Body, response.Header.Get("Content-Type")
}

func (client ClientV2) baseURLF(path string, a ...interface{}) string {
	return client.baseURL + fmt.Sprintf(path, a...)
}
