package heightmap

import . "github.com/cruftbusters/painkiller-gallery/types"

type Service interface {
	get(id string) *Metadata
	post(metadata Metadata) Metadata
}

type DefaultService struct {
	uuidService UUIDService
	metadata    *Metadata
}

func (service *DefaultService) get(id string) *Metadata {
	return service.metadata
}

func (service *DefaultService) post(metadata Metadata) Metadata {
	newMetadata := &metadata
	newMetadata.Id = service.uuidService.NewUUID()
	service.metadata = newMetadata
	return *newMetadata
}
