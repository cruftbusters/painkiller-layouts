package maps

import "github.com/julienschmidt/httprouter"

func Handler() *httprouter.Router {
	router := httprouter.New()
	service := NewService(
		&DefaultUUIDService{},
	)
	Controller{service}.AddRoutes(router)
	MapController{service}.AddRoutes(router)
	return router
}
