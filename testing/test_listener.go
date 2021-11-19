package testing

import (
	"fmt"
	"net"
)

func TestListener() (net.Listener, string, string) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	return listener, fmt.Sprintf("http://localhost:%d", port), fmt.Sprintf("ws://localhost:%d", port)
}
