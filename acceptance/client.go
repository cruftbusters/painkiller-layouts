package acceptance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type Client struct {
	t       testing.TB
	BaseUrl string
}

func (client Client) GetMetadata() Metadata {
	response, err := http.Get(client.baseUrl("/v1/heightmaps/deadbeef"))
	assertNoError(client.t, err)
	assertStatusCode(client.t, response, 200)

	return decodeMetadata(client.t, response)
}

func (client Client) GetMetadataExpectNotFound() {
	response, err := http.Get(client.baseUrl("/v1/heightmaps/deadbeef"))
	assertNoError(client.t, err)
	assertStatusCode(client.t, response, 404)
}

func (client Client) Create(metadata Metadata) Metadata {
	up := encodeMetadata(client.t, metadata)
	response, err := http.Post(client.baseUrl("/v1/heightmaps"), "", up)
	assertNoError(client.t, err)
	assertStatusCode(client.t, response, 201)

	return decodeMetadata(client.t, response)
}

func (client Client) baseUrl(path string) string {
	return client.BaseUrl + path
}

func encodeMetadata(t testing.TB, metadata Metadata) *bytes.Buffer {
	up := &bytes.Buffer{}
	err := json.NewEncoder(up).Encode(metadata)
	assertNoError(t, err)
	return up
}

func assertNoError(t testing.TB, err error) {
	if err != nil {
		t.Fatal("got error wanted no error", err)
	}
}

func assertStatusCode(t testing.TB, response *http.Response, statusCode int) {
	t.Helper()
	if response.StatusCode != statusCode {
		t.Fatalf("got status code %d want %d", response.StatusCode, statusCode)
	}
}

func decodeMetadata(t testing.TB, response *http.Response) Metadata {
	down := &Metadata{}
	err := json.NewDecoder(response.Body).Decode(down)
	assertNoError(t, err)
	return *down
}
