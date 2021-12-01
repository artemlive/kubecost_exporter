package kubecost_api

import (
	"bytes"
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

type AssetItem interface {
}

func NewApiClient(apiUrl *url.URL, userAgent string) *Client {
	return &Client{
		BaseURL:    apiUrl,
		UserAgent:  userAgent,
		httpClient: new(http.Client),
	}

}

func (c *Client) ListAssets(extraQueryParams []string) (AssetItem, error) {
	req, err := c.newRequest("GET", ListAssetsURI, strings.Join(extraQueryParams, "&"), nil)
	if err != nil {
		return nil, err
	}
	var assets AssetItem
	_, err = c.do(req, &assets)
	return assets, err
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
