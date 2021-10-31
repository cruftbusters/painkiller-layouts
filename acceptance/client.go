package acceptance

import (
	"bytes"
	"encoding/json"
	"net/http"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type Client struct {
	BaseUrl string
}

func (client Client) getMetadata() (*http.Response, error) {
	return http.Get(client.baseUrl("/v1/heightmaps/deadbeef"))
}

func (client Client) post(metadata Metadata) (*http.Response, error) {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(metadata)
	return http.Post(client.baseUrl("/v1/heightmaps"), "", body)
}

func (client Client) baseUrl(path string) string {
	return client.BaseUrl + path
}
