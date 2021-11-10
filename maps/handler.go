package maps

import "github.com/julienschmidt/httprouter"

func Handler() *httprouter.Router {
	router := httprouter.New()
	mapService := NewService(
		&DefaultUUIDService{},
	)
	MapController{
		mapService,
		DefaultHeightmapService{mapService},
	}.AddRoutes(router)
	return router
}
