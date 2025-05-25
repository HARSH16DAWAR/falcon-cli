package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

// FalconClient represents a client for the Falcon API
type FalconClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

// NewFalconClient creates a new Falcon API client
func NewFalconClient() (*FalconClient, error) {
	// Get cloud region
	cloudRegion := viper.GetString("falcon.cloud_region")

	// Get base URL for the region
	baseURL, ok := RegionBaseURL[cloudRegion]
	if !ok {
		return nil, fmt.Errorf("invalid cloud region: %s", cloudRegion)
	}

	// Get bearer token
	token, err := GetBearerToken()
	if err != nil {
		return nil, fmt.Errorf("error getting bearer token: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &FalconClient{
		BaseURL: baseURL,
		Token:   token,
		Client:  client,
	}, nil
}

// Get makes a GET request to the Falcon API
func (fc *FalconClient) Get(endpoint string, params map[string]string) (*http.Response, error) {
	// Build URL with query parameters
	apiURL := fmt.Sprintf("%s%s", fc.BaseURL, endpoint)
	if len(params) > 0 {
		query := url.Values{}
		for key, value := range params {
			query.Add(key, value)
		}
		apiURL = fmt.Sprintf("%s?%s", apiURL, query.Encode())
	}

	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+fc.Token)

	// Make request
	resp, err := fc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	// Check for successful response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("error: status code %d, body: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// Post makes a POST request to the Falcon API
func (fc *FalconClient) Post(endpoint string, body io.Reader) (*http.Response, error) {
	// Build URL
	apiURL := fmt.Sprintf("%s%s", fc.BaseURL, endpoint)

	// Create request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+fc.Token)

	// Make request
	resp, err := fc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	// Check for successful response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("error: status code %d, body: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// ParseResponse parses the response body into the provided struct
func (fc *FalconClient) ParseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}

	return nil
}
