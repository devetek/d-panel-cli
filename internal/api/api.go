package api

import (
	"fmt"
	"os"
	"path"
)

var (
	BaseURL = "https://pawon.terpusat.com"
)

type Client struct {
	BaseURL string
}

func NewClient() *Client {
	// for development purpose, we can set base URL from env variable
	apiURL := os.Getenv("DPANEL_API_BASE_URL")
	if apiURL != "" {
		BaseURL = apiURL
	}

	client := &Client{
		BaseURL: BaseURL,
	}
	return client
}

// func check session
func (c *Client) CheckSessionExist() error {
	// read cookieValue from file
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return err
	}

	// check if cookieValue is empty
	if cookieValue == "" {
		return fmt.Errorf("cookie is empty")
	}

	return nil
}

// func to write file cookieValue to file
func (c *Client) writeCookieToFile(cookieValue string) error {
	// check if folder .devetek exist in home directory
	if !checkDevetekFolderExist() {
		// create folder .devetek in home directory
		err := createDevetekFolder()
		if err != nil {
			return err
		}
	}

	// write cookieValue to file
	devetekDir, err := getDevetekDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(devetekDir, "session"), []byte(cookieValue), 0644)
	if err != nil {
		return err
	}
	return nil
}

// func to read cookie from file
func (c *Client) readCookieFromFile() (string, error) {
	// read cookieValue from file
	devetekDir, err := getDevetekDir()
	if err != nil {
		return "", err
	}

	cookieValue, err := os.ReadFile(path.Join(devetekDir, "session"))
	if err != nil {
		return "", err
	}
	return string(cookieValue), nil
}
