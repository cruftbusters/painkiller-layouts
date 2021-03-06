package v1

import (
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func NewLayoutAwaitingLayerWire(
	layoutService LayoutService,
	awaitingHeightmap AwaitingLayerService,
	awaitingHillshade AwaitingLayerService,
) LayoutService {
	return &DefaultLayoutAwaitingLayerWire{
		layoutService,
		awaitingHeightmap,
		awaitingHillshade,
	}
}

type DefaultLayoutAwaitingLayerWire struct {
	layoutService     LayoutService
	awaitingHeightmap AwaitingLayerService
	awaitingHillshade AwaitingLayerService
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

func (s *DefaultLayoutAwaitingLayerWire) Patch(id string, patch Layout) (Layout, error) {
	down, err := s.layoutService.Patch(id, patch)
	if patch.HiResHeightmapURL != "" || patch.Scale != 0 {
		if err := s.awaitingHillshade.Enqueue(down); err != nil {
			return down, err
		}
	}
	return down, err
}

func (s *DefaultLayoutAwaitingLayerWire) Delete(id string) error {
	return s.layoutService.Delete(id)
}
