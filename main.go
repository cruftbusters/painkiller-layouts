package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cruftbusters/painkiller-layouts/layouts"
)

func main() {
	http.ListenAndServe(
		fmt.Sprintf(":%s", os.Args[1]),
		layouts.Handler(os.Args[2], os.Args[3]),
	)
}
