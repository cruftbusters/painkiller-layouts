package maps

import (
	"errors"
)

type HeightmapService interface {
	Put(id string, heightmap []byte) error
	Get(id string) ([]byte, error)
}

func NewHeightmapService(mapService Service) HeightmapService {
	return &DefaultHeightmapService{
		mapService,
		make(map[string][]byte),
	}
}

type DefaultHeightmapService struct {
	mapService    Service
	heightmapByID map[string][]byte
}

var HeightmapNotFoundError = errors.New("heightmap not found")

func (s *DefaultHeightmapService) Put(id string, heightmap []byte) error {
	_, err := s.mapService.Get(id)
	if err != nil {
		return err
	}
	s.heightmapByID[id] = heightmap
	return nil
}

func (s *DefaultHeightmapService) Get(id string) ([]byte, error) {
	if _, err := s.mapService.Get(id); err != nil {
		return nil, err
	} else if heightmap := s.heightmapByID[id]; heightmap != nil {
		return heightmap, nil
	}
	return nil, HeightmapNotFoundError
}
