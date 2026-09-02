package ebecasv1client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	maxPageSize    = 1000
)

// Client is a client for interacting with the eBECAS V1 API.
type Client struct {
	baseURL     string
	collegeCode string
	username    string
	authToken   string
	httpClient  *http.Client
}

// Config contains the configuration required to create an eBECAS V1 API client.
type Config struct {
	// BaseURL is the eBECAS V1 API URL.
	BaseURL string

	// CollegeCode is the eBECAS college code used for authentication.
	CollegeCode string

	// Username is the eBECAS username used for authentication.
	Username string

	// AuthToken is the eBECAS V1 API authentication token.
	AuthToken string

	// HTTPClient is an optional custom HTTP client.
	// If nil, a client with a 30-second timeout is used.
	HTTPClient *http.Client
}

// NewClient creates a new eBECAS V1 API client using the provided configuration.
func NewClient(config Config) (*Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, errors.New("base URL is required")
	}

	if strings.TrimSpace(config.CollegeCode) == "" {
		return nil, errors.New("college code is required")
	}

	if strings.TrimSpace(config.Username) == "" {
		return nil, errors.New("username is required")
	}

	if strings.TrimSpace(config.AuthToken) == "" {
		return nil, errors.New("auth token is required")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return &Client{
		baseURL:     strings.TrimRight(strings.TrimSpace(config.BaseURL), "/"),
		collegeCode: strings.TrimSpace(config.CollegeCode),
		username:    strings.TrimSpace(config.Username),
		authToken:   strings.TrimSpace(config.AuthToken),
		httpClient:  httpClient,
	}, nil
}

// do executes an authenticated eBECAS V1 API request.
//
// It returns the response body, HTTP status code, and any request or eBECAS V1 API error.
func (c *Client) do(req *http.Request) ([]byte, int, error) {
	req.Header.Set("COLLEGECODE", c.collegeCode)
	req.Header.Set("USERNAME", c.username)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+c.authToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("execute eBECAS V1 API request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, fmt.Errorf("read eBECAS V1 API response body: %w", err)
	}

	return data, res.StatusCode, nil
}
