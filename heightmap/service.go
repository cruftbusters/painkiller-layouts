package heightmap

import . "github.com/cruftbusters/painkiller-gallery/types"

type Service interface {
	Get(id string) *Metadata
	Post(metadata Metadata) Metadata
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

func (service *DefaultService) Get(id string) *Metadata {
	if metadata, hasKey := service.metadata[id]; hasKey {
		return &metadata
	} else {
		return nil
	}
}

func (service *DefaultService) Post(metadata Metadata) Metadata {
	id := service.uuidService.NewUUID()
	newMetadata := &metadata
	newMetadata.Id = id
	service.metadata[id] = *newMetadata
	return *newMetadata
}

func (service *DefaultService) Delete(id string) error {
	delete(service.metadata, id)
	return nil
}
