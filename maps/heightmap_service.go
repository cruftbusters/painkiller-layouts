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
		make(map[string]string),
	}
}

type DefaultHeightmapService struct {
	mapService    Service
	heightmapByID map[string]string
}

var HeightmapNotFoundError = errors.New("heightmap not found")

func (s *DefaultHeightmapService) Put(id string, heightmap []byte) error {
	_, err := s.mapService.Get(id)
	if err != nil {
		return err
	}
	s.heightmapByID[id] = string(heightmap)
	return nil
}

func (s *DefaultHeightmapService) Get(id string) ([]byte, error) {
	_, err := s.mapService.Get(id)
	if err != nil {
		return nil, err
	}
	heightmap := s.heightmapByID[id]
	if heightmap != "" {
		return []byte(s.heightmapByID[id]), nil
	} else {
		return nil, HeightmapNotFoundError
	}
}
