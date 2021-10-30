package heightmap

import "testing"

func TestService(t *testing.T) {
	t.Run("get when missing", func(t *testing.T) {
		assertGet(t, &DefaultService{}, false)
	})

	t.Run("get after post", func(t *testing.T) {
		service := &DefaultService{}
		service.post()
		assertGet(t, service, true)
	})
}

func assertGet(t testing.TB, service Service, want bool) {
	t.Helper()
	got := service.get()
	if got != want {
		t.Fatalf("got %t want %t", got, want)
	}
}
