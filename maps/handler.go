package maps

import "github.com/julienschmidt/httprouter"

func Handler(baseURL string) *httprouter.Router {
	router := httprouter.New()
	mapService := NewService(
		&DefaultUUIDService{},
	)
	MapController{
		mapService,
		NewHeightmapService(baseURL, mapService),
	}.AddRoutes(router)
	return router
}
