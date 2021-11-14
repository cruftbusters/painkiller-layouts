package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/types"
)

func encode(t testing.TB, layout types.Layout) *bytes.Buffer {
	up := &bytes.Buffer{}
	err := json.NewEncoder(up).Encode(layout)
	AssertNoError(t, err)
	return up
}

func decode(t testing.TB, response *http.Response) types.Layout {
	down := &types.Layout{}
	err := json.NewDecoder(response.Body).Decode(down)
	AssertNoError(t, err)
	return *down
}

func decodeLayouts(t testing.TB, response *http.Response) []types.Layout {
	var down *[]types.Layout = &[]types.Layout{}
	err := json.NewDecoder(response.Body).Decode(down)
	AssertNoError(t, err)
	return *down
}
