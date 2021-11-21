package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
)

type LayoutController struct {
	layoutService LayoutService
	newLayouts    chan types.Layout
}

func (c LayoutController) AddRoutes(router *httprouter.Router) {
	router.POST("/v1/layouts", c.Create)
	router.GET("/v1/layouts/:id", c.Get)
	router.GET("/v1/layouts", c.GetAll)
	router.PATCH("/v1/layouts/:id", c.Patch)
	router.DELETE("/v1/layouts/:id", c.Delete)
}

func (c LayoutController) Create(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	up := &types.Layout{}
	json.NewDecoder(request.Body).Decode(up)
	down := c.layoutService.Create(*up)
	select {
	case c.newLayouts <- down:
		response.WriteHeader(201)
		json.NewEncoder(response).Encode(down)
	default:
		response.WriteHeader(500)
		log.Print("new layouts channel is not accepting emissions")
	}
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
	} else if request.URL.Query().Get("withHeightmapWithoutHillshade") == "true" {
		layouts = c.layoutService.GetAllWithHeightmapWithoutHillshade()
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
	} else {
		response.WriteHeader(204)
	}
}
