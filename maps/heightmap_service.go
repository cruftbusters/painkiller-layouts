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
		make([]string, 0),
	}
}

type DefaultHeightmapService struct {
	mapService Service
	ids        []string
}

var HeightmapNotFoundError = errors.New("heightmap not found")

func (s *DefaultHeightmapService) Put(id string) error {
	_, err := s.mapService.Get(id)
	if err != nil {
		return err
	}
	s.ids = append(s.ids, id)
	return nil
}

func (s *DefaultHeightmapService) Get(id string) (string, error) {
	_, err := s.mapService.Get(id)
	if err != nil {
		return "", err
	}
	for _, _id := range s.ids {
		if _id == id {
			return "fixed bytes", nil
		}
	}
	return "", HeightmapNotFoundError
}
