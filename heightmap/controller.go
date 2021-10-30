package heightmap

import (
	"encoding/json"
	"net/http"
)

type Controller struct {
	Service Service
}

func (controller Controller) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		controller.Service.post()
		response.WriteHeader(201)
	} else if metadata := controller.Service.get(); metadata == nil {
		response.WriteHeader(404)
	} else {
		response.WriteHeader(200)
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(controller.Service.get())
	}
}
