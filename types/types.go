package types

import (
	"encoding/json"
	"io"
	"net/http"
)

type Layout struct {
	Id           string `json:"id"`
	Size         Size   `json:"size"`
	Bounds       Bounds `json:"bounds"`
	HeightmapURL string `json:"heightmapURL"`
	HillshadeURL string `json:"hillshadeURL"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Bounds struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
}

type Version struct {
	Version string `json:"version"`
}

func EncodeVersion(response io.Writer, version Version) error {
	return json.NewEncoder(response).Encode(version)
}

func DecodeVersion(response *http.Response) (Version, error) {
	down := &Version{}
	if err := json.NewDecoder(response.Body).Decode(down); err != nil {
		return Version{}, err
	}
	return *down, nil
}
