package helper

import (
	"net"
	"time"
)

func IsPortUsed(host string, port string) bool {
	timeout := time.Second * 2 // Set a timeout for the connection attempt

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		return false
	} else {
		defer conn.Close() // Close the connection if successful
		return true
	}
}
