package heightmap

import . "github.com/cruftbusters/painkiller-gallery/types"

type Service interface {
	get(id string) *Metadata
	post(metadata Metadata) Metadata
}

type DefaultService struct {
	metadata *Metadata
}

func (service *DefaultService) get(id string) *Metadata {
	return service.metadata
}

func (service *DefaultService) post(metadata Metadata) Metadata {
	newMetadata := &metadata
	newMetadata.Id = "deadbeef"
	service.metadata = newMetadata
	return *newMetadata
}
