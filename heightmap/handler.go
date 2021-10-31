package heightmap

func Handler() Controller {
	return Controller{
		&DefaultService{
			uuidService: &DefaultUUIDService{},
		},
	}
}
