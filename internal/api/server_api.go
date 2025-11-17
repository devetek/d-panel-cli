package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devetek/d-panel/pkg/dmachine"
)

type jsonResponseServer struct {
	Code   int                         `json:"code"`
	Status string                      `json:"status,omitempty"`
	Data   dmachine.ResponseForPrivate `json:"data,omitempty"`
	Error  string                      `json:"error,omitempty"`
}

type jsonResponseSetup struct {
	Code   int    `json:"code"`
	Status string `json:"status,omitempty"`
	Data   string `json:"data,omitempty"`
	Error  any    `json:"error,omitempty"`
}

// register new server
func (c *Client) RegisterServer(newServer dmachine.Payload) (*jsonResponseServer, error) {

	// read session from file
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return nil, err
	}

	jsonStr, err := json.Marshal(newServer)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + "/api/v1/server/create"
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonStr)))
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
	var servers = new(jsonResponseServer)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

// setup server
func (c *Client) SetupServer(serverID int) (*jsonResponseSetup, error) {
	// read session from file
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/server/setup/%d", c.BaseURL, serverID)
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("POST", url, nil)
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
	var setup = new(jsonResponseSetup)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(setup)
	if err != nil {
		return nil, err
	}

	return setup, nil
}
