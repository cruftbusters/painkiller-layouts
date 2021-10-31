package acceptance

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type Client struct {
	t       testing.TB
	BaseUrl string
}

func (client Client) GetMetadata(id string) Metadata {
	response, err := http.Get(client.baseUrl("/v1/heightmaps/%s", id))
	assertNoError(client.t, err)
	assertStatusCode(client.t, response, 200)

	return decode(client.t, response)
}

func (client Client) GetMetadataExpectNotFound() {
	response, err := http.Get(client.baseUrl("/v1/heightmaps/deadbeef"))
	assertNoError(client.t, err)
	assertStatusCode(client.t, response, 404)
}

func (client Client) Create(metadata Metadata) Metadata {
	up := encode(client.t, metadata)
	response, err := http.Post(client.baseUrl("/v1/heightmaps"), "", up)
	assertNoError(client.t, err)
	assertStatusCode(client.t, response, 201)

	return decode(client.t, response)
}

func (client Client) baseUrl(path string, a ...interface{}) string {
	return client.BaseUrl + fmt.Sprintf(path, a...)
}
