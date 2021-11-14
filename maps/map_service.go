package maps

import (
	"errors"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type MapService interface {
	Create(metadata Metadata) Metadata
	Get(id string) (Metadata, error)
	GetAll(excludeMapsWithHeightmap bool) []Metadata
	Patch(id string, metadata Metadata) (Metadata, error)
	Delete(id string) error
}

type DefaultMapService struct {
	uuidService UUIDService

	metadata map[string]Metadata
}

func NewMapService(uuidService UUIDService) MapService {
	return &DefaultMapService{
		uuidService: uuidService,
		metadata:    make(map[string]Metadata),
	}
}

var ErrMapNotFound = errors.New("map not found")

func (service *DefaultMapService) Create(metadata Metadata) Metadata {
	id := service.uuidService.NewUUID()
	newMetadata := &metadata
	newMetadata.Id = id
	service.metadata[id] = *newMetadata
	return *newMetadata
}

func (service *DefaultMapService) Get(id string) (Metadata, error) {
	if metadata, hasKey := service.metadata[id]; hasKey {
		return metadata, nil
	} else {
		return Metadata{}, ErrMapNotFound
	}
}

func (service *DefaultMapService) GetAll(excludeMapsWithHeightmap bool) []Metadata {
	all := make([]Metadata, 0, len(service.metadata))
	for _, metadata := range service.metadata {
		if !excludeMapsWithHeightmap || metadata.HeightmapURL == "" {
			all = append(all, metadata)
		}
	}
	return all
}

func (service *DefaultMapService) Patch(id string, patch Metadata) (Metadata, error) {
	if oldMetadata, ok := service.metadata[id]; ok {
		newMetadata := &oldMetadata
		newMetadata.HeightmapURL = patch.HeightmapURL
		service.metadata[id] = *newMetadata
		return *newMetadata, nil
	} else {
		return Metadata{}, ErrMapNotFound
	}
}

func (service *DefaultMapService) Delete(id string) error {
	delete(service.metadata, id)
	return nil
}
