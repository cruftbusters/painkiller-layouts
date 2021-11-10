package main

import (
	"net/http"

	. "github.com/cruftbusters/painkiller-gallery/maps"
)

func main() {
	http.ListenAndServe(":8080", Handler())
}
