package heightmap

type UUIDService interface {
	NewUUID() string
}

type DefaultUUIDService struct{}

func (service DefaultUUIDService) NewUUID() string {
	return "deadbeef"
}
