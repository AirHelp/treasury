package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	// version is the version of this client
	version = "0.0.1"
	// default value for UserAgent header
	defaultUserAgent   = "AirHelp/treasury/v" + version
	defaultAPIVersion  = "v1"
	defaultContentType = "application/json"
	defaultMediaType   = "application/json"
)

// NewClient initializes a new API client for the given endpoint and API version.
// It uses the given http client as transport.
// It also initializes the custom http headers to add to each request.
func NewClient(endpoint string, options Options) (*Client, error) {
	if options.apiVersion == "" {
		options.apiVersion = defaultAPIVersion
	}

	validatedURL, err := url.Parse(endpoint + "/" + options.apiVersion + "/")
	if err != nil {
		return nil, err
	}

	if options.userAgent == "" {
		options.userAgent = defaultUserAgent
	}

	if options.httpClient == nil {
		options.httpClient = http.DefaultClient
	}

	return &Client{
		httpClient: options.httpClient,
		endpoint:   validatedURL,
		userAgent:  options.userAgent,
		version:    options.version,
	}, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.endpoint.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", defaultContentType)
	}
	req.Header.Set("Accept", defaultMediaType)
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}
