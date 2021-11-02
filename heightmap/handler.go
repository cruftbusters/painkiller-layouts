package heightmap

import "github.com/julienschmidt/httprouter"

func Handler() *httprouter.Router {
	return NewController(
		NewService(
			&DefaultUUIDService{},
		),
	)
}
