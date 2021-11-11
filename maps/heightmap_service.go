package maps

import (
	"errors"
)

type HeightmapService interface {
	Put(id string) error
	Get(id string) (string, error)
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

func (s *DefaultHeightmapService) Put(id string) error {
	_, err := s.mapService.Get(id)
	if err != nil {
		return err
	}
	s.heightmapByID[id] = "fixed bytes"
	return nil
}

func (s *DefaultHeightmapService) Get(id string) (string, error) {
	_, err := s.mapService.Get(id)
	if err != nil {
		return "", err
	}
	heightmap := s.heightmapByID[id]
	if heightmap != "" {
		return s.heightmapByID[id], nil
	} else {
		return "", HeightmapNotFoundError
	}
}
