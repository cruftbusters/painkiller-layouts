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
	router.PUT("/v1/layouts/:id/heightmap.jpg", c.Put)
	router.GET("/v1/layouts/:id/heightmap.jpg", c.Get)
}

func (c LayerController) Put(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	layer, _ := io.ReadAll(request.Body)
	if c.layerService.Put(ps.ByName("id"), layer) != nil {
		response.WriteHeader(404)
	}
}

func (c LayerController) Get(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	layer, contentType, err := c.layerService.Get(ps.ByName("id"))
	if err != nil {
		response.WriteHeader(404)
	}
	response.Header().Add("Content-Type", contentType)
	response.Write(layer)
}
