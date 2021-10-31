package heightmap

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubService struct {
	t                  testing.TB
	whenGetCalledWith  string
	getWillReturn      *Metadata
	whenPostCalledWith Metadata
	postWillReturn     Metadata
}

func (stub *StubService) get(got string) *Metadata {
	stub.t.Helper()
	want := stub.whenGetCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.getWillReturn
}

func (stub *StubService) post(got Metadata) Metadata {
	stub.t.Helper()
	want := stub.whenPostCalledWith
	if got != want {
		stub.t.Fatalf("got %#v want %#v", got, want)
	}
	return stub.postWillReturn
}

func TestController(t *testing.T) {
	stubService := &StubService{t: t}
	controller := Controller{
		stubService,
	}

	t.Run("get missing heightmap", func(t *testing.T) {
		stubService.whenGetCalledWith = "deadbeef"
		stubService.getWillReturn = nil

		request, _ := http.NewRequest(http.MethodGet, "/v1/heightmaps/deadbeef", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		assertStatusCode(t, response, 404)
	})

	t.Run("create heightmap", func(t *testing.T) {
		up, down := Metadata{Id: "up"}, Metadata{Id: "down"}
		stubService.whenPostCalledWith = up
		stubService.postWillReturn = down

		body := &bytes.Buffer{}
		json.NewEncoder(body).Encode(up)
		request, _ := http.NewRequest(http.MethodPost, "/v1/heightmaps", body)
		response := httptest.NewRecorder()
		controller.ServeHTTP(response, request)

		assertStatusCode(t, response, 201)
		assertBody(t, response, down)
	})

	t.Run("get heightmap", func(t *testing.T) {
		stubService.whenGetCalledWith = "path-id"
		stubService.getWillReturn = &Metadata{Id: "beefdead"}

		request, _ := http.NewRequest(http.MethodGet, "/v1/heightmaps/path-id", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		assertStatusCode(t, response, 200)
		assertBody(t, response, Metadata{Id: "beefdead"})
	})
}

func assertStatusCode(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Fatalf("got status code %d want %d", response.Code, want)
	}
}

func assertBody(t testing.TB, response *httptest.ResponseRecorder, want Metadata) {
	t.Helper()
	got := &Metadata{}
	if err := json.NewDecoder(response.Body).Decode(got); err != nil {
		t.Fatal("got error json decoding body", err)
	}
	if *got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
