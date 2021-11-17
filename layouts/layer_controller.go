package layouts

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LayerController struct {
	layerService LayerService
}

func (c LayerController) AddRoutes(router *httprouter.Router) {
	router.PUT("/v1/layouts/:id/:name", c.Put)
	router.GET("/v1/layouts/:id/:name", c.Get)
	router.DELETE("/v1/layouts/:id/:name", c.Delete)
}

func (c LayerController) Put(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	if isUnacceptableName(name) {
		response.WriteHeader(400)
		return
	}
	layer, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	if c.layerService.Put(ps.ByName("id"), name, layer) != nil {
		response.WriteHeader(404)
	}
}

func (c LayerController) Get(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	layer, contentType, err := c.layerService.Get(ps.ByName("id"), ps.ByName("name"))
	if err != nil {
		response.WriteHeader(404)
	}
	response.Header().Add("Content-Type", contentType)
	response.Write(layer)
}

func (c LayerController) Delete(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	if err := c.layerService.Delete(ps.ByName("id"), ps.ByName("name")); err != nil {
		response.WriteHeader(500)
	} else {
		response.WriteHeader(204)
	}
}

func isUnacceptableName(name string) bool {
	return name != "heightmap.jpg" && name != "hillshade.jpg"
}
