package assertions

import "net"

func RandomPortListener() (net.Listener, int) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	return listener, port
}
