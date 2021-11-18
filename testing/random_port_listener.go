package testing

import (
	"fmt"
	"net"
)

func RandomPortListener() (net.Listener, string) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	return listener, fmt.Sprintf("localhost:%d", port)
}
