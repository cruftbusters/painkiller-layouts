package acceptance

import (
	"net/http"
)

type Client struct {
	BaseUrl string
}

func (client Client) getMetadata() (*http.Response, error) {
	return http.Get(client.baseUrl("/v1/heightmaps/deadbeef"))
}

func (client Client) postMetadata() (*http.Response, error) {
	return http.Post(client.baseUrl("/v1/heightmaps"), "", nil)
}

func (client Client) baseUrl(path string) string {
	return client.BaseUrl + path
}
