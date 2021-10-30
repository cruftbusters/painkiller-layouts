package heightmap

import "net/http"

type Controller struct {
	Service Service
}

func (controller Controller) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		controller.Service.post()
		response.WriteHeader(201)
	} else if controller.Service.get() {
		response.WriteHeader(200)
	} else {
		response.WriteHeader(404)
	}
}
