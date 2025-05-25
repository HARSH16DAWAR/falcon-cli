package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	tokenEndpoint = "/oauth2/token"
)

// RegionBaseURL maps region codes to their base URLs
var RegionBaseURL = map[string]string{
	"us-1":     "https://api.crowdstrike.com",
	"us-2":     "https://api.us-2.crowdstrike.com",
	"eu-1":     "https://api.eu-1.crowdstrike.com",
	"us-gov-1": "https://api.laggar.gcw.crowdstrike.com",
	"us-gov-2": "https://api.falcon.us-gov-2.crowdstrike.mil",
}

// TokenResponse represents the OAuth2 token response from Falcon
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// TokenManager handles OAuth2 token management
type TokenManager struct {
	token     string
	expiresAt time.Time
	mu        sync.RWMutex
}

var (
	tokenManager *TokenManager
	once         sync.Once
)

// GetTokenManager returns a singleton instance of TokenManager
func GetTokenManager() *TokenManager {
	once.Do(func() {
		tokenManager = &TokenManager{}
	})
	return tokenManager
}

// GetToken returns a valid bearer token, refreshing if necessary
func (tm *TokenManager) GetToken() (string, error) {
	tm.mu.RLock()
	if tm.isTokenValid() {
		token := tm.token
		tm.mu.RUnlock()
		return token, nil
	}
	tm.mu.RUnlock()

	// Token is invalid or expired, get a new one
	return tm.refreshToken()
}

// isTokenValid checks if the current token is still valid
func (tm *TokenManager) isTokenValid() bool {
	return tm.token != "" && time.Now().Before(tm.expiresAt)
}

// refreshToken gets a new token from the Falcon API
func (tm *TokenManager) refreshToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Double check if token was refreshed by another goroutine
	if tm.isTokenValid() {
		return tm.token, nil
	}

	clientID := viper.GetString("falcon.client_id")
	clientSecret := viper.GetString("falcon.client_secret")
	cloudRegion := viper.GetString("falcon.cloud_region")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("Falcon credentials not found. Please run 'falcon-cli init' first")
	}

	// Get the base URL for the region
	baseURL, ok := RegionBaseURL[cloudRegion]
	if !ok {
		return "", fmt.Errorf("invalid cloud region: %s", cloudRegion)
	}

	// Construct the token request URL
	tokenURL := fmt.Sprintf("%s/oauth2/token", baseURL)

	// Create form data
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")

	// Create the request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating token request: %v", err)
	}

	// Add required headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting token: status code %d", resp.StatusCode)
	}

	// Parse the response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error parsing token response: %v", err)
	}

	// Update the token manager
	tm.token = tokenResp.AccessToken
	tm.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // Subtract 60s for safety margin

	return tm.token, nil
}

// GetBearerToken gets a new bearer token from the Falcon API
func GetBearerToken() (string, error) {
	clientID := viper.GetString("falcon.client_id")
	clientSecret := viper.GetString("falcon.client_secret")
	cloudRegion := viper.GetString("falcon.cloud_region")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("Falcon credentials not found. Please run 'falcon-cli init' first")
	}

	// Get the base URL for the region
	baseURL, ok := RegionBaseURL[cloudRegion]
	if !ok {
		return "", fmt.Errorf("invalid cloud region: %s", cloudRegion)
	}

	// Construct the token request URL
	tokenURL := fmt.Sprintf("%s/oauth2/token", baseURL)

	// Create payload with client credentials
	payload := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s", clientID, clientSecret))

	// Create the request
	req, err := http.NewRequest("POST", tokenURL, payload)
	if err != nil {
		return "", fmt.Errorf("error creating token request: %v", err)
	}

	// Add required headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer null")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting token: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful response (201 Created is expected for token creation)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error getting token: status code %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error parsing token response: %v", err)
	}

	return tokenResp.AccessToken, nil
}
