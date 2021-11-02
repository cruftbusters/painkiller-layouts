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

func (client Client) GetMetadata(id string) Metadata {
	response, err := http.Get(client.baseURLF("/v1/heightmaps/%s", id))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client Client) GetMetadataExpectNotFound() {
	response, err := http.Get(client.baseURLF("/v1/heightmaps/deadbeef"))
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 404)
}

func (client Client) Create(metadata Metadata) Metadata {
	up := encode(client.t, metadata)
	response, err := http.Post(client.baseURLF("/v1/heightmaps"), "", up)
	AssertNoError(client.t, err)
	AssertStatusCode(client.t, response, 201)

	return decode(client.t, response)
}

func (client Client) baseURLF(path string, a ...interface{}) string {
	return client.baseURL + fmt.Sprintf(path, a...)
}
