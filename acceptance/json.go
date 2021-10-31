package acceptance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/types"
)

func encode(t testing.TB, metadata types.Metadata) *bytes.Buffer {
	up := &bytes.Buffer{}
	err := json.NewEncoder(up).Encode(metadata)
	assertNoError(t, err)
	return up
}

func decode(t testing.TB, response *http.Response) types.Metadata {
	down := &types.Metadata{}
	err := json.NewDecoder(response.Body).Decode(down)
	assertNoError(t, err)
	return *down
}
