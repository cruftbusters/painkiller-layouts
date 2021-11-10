package maps

import "github.com/julienschmidt/httprouter"

func Handler() *httprouter.Router {
	router := httprouter.New()
	MapController{
		NewService(
			&DefaultUUIDService{},
		),
	}.AddRoutes(router)
	return router
}
