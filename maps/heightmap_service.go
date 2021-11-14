package maps

import (
	"errors"
	"fmt"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type HeightmapService interface {
	Put(id string, heightmap []byte) error
	Get(id string) ([]byte, string, error)
}

func NewHeightmapService(baseURL string, mapService Service) HeightmapService {
	return &DefaultHeightmapService{
		baseURL,
		mapService,
		make(map[string][]byte),
	}
}

type DefaultHeightmapService struct {
	baseURL       string
	mapService    Service
	heightmapByID map[string][]byte
}

var ErrHeightmapNotFound = errors.New("heightmap not found")

func (s *DefaultHeightmapService) Put(id string, heightmap []byte) error {
	_, err := s.mapService.Get(id)
	if err != nil {
		return err
	}
	s.heightmapByID[id] = heightmap
	heightmapURL := fmt.Sprintf("%s/v1/maps/%s/heightmap.jpg", s.baseURL, id)
	_, err = s.mapService.Patch(id, Metadata{HeightmapURL: heightmapURL})
	return err
}

func (s *DefaultHeightmapService) Get(id string) ([]byte, string, error) {
	if _, err := s.mapService.Get(id); err != nil {
		return nil, "", err
	} else if heightmap := s.heightmapByID[id]; heightmap != nil {
		return heightmap, "image/jpeg", nil
	}
	return nil, "", ErrHeightmapNotFound
}
