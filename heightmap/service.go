package heightmap

import . "github.com/cruftbusters/painkiller-gallery/types"

type Service interface {
	get() *Metadata
	post()
}

type DefaultService struct {
	metadata *Metadata
}

func (service *DefaultService) get() *Metadata {
	return service.metadata
}

func (service *DefaultService) post() {
	service.metadata = &Metadata{Id: "deadbeef"}
}
