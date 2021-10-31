package heightmap

import . "github.com/cruftbusters/painkiller-gallery/types"

type Service interface {
	get(id string) *Metadata
	post(metadata Metadata) Metadata
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

func (service *DefaultService) get(id string) *Metadata {
	if metadata, hasKey := service.metadata[id]; hasKey {
		return &metadata
	} else {
		return nil
	}
}

func (service *DefaultService) post(metadata Metadata) Metadata {
	id := service.uuidService.NewUUID()
	newMetadata := &metadata
	newMetadata.Id = id
	service.metadata[id] = *newMetadata
	return *newMetadata
}
