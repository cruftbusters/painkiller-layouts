package main

import (
	"fmt"
	"net/http"
	"os"

	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	v1.Handler(router, os.Args[2], os.Args[3])
	http.ListenAndServe(
		fmt.Sprintf(":%s", os.Args[1]),
		router,
	)
}
