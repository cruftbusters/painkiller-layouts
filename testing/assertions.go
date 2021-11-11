package testing

import (
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/types"
)

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("got error wanted no error", err)
	}
}

func AssertStatusCode(t testing.TB, response *http.Response, statusCode int) {
	t.Helper()
	if response.StatusCode != statusCode {
		t.Fatalf("got status code %d want %d", response.StatusCode, statusCode)
	}
}

func AssertMetadata(t testing.TB, got, want types.Metadata) {
	t.Helper()
	if got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func AssertAllMetadataUnordered(t testing.TB, got, want []types.Metadata) {
	t.Helper()
	sort.SliceStable(got, func(i, j int) bool {
		return got[i].Id < got[j].Id
	})
	sort.SliceStable(want, func(i, j int) bool {
		return want[i].Id < want[j].Id
	})
	AssertAllMetadata(t, got, want)
}

func AssertAllMetadata(t testing.TB, got, want []types.Metadata) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
