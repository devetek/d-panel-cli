package api

import (
	"os"
	"path"
)

// get .devetek directory path
func getDevetekDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	devetekDir := path.Join(homeDir, ".devetek")
	return devetekDir, nil
}

// function to check if folder .devetek exist in home directory
func checkDevetekFolderExist() bool {
	devetekDir, err := getDevetekDir()
	if err != nil {
		return false
	}

	if _, err := os.Stat(devetekDir); err == nil {
		return true
	} else {
		return false
	}
}

// create folder .devetek in home directory
func createDevetekFolder() error {
	devetekDir, err := getDevetekDir()
	if err != nil {
		return err
	}

	if err := os.Mkdir(devetekDir, 0755); err != nil {
		return err
	}
	return nil
}
