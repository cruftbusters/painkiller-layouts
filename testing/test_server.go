package testing

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var overrideBaseURL string

func init() { flag.StringVar(&overrideBaseURL, "overrideBaseURL", "", "override base URL") }

func TestServer(handler func(router *httprouter.Router, sqlite3Connection, httpBaseURL string)) (string, string) {
	if overrideBaseURL != "" {
		if !strings.HasPrefix(overrideBaseURL, "http://") && !strings.HasPrefix(overrideBaseURL, "https://") {
			log.Fatalf("Malformed protocol in override base URL: %s", overrideBaseURL)
		}
		httpBaseURL := overrideBaseURL
		wsBaseURL := "ws" + strings.TrimPrefix(overrideBaseURL, "http")
		return httpBaseURL, wsBaseURL
	} else {
		listener, httpBaseURL, wsBaseURL := TestListener()
		router := httprouter.New()
		handler(router, "file::memory:?cache=shared", httpBaseURL)
		go func() { http.Serve(listener, router) }()
		return httpBaseURL, wsBaseURL
	}
}
