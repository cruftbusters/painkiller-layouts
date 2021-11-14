package layouts

import "github.com/julienschmidt/httprouter"

func Handler(baseURL string) *httprouter.Router {
	router := httprouter.New()
	layoutService := NewLayoutService(
		&DefaultUUIDService{},
	)
	LayoutController{
		layoutService,
		NewHeightmapService(baseURL, layoutService),
	}.AddRoutes(router)
	VersionController{}.AddRoutes(router)
	return router
}
