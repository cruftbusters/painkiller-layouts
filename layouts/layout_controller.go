package layouts

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/julienschmidt/httprouter"
)

type LayoutController struct {
	mapService       LayoutService
	heightmapService HeightmapService
}

func (c LayoutController) AddRoutes(router *httprouter.Router) {
	router.POST("/v1/maps", c.Create)
	router.POST("/v1/layouts", c.Create)
	router.GET("/v1/maps/:id", c.Get)
	router.GET("/v1/layouts/:id", c.Get)
	router.GET("/v1/maps", c.GetAll)
	router.GET("/v1/layouts", c.GetAll)
	router.PATCH("/v1/maps/:id", c.Patch)
	router.PATCH("/v1/layouts/:id", c.Patch)
	router.DELETE("/v1/maps/:id", c.Delete)
	router.DELETE("/v1/layouts/:id", c.Delete)
	router.PUT("/v1/maps/:id/heightmap.jpg", c.PutHeightmap)
	router.PUT("/v1/layouts/:id/heightmap.jpg", c.PutHeightmap)
	router.GET("/v1/maps/:id/heightmap.jpg", c.GetHeightmap)
	router.GET("/v1/layouts/:id/heightmap.jpg", c.GetHeightmap)
}

func (c LayoutController) Create(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.WriteHeader(201)
	up := &types.Layout{}
	json.NewDecoder(request.Body).Decode(up)
	down := c.mapService.Create(*up)
	json.NewEncoder(response).Encode(down)
}

func (c LayoutController) Get(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	layout, err := c.mapService.Get(ps.ByName("id"))
	if err != nil {
		response.WriteHeader(404)
	} else {
		response.WriteHeader(200)
		json.NewEncoder(response).Encode(layout)
	}
}

func (c LayoutController) GetAll(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	excludeMapsWithHeightmap := request.URL.Query().Get("excludeLayoutsWithHeightmap") == "true" ||
		request.URL.Query().Get("excludeMapsWithHeightmap") == "true"
	layouts := c.mapService.GetAll(excludeMapsWithHeightmap)
	json.NewEncoder(response).Encode(layouts)
}

func (c LayoutController) Patch(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	up := &types.Layout{}
	json.NewDecoder(request.Body).Decode(up)
	down, err := c.mapService.Patch(ps.ByName("id"), *up)
	if err != nil {
		response.WriteHeader(404)
	} else {
		json.NewEncoder(response).Encode(down)
	}
}

func (c LayoutController) Delete(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	if err := c.mapService.Delete(ps.ByName("id")); err != nil {
		response.WriteHeader(500)
	}
	response.WriteHeader(204)
}

func (c LayoutController) PutHeightmap(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	heightmap, _ := io.ReadAll(request.Body)
	if c.heightmapService.Put(ps.ByName("id"), heightmap) != nil {
		response.WriteHeader(404)
	}
}

func (c LayoutController) GetHeightmap(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	heightmap, contentType, err := c.heightmapService.Get(ps.ByName("id"))
	if err != nil {
		response.WriteHeader(404)
	}
	response.Header().Add("Content-Type", contentType)
	response.Write(heightmap)
}
