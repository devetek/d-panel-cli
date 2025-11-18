package helper

import "os"

// Check if file ~/.ssh/authorized_keys exist
func IsAuthorizedKeysExist() bool {
	_, err := os.Stat(os.Getenv("HOME") + "/.ssh/authorized_keys")

	return os.IsNotExist(err)
}
