package maps

import (
	"encoding/json"
	"net/http"

	. "github.com/cruftbusters/painkiller-gallery/types"
	"github.com/julienschmidt/httprouter"
)

type MapController struct {
	service Service
}

func (c MapController) AddRoutes(router *httprouter.Router) {
	router.POST("/v1/maps", c.Create)
	router.GET("/v1/maps/:id", c.Get)
	router.GET("/v1/maps", c.GetAll)
	router.PATCH("/v1/maps/:id", c.Patch)
	router.DELETE("/v1/maps/:id", c.Delete)
}

func (c MapController) Create(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.WriteHeader(201)
	up := &Metadata{}
	json.NewDecoder(request.Body).Decode(up)
	down := c.service.Post(*up)
	json.NewEncoder(response).Encode(down)
}

func (c MapController) Get(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	metadata, err := c.service.Get(ps.ByName("id"))
	if err != nil {
		response.WriteHeader(404)
	} else {
		response.WriteHeader(200)
		json.NewEncoder(response).Encode(metadata)
	}
}

func (c MapController) GetAll(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	allMetadata := c.service.GetAll()
	json.NewEncoder(response).Encode(allMetadata)
}

func (c MapController) Patch(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	up := &Metadata{}
	json.NewDecoder(request.Body).Decode(up)
	down := c.service.Patch(ps.ByName("id"), *up)
	json.NewEncoder(response).Encode(down)
}

func (c MapController) Delete(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	if err := c.service.Delete(ps.ByName("id")); err != nil {
		response.WriteHeader(500)
	}
	response.WriteHeader(204)
}
