package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devetek/d-panel/pkg/duser"
)

// Deprecated login API v0
type loginAPIv0Deprecated struct {
	ID           int    `json:"id"`
	Fullname     string `json:"fullname"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	TokenVersion string `json:"token_version"`
	RegisterFrom string `json:"register_from"`
	AvatarURL    string `json:"avatar_url"`
}

type jsonResponseLogin struct {
	Code   int                  `json:"code"`
	Status string               `json:"status,omitempty"`
	Data   loginAPIv0Deprecated `json:"data,omitempty"`
	Error  string               `json:"error,omitempty"`
}

type jsonResponseUser struct {
	Code   int                      `json:"code"`
	Status string                   `json:"status,omitempty"`
	Data   duser.ResponseForPrivate `json:"data,omitempty"`
	Error  string                   `json:"error,omitempty"`
}

func (c *Client) Login(email, password string) (*jsonResponseLogin, error) {
	url := c.BaseURL + "/api/v0/user/login"
	jsonStr := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response body with json decoder
	var loginStatus = new(jsonResponseLogin)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(loginStatus)
	if err != nil {
		return nil, err
	}

	if loginStatus.Error != "" {
		return nil, fmt.Errorf("%s", loginStatus.Error)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	// read response header
	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		return nil, fmt.Errorf("Failed to set session cookie")
	}

	// parse cookie
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		return nil, fmt.Errorf("no session found")
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
		return nil, err
	}

	return loginStatus, nil
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

	if profile.Error != "" {
		return nil, fmt.Errorf("%s", profile.Error)
	}

	return profile, nil
}
