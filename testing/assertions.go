package testing

import (
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/types"
)

func AssertNoError(t testing.TB, err error) {
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
