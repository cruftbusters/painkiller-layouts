package main

import (
	"net/http"
	"os"

	. "github.com/cruftbusters/painkiller-gallery/maps"
)

func main() {
	http.ListenAndServe(":8080", Handler(os.Args[1]))
}
