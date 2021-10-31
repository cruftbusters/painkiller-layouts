package heightmap

import (
	"encoding/json"
	"net/http"

	. "github.com/cruftbusters/painkiller-gallery/types"
)

type Controller struct {
	Service Service
}

func (controller Controller) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		response.WriteHeader(201)
		json.NewEncoder(response).Encode(controller.Service.post(Metadata{}))
	} else if metadata := controller.Service.get(); metadata == nil {
		response.WriteHeader(404)
	} else {
		response.WriteHeader(200)
		json.NewEncoder(response).Encode(controller.Service.get())
	}
}
