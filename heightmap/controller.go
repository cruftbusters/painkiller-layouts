package heightmap

import "net/http"

func HeightmapController(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		response.WriteHeader(201)
	} else {
		response.WriteHeader(404)
	}
}
