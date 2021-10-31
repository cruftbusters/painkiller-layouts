package heightmap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type StubService struct {
	getWillReturn  *Metadata
	postWillReturn Metadata
}

func (stub *StubService) get() *Metadata {
	return stub.getWillReturn
}

func (stub *StubService) post(metadata Metadata) Metadata {
	return stub.postWillReturn
}

func TestController(t *testing.T) {
	stubService := &StubService{}
	controller := Controller{
		stubService,
	}

	t.Run("get missing heightmap", func(t *testing.T) {
		stubService.getWillReturn = nil

		request, _ := http.NewRequest(http.MethodGet, "/v1/heightmaps/deadbeef", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		assertStatusCode(t, response, 404)
	})

	t.Run("create heightmap", func(t *testing.T) {
		down := Metadata{Id: "down"}
		stubService.postWillReturn = down
		request, _ := http.NewRequest(http.MethodPost, "/v1/heightmaps", nil)
		response := httptest.NewRecorder()

		controller.ServeHTTP(response, request)

		assertStatusCode(t, response, 201)
		assertBody(t, response, down)
	})

	t.Run("get heightmap", func(t *testing.T) {
		stubService.getWillReturn = &Metadata{Id: "beefdead"}

		request, _ := http.NewRequest(http.MethodGet, "/v1/heightmaps/deadbeef", nil)
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
