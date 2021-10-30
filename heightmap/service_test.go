package heightmap

import (
	"reflect"
	"testing"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestService(t *testing.T) {
	t.Run("get when missing", func(t *testing.T) {
		assertGetIsNil(t, &DefaultService{})
	})

	t.Run("get after post", func(t *testing.T) {
		service := &DefaultService{}
		service.post()
		assertGet(t, service, &Metadata{Id: "deadbeef"})
	})
}

func assertGetIsNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.get()
	if metadata != nil {
		t.Fatal("got metadata but want nil")
	}
}

func assertGet(t testing.TB, service Service, want *Metadata) {
	t.Helper()
	got := service.get()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
