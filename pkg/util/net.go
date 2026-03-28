package util

import (
	"net"
)

func GetAvailablePort(addrList []string) (string, bool) {
	for _, addr := range addrList {
		if IsPortAvailable(addr) {
			return addr, true
		}
	}
	return "", false
}

func IsPortAvailable(addr string) bool {
	listener, err := net.Listen("tcp", addr)
	if err != nil { // port in use
		return false
	}
	defer listener.Close()
	return true
}
