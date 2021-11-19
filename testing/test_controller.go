package testing

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Controller interface {
	AddRoutes(*httprouter.Router)
}

func TestController(controller Controller) (string, string) {
	listener, httpBaseURL, wsBaseURL := TestListener()
	router := httprouter.New()
	controller.AddRoutes(router)
	go func() { http.Serve(listener, router) }()
	return httpBaseURL, wsBaseURL
}
