package v1

import (
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
)

func TestAwaitingLayerService(t *testing.T) {
	service := NewAwaitingLayerService()

	layout := types.Layout{Id: "enqueue this"}
	if err := service.Enqueue(layout); err != nil {
		t.Fatal(err)
	}
	got := service.Dequeue()
	AssertLayout(t, got, layout)
}
