package heightmap

type Service interface {
	get() bool
	post()
}

type DefaultService struct {
	isCreated bool
}

func (service *DefaultService) get() bool {
	return service.isCreated
}

func (service *DefaultService) post() {
	service.isCreated = true
}
