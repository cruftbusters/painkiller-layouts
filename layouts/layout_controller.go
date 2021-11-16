package layouts

import (
	"encoding/json"
	"net/http"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
)

type LayoutController struct {
	layoutService LayoutService
}

func (c LayoutController) AddRoutes(router *httprouter.Router) {
	router.POST("/v1/layouts", c.Create)
	router.GET("/v1/layouts/:id", c.Get)
	router.GET("/v1/layouts", c.GetAll)
	router.PATCH("/v1/layouts/:id", c.Patch)
	router.DELETE("/v1/layouts/:id", c.Delete)
}

func (c LayoutController) Create(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.WriteHeader(201)
	up := &types.Layout{}
	json.NewDecoder(request.Body).Decode(up)
	down := c.layoutService.Create(*up)
	json.NewEncoder(response).Encode(down)
}

func (c LayoutController) Get(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	layout, err := c.layoutService.Get(ps.ByName("id"))
	if err != nil {
		response.WriteHeader(404)
	} else {
		response.WriteHeader(200)
		json.NewEncoder(response).Encode(layout)
	}
}

func (c LayoutController) GetAll(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	var layouts []types.Layout
	if request.URL.Query().Get("excludeLayoutsWithHeightmap") == "true" {
		layouts = c.layoutService.GetAllWithNoHeightmap()
	} else {
		layouts = c.layoutService.GetAll()
	}
	json.NewEncoder(response).Encode(layouts)
}

func (c LayoutController) Patch(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	up := &types.Layout{}
	json.NewDecoder(request.Body).Decode(up)
	down, err := c.layoutService.Patch(ps.ByName("id"), *up)
	if err != nil {
		response.WriteHeader(404)
	} else {
		json.NewEncoder(response).Encode(down)
	}
}

func (c LayoutController) Delete(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	if err := c.layoutService.Delete(ps.ByName("id")); err != nil {
		response.WriteHeader(500)
	}
	response.WriteHeader(204)
}
