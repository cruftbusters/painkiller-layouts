package heightmap

import (
	"encoding/json"
	"net/http"
	"strings"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type Controller struct {
	Service Service
}

func (controller Controller) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	id := strings.TrimPrefix(request.URL.Path, "/v1/heightmaps/")
	if request.Method == "POST" {
		response.WriteHeader(201)
		up := &Metadata{}
		json.NewDecoder(request.Body).Decode(up)
		down := controller.Service.post(*up)
		json.NewEncoder(response).Encode(down)
	} else if request.Method == http.MethodDelete {
		if err := controller.Service.Delete(id); err != nil {
			response.WriteHeader(500)
		}
		response.WriteHeader(204)
	} else {
		metadata := controller.Service.get(id)
		if metadata == nil {
			response.WriteHeader(404)
		} else {
			response.WriteHeader(200)
			json.NewEncoder(response).Encode(metadata)
		}
	}
}
