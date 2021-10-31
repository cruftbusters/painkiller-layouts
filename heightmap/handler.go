package heightmap

func Handler() Controller {
	return Controller{
		NewService(
			&DefaultUUIDService{},
		),
	}
}
