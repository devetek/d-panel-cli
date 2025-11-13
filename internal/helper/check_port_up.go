package helper

import (
	"net"
	"time"
)

func CheckPortUp(port string) bool {
	address := net.JoinHostPort("localhost", port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(5)*time.Second)
	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}
