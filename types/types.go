package types

type Metadata struct {
	Id       string `json:"id"`
	Size     Size   `json:"size"`
	ImageURL string `json:"imageURL"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
