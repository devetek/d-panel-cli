package helper

import (
	"fmt"
	"net"
)

// func to check port from 9000 to 10000
func FindAvailablePort() (int, error) {
	for port := 9000; port < 10000; port++ {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return port, nil
		}
		conn.Close()
	}

	return 0, fmt.Errorf("no available port found")
}
