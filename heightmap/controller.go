package heightmap

import (
	"encoding/json"
	"net/http"

	. "github.com/cruftbusters/painkiller-gallery/types"
	"github.com/julienschmidt/httprouter"
)

type Controller struct {
	service Service
}

func NewController(service Service) *httprouter.Router {
	c := &Controller{service}
	router := httprouter.New()
	router.POST("/v1/heightmaps", c.Create)
	router.GET("/v1/heightmaps/:id", c.Get)
	router.GET("/v1/heightmaps", c.GetAll)
	router.DELETE("/v1/heightmaps/:id", c.Delete)
	return router
}

func (c Controller) Create(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.WriteHeader(201)
	up := &Metadata{}
	json.NewDecoder(request.Body).Decode(up)
	down := c.service.Post(*up)
	json.NewEncoder(response).Encode(down)
}

func (c Controller) Get(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	metadata := c.service.Get(ps.ByName("id"))
	if metadata == nil {
		response.WriteHeader(404)
	} else {
		response.WriteHeader(200)
		json.NewEncoder(response).Encode(metadata)
	}
}

func (c Controller) GetAll(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	allMetadata := c.service.GetAll()
	json.NewEncoder(response).Encode(allMetadata)
}

func (c Controller) Delete(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	if err := c.service.Delete(ps.ByName("id")); err != nil {
		response.WriteHeader(500)
	}
	response.WriteHeader(204)
}
