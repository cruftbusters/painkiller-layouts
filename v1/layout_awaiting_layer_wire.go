package v1

import (
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func NewLayoutAwaitingLayerWire(
	layoutService LayoutService,
	awaitingHeightmap AwaitingLayerService,
) LayoutService {
	return &DefaultLayoutAwaitingLayerWire{
		layoutService,
		awaitingHeightmap,
	}
}

type DefaultLayoutAwaitingLayerWire struct {
	layoutService     LayoutService
	awaitingHeightmap AwaitingLayerService
}

func (s *DefaultLayoutAwaitingLayerWire) Create(layout Layout) Layout {
	down := s.layoutService.Create(layout)
	if err := s.awaitingHeightmap.Enqueue(down); err != nil {
		panic(err)
	}
	return down
}

func (s *DefaultLayoutAwaitingLayerWire) Get(id string) (Layout, error) {
	return s.layoutService.Get(id)
}

func (s *DefaultLayoutAwaitingLayerWire) GetAll() []Layout {
	return s.layoutService.GetAll()
}

func (s *DefaultLayoutAwaitingLayerWire) GetAllWithNoHeightmap() []Layout {
	return s.layoutService.GetAllWithNoHeightmap()
}

func (s *DefaultLayoutAwaitingLayerWire) GetAllWithHeightmapWithoutHillshade() []Layout {
	return s.layoutService.GetAllWithHeightmapWithoutHillshade()
}

func (s *DefaultLayoutAwaitingLayerWire) Patch(id string, patch Layout) (Layout, error) {
	return s.layoutService.Patch(id, patch)
}

func (s *DefaultLayoutAwaitingLayerWire) Delete(id string) error {
	return s.layoutService.Delete(id)
}
