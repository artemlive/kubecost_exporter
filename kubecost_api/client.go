package kubecost_api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

const ListAssetsURI = "model/assets"
const AllocationURI = "model/allocation"

func NewApiClient(apiUrl *url.URL, userAgent string, skipTLSVerify bool) *Client {
	// Disable TLS verification globally
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: skipTLSVerify}
	return &Client{
		BaseURL:    apiUrl,
		UserAgent:  userAgent,
		httpClient: new(http.Client),
	}
}

// This method returns interface, because the /model/assets endpoint returns array of different objects
// that can't be mapped here
func (c *Client) ListAssets(extraQueryParams []string) (interface{}, error) {
	req, err := c.newRequest("GET", ListAssetsURI, strings.Join(extraQueryParams, "&"), nil)
	if err != nil {
		return nil, err
	}
	var assets interface{}
	_, err = c.do(req, &assets)
	return assets, err
}

// Method for getting information about namespace costs
// Such a strange response structure: Array with one element, which has a map inside
func (c *Client) GetAllocation(extraQueryParams []string) (*CostDataResponse, error) {
	req, err := c.newRequest("GET", AllocationURI, strings.Join(extraQueryParams, "&"), nil)
	if err != nil {
		return nil, err
	}
	var assets CostDataResponse
	_, err = c.do(req, &assets)
	if err != nil {
		return nil, err
	}
	return &assets, err
}

func (c *Client) newRequest(method, path string, query string, body interface{}) (*http.Request, error) {
	c.BaseURL.Path = path
	c.BaseURL.RawQuery = query
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, c.BaseURL.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
