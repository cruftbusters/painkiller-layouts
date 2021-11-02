package heightmap

import (
	"encoding/json"
	"net/http"

	. "github.com/cruftbusters/painkiller-gallery/types"
	"github.com/julienschmidt/httprouter"
)

type Controller struct {
	Service Service
}

func NewController(service Service) *httprouter.Router {
	controller := &Controller{service}
	router := httprouter.New()
	router.POST("/v1/heightmaps", controller.Create)
	router.GET("/v1/heightmaps/:id", controller.Get)
	router.DELETE("/v1/heightmaps/:id", controller.Delete)
	return router
}

func (controller Controller) Create(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.WriteHeader(201)
	up := &Metadata{}
	json.NewDecoder(request.Body).Decode(up)
	down := controller.Service.Post(*up)
	json.NewEncoder(response).Encode(down)
}

func (controller Controller) Get(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	metadata := controller.Service.Get(ps.ByName("id"))
	if metadata == nil {
		response.WriteHeader(404)
	} else {
		response.WriteHeader(200)
		json.NewEncoder(response).Encode(metadata)
	}
}

func (controller Controller) Delete(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	if err := controller.Service.Delete(ps.ByName("id")); err != nil {
		response.WriteHeader(500)
	}
	response.WriteHeader(204)
}
