package tunnel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type latestResponse struct {
	TagName string `json:"tag_name"`
}

func getLatestReleaseVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/devetek/tuman/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release latestResponse
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return release.TagName, nil
}
