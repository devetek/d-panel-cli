package helper

import "os"

// function to append ssh key to authorized_keys file
func AppendAuthorizedKey(sshKey string) error {
	// create file if not exist
	_, err := os.Stat(os.Getenv("HOME") + "/.ssh/authorized_keys")
	if os.IsNotExist(err) {
		err = os.MkdirAll(os.Getenv("HOME")+"/.ssh", 0700)
		if err != nil {
			return err
		}
		err = os.WriteFile(os.Getenv("HOME")+"/.ssh/authorized_keys", []byte{}, 0644)
		if err != nil {
			return err
		}
	}

	// read authorized_keys from current user
	authorizedKeys, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/authorized_keys")
	if err != nil {
		return err
	}

	// append ssh key to authorized_keys file
	authorizedKeys = append(authorizedKeys, []byte(sshKey+"\n")...)

	// write authorized_keys file
	err = os.WriteFile(os.Getenv("HOME")+"/.ssh/authorized_keys", authorizedKeys, 0644)
	if err != nil {
		return err
	}

	return nil
}
