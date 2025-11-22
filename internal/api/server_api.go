package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func (c *Client) IsRegistered() bool {
	devetekDir, err := getDevetekDir()
	if err != nil {
		return false
	}

	machineConfig := filepath.Join(devetekDir, "machine.json")

	machineContent, err := os.ReadFile(machineConfig)
	if err != nil {
		return false
	}

	var machine *dmachine.ResponseForPrivate
	err = json.Unmarshal(machineContent, &machine)
	if err != nil {
		return false
	}

	if machine == nil {
		return false
	}

	if machine.ID == 0 {
		return false
	}

	// read session from file
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return false
	}

	// fetch to validate to dPanel
	url := c.BaseURL + "/api/v1/server/detail/" + strconv.FormatInt(int64(machine.GetUint64ID()), 10)
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	// set cookie to request header, with cookie name dcloud_sid
	req.Header.Set("Cookie", "dcloud_sid="+cookieValue)

	resp, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// read response header
	if resp.StatusCode != http.StatusOK {
		return false
	}

	// read response body with json decoder
	var server = new(jsonResponseServer)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(server)
	if err != nil {
		return false
	}

	if server.Error != "" {
		return false
	}

	if server.Code != 200 {
		return false
	}

	if server.Data.ID == 0 {
		return false
	}

	return true
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
