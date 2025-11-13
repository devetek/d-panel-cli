package helper

import (
	"encoding/json"
	"net/http"
)

// struct to decode json response
type jsonResponseIP struct {
	IP string `json:"ip"`
}

// func to get my public IP
func GetMyIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read response body with json decoder
	var data = new(jsonResponseIP)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return "", err
	}

	return data.IP, nil
}
