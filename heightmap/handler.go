package heightmap

import "github.com/julienschmidt/httprouter"

func Handler() *httprouter.Router {
	return Controller{
		NewService(
			&DefaultUUIDService{},
		),
	}.Router()
}
