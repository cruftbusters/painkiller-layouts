package maps

import "github.com/julienschmidt/httprouter"

func Handler(baseURL string) *httprouter.Router {
	router := httprouter.New()
	mapService := NewMapService(
		&DefaultUUIDService{},
	)
	MapController{
		mapService,
		NewHeightmapService(baseURL, mapService),
	}.AddRoutes(router)
	VersionController{}.AddRoutes(router)
	return router
}
