package maps

type HeightmapService interface {
	put(id string) error
}

type DefaultHeightmapService struct {
	mapService Service
}

func (s DefaultHeightmapService) put(id string) error {
	_, err := s.mapService.Get(id)
	return err
}
