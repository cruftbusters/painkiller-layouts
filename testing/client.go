package testing

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type ClientV2 struct {
	t       testing.TB
	baseURL string
}

func NewClientV2(t testing.TB, baseURL string) ClientV2 {
	return ClientV2{t: t, baseURL: baseURL}
}

func (client ClientV2) Get(id string) Metadata {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client ClientV2) GetAll() []Metadata {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decodeAllMetadata(client.t, response)
}

func (client ClientV2) GetExpectNotFound(id string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/maps/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client ClientV2) Create(metadata Metadata) Metadata {
	client.t.Helper()
	up := encode(client.t, metadata)
	response, err := http.Post(client.baseURLF("/v1/maps"), "", up)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 201)

	return decode(client.t, response)
}

func (client ClientV2) Patch(id string, metadata Metadata) Metadata {
	client.t.Helper()

	up := encode(client.t, metadata)
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

func (client ClientV2) PutHeightmap(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/maps/%s/heightmap.jpg", id)
	request, err := http.NewRequest(http.MethodPut, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)
}

func (client ClientV2) baseURLF(path string, a ...interface{}) string {
	return client.baseURL + fmt.Sprintf(path, a...)
}
