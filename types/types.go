package types

type Metadata struct {
	Id           string `json:"id"`
	Size         Size   `json:"size"`
	Bounds       Bounds `json:"bounds"`
	HeightmapURL string `json:"heightmapURL"`
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
