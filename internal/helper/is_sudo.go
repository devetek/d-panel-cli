package helper

import (
	"os/exec"
)

func IsSudo() bool {
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()

	if err != nil {
		return false
	} else {
		return true
	}
}
