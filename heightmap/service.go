package heightmap

import . "github.com/cruftbusters/painkiller-gallery/types"

type Service interface {
	get() *Metadata
	post(metadata Metadata) Metadata
}

type DefaultService struct {
	metadata *Metadata
}

func (service *DefaultService) get() *Metadata {
	return service.metadata
}

func (service *DefaultService) post(metadata Metadata) Metadata {
	service.metadata = &Metadata{Id: "deadbeef"}
	return *service.metadata
}
