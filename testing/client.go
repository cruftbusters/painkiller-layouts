package testing

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type Client struct {
	t       testing.TB
	baseURL string
}

func NewClient(t testing.TB, baseURL string) Client {
	return Client{t: t, baseURL: baseURL}
}

func (client Client) Get(id string) Metadata {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/heightmaps/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client Client) GetAll() []Metadata {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/heightmaps"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decodeAllMetadata(client.t, response)
}

func (client Client) GetExpectNotFound(id string) {
	client.t.Helper()
	response, err := http.Get(client.baseURLF("/v1/heightmaps/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client Client) Create(metadata Metadata) Metadata {
	client.t.Helper()
	up := encode(client.t, metadata)
	response, err := http.Post(client.baseURLF("/v1/heightmaps"), "", up)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 201)

	return decode(client.t, response)
}

func (client Client) Delete(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/heightmaps/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 204)
}

func (client Client) DeleteExpectInternalServerError(id string) {
	client.t.Helper()
	requestURL := client.baseURLF("/v1/heightmaps/%s", id)
	request, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	AssertNoError(client.t, err)
	response, err := (&http.Client{}).Do(request)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 500)
}

func (client Client) baseURLF(path string, a ...interface{}) string {
	return client.baseURL + fmt.Sprintf(path, a...)
}
