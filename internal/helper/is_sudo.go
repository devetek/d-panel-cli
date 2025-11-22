package helper

import (
	"os/exec"
	"os/user"
)

func IsSudo() bool {
	/**
	*	Todo:
	*	Remove require root as dpanel executor. But it need to update some patchs:
	*	- dpanel-init role
	*	- dpanel-agent role
	 */
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	if currentUser.Username == "root" {
		return true
	}

	cmd := exec.Command("sudo", "-n", "true")
	err = cmd.Run()

	if err != nil {
		return false
	} else {
		return true
	}
}
