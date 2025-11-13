package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devetek/d-panel/pkg/duser"
)

type jsonResponseUser struct {
	Code   int                      `json:"code"`
	Status string                   `json:"status,omitempty"`
	Data   duser.ResponseForPrivate `json:"data,omitempty"`
	Error  any                      `json:"error,omitempty"`
}

func (c *Client) Login(email, password string) error {
	url := c.BaseURL + "/api/v0/user/login"
	jsonStr := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read response header
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// read response header
	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		return fmt.Errorf("unexpected cookie: %s", cookie)
	}

	// parse cookie
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		return fmt.Errorf("unexpected cookie: %s", cookie)
	}

	// get cookie dcloud_sid
	var cookieValue string
	for _, cookie := range cookies {
		if cookie.Name == "dcloud_sid" {
			cookieValue = cookie.Value
			break
		}
	}

	err = c.writeCookieToFile(cookieValue)
	if err != nil {
		return err
	}

	return nil
}

// fetch API get user profile
func (c *Client) GetProfile() (*jsonResponseUser, error) {
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + "/api/v1/user/profile"
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set cookie to request header, with cookie name dcloud_sid
	req.Header.Set("Cookie", "dcloud_sid="+cookieValue)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response body
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// read response body with json decoder
	var profile = new(jsonResponseUser)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
