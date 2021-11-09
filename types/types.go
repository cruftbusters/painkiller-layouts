package types

type Metadata struct {
	Id       string   `json:"id"`
	Size     Size     `json:"size"`
	Position Position `json:"position"`
	ImageURL string   `json:"imageURL"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Position struct {
	Top  float64 `json:"top"`
	Left float64 `json:"left"`
}
