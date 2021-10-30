package heightmap

import (
	"testing"
)

func TestService(t *testing.T) {
	t.Run("get when missing", func(t *testing.T) {
		assertGetIsNil(t, &DefaultService{})
	})

	t.Run("get after post", func(t *testing.T) {
		service := &DefaultService{}
		service.post()
		assertGetIsNotNil(t, service)
	})
}

func assertGetIsNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.get()
	if metadata != nil {
		t.Fatal("got metadata but want nil")
	}
}

func assertGetIsNotNil(t testing.TB, service Service) {
	t.Helper()
	metadata := service.get()
	if metadata == nil {
		t.Fatal("got nil but want metadata", metadata)
	}
}
