package maps

import (
	"errors"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type Service interface {
	Post(metadata Metadata) Metadata
	Get(id string) (Metadata, error)
	GetAll() []Metadata
	Patch(id string, metadata Metadata) Metadata
	Delete(id string) error
}

type DefaultService struct {
	uuidService UUIDService

	metadata map[string]Metadata
}

func NewService(uuidService UUIDService) Service {
	return &DefaultService{
		uuidService: uuidService,
		metadata:    make(map[string]Metadata),
	}
}

var MapNotFoundError = errors.New("map not found")

func (service *DefaultService) Post(metadata Metadata) Metadata {
	id := service.uuidService.NewUUID()
	newMetadata := &metadata
	newMetadata.Id = id
	service.metadata[id] = *newMetadata
	return *newMetadata
}

func (service *DefaultService) Get(id string) (Metadata, error) {
	if metadata, hasKey := service.metadata[id]; hasKey {
		return metadata, nil
	} else {
		return Metadata{}, MapNotFoundError
	}
}

func (service *DefaultService) GetAll() []Metadata {
	all := make([]Metadata, 0, len(service.metadata))
	for _, metadata := range service.metadata {
		all = append(all, metadata)
	}
	return all
}

func (service *DefaultService) Patch(id string, patch Metadata) Metadata {
	oldMetadata := service.metadata[id]
	newMetadata := &oldMetadata
	newMetadata.ImageURL = patch.ImageURL
	service.metadata[id] = *newMetadata
	return *newMetadata
}

func (service *DefaultService) Delete(id string) error {
	delete(service.metadata, id)
	return nil
}
