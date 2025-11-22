package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devetek/d-panel/pkg/drouter"
)

type jsonResponseRouter struct {
	Code   int                    `json:"code"`
	Status string                 `json:"status,omitempty"`
	Data   drouter.ResponseRouter `json:"data,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// create new router proxy for dPanel agent
func (c *Client) CreateRouter(payload drouter.PayloadRouter) (*jsonResponseRouter, error) {
	// get cookie session
	cookieValue, err := c.readCookieFromFile()
	if err != nil {
		return nil, err
	}

	// convert payload to json
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// http post create secret ssh key
	url := c.BaseURL + "/api/v1/router/create"
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
	var data = new(jsonResponseRouter)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	if data.Code != 200 {
		return nil, fmt.Errorf("data return %d, with error %s", data.Code, data.Error)
	}

	if data.Error != "" {
		return nil, errors.New(data.Error)
	}

	return data, nil
}
