package helper

import (
	"os"
	"strings"
)

// func to check if SSH pub key authorized
func IsSSHAuthorized(str string) bool {
	// read file content
	content, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/authorized_keys")
	if err != nil {
		return false
	}

	// check if key exist in file content
	if strings.Contains(string(content), str) {
		return true
	}

	return false
}
