package types

type Metadata struct {
	Id       string   `json:"id"`
	Size     Size     `json:"size"`
	Position Position `json:"position"`
	Bounds   Bounds   `json:"bounds"`
	ImageURL string   `json:"imageURL"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Position struct {
	Left float64 `json:"left"`
	Top  float64 `json:"top"`
}

type Bounds struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
}
