package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	layouts.Handler(router, os.Args[2], os.Args[3])
	http.ListenAndServe(
		fmt.Sprintf(":%s", os.Args[1]),
		router,
	)
}
