package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devetek/d-panel/pkg/dsecret"
)

type jsonResponseSecretSSHList struct {
	Code   int                  `json:"code"`
	Status string               `json:"status,omitempty"`
	Data   dsecret.ResponseList `json:"data,omitempty"`
	Error  any                  `json:"error,omitempty"`
}

type jsonResponseSecretSSH struct {
	Code   int              `json:"code"`
	Status string           `json:"status,omitempty"`
	Data   dsecret.Response `json:"data,omitempty"`
	Error  any              `json:"error,omitempty"`
}

// create ssh key
func (c *Client) CreateSecretSSH() (*jsonResponseSecretSSH, error) {
	// get cookie session
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return nil, err
	}

	// set payload
	var payload = dsecret.Payload{
		KeySize:   4096,
		Name:      fmt.Sprintf("cli-ssh-key-%s", time.Now().Format("20060102150405")), // add date time
		Type:      "ssh-key",
		KeyPrefix: "",
	}

	// convert payload to json
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// http post create secret ssh key
	url := c.BaseURL + "/api/v1/secret/ssh-key/create"
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonStr)))
	if err != nil {
		return nil, err
	}

	// set cookie to request header, with cookie name dcloud_sid
	req.Header.Set("Cookie", "dcloud_sid="+cookieValue)

	// do request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response header
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// read response body with json decoder
	var data = new(jsonResponseSecretSSH)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// get list secret ssh
func (c *Client) GetListSecretSSH() (*jsonResponseSecretSSHList, error) {
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + "/api/v1/secret/ssh-key/find"
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

	// read response header
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// read response body with json decoder
	var data = new(jsonResponseSecretSSHList)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// get secret ssh by id
func (c *Client) GetSecretSSHByID(secretID string) (*jsonResponseSecretSSH, error) {
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return nil, err
	}

	// get secret ssh by id
	url := fmt.Sprintf("%s/api/v1/secret/ssh-key/detail/%s", c.BaseURL, secretID)
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

	// read response header
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// read response body with json decoder
	var data = new(jsonResponseSecretSSH)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
