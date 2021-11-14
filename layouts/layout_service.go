package layouts

import (
	"errors"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type LayoutService interface {
	Create(metadata Metadata) Metadata
	Get(id string) (Metadata, error)
	GetAll(excludeMapsWithHeightmap bool) []Metadata
	Patch(id string, metadata Metadata) (Metadata, error)
	Delete(id string) error
}

type DefaultLayoutService struct {
	uuidService UUIDService

	metadata map[string]Metadata
}

func NewLayoutService(uuidService UUIDService) LayoutService {
	return &DefaultLayoutService{
		uuidService: uuidService,
		metadata:    make(map[string]Metadata),
	}
}

var ErrLayoutNotFound = errors.New("layout not found")

func (service *DefaultLayoutService) Create(metadata Metadata) Metadata {
	id := service.uuidService.NewUUID()
	newMetadata := &metadata
	newMetadata.Id = id
	service.metadata[id] = *newMetadata
	return *newMetadata
}

func (service *DefaultLayoutService) Get(id string) (Metadata, error) {
	if metadata, hasKey := service.metadata[id]; hasKey {
		return metadata, nil
	} else {
		return Metadata{}, ErrLayoutNotFound
	}
}

func (service *DefaultLayoutService) GetAll(excludeMapsWithHeightmap bool) []Metadata {
	all := make([]Metadata, 0, len(service.metadata))
	for _, metadata := range service.metadata {
		if !excludeMapsWithHeightmap || metadata.HeightmapURL == "" {
			all = append(all, metadata)
		}
	}
	return all
}

func (service *DefaultLayoutService) Patch(id string, patch Metadata) (Metadata, error) {
	if oldMetadata, ok := service.metadata[id]; ok {
		newMetadata := &oldMetadata
		newMetadata.HeightmapURL = patch.HeightmapURL
		service.metadata[id] = *newMetadata
		return *newMetadata, nil
	} else {
		return Metadata{}, ErrLayoutNotFound
	}
}

func (service *DefaultLayoutService) Delete(id string) error {
	delete(service.metadata, id)
	return nil
}
