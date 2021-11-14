package main

import (
	"net/http"
	"os"

	"github.com/cruftbusters/painkiller-gallery/layouts"
)

func main() {
	http.ListenAndServe(":8080", layouts.Handler(os.Args[1]))
}
