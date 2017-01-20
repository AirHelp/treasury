package client

import (
	"net/http"
	"net/url"
)

// Client is the API client that performs all operations
// against a treasury server.
type Client struct {
	// Base URL for API requests.
	endpoint *url.URL

	// User agent used when communicating with the GitHub API.
	userAgent string

	// version of the Client
	version string

	// HTTP client used to communicate with the API.
	httpClient *http.Client
}

// Secret contains response of API:
// GET "/secret?key="
type Secret struct {
	Key     string  `json:"key"`
	Value   string  `json:"value"`
	KmsARN  string  `json:"kms_arn"`
	Author  string  `json:"author"`
	Version float32 `json:"version"`
}

// WriteMessage contains write response
// POST "/secret"
type WriteMessage struct {
	Message string `json:"message"`
}

// Options for client
type Options struct {
	version    string
	apiVersion string
	userAgent  string
	httpClient *http.Client
}
