package staketab

import (
	"encoding/json"
	"net/http"
)

// Client interacts with the staketab API
type Client struct {
	endpoint string
	client   *http.Client
}

// NewClient returns a new staketab service client
func NewClient(httpClient *http.Client, endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		client:   httpClient,
	}
}

// NewDefaultClient returns a default archive service client
func NewDefaultClient(endpoint string) *Client {
	return NewClient(http.DefaultClient, endpoint)
}

// GetAllProviders returns all providers
func (c Client) GetAllProviders() (Providers, error) {
	resp, err := c.client.Get(c.endpoint + "/get_all_providers")
	if err != nil {
		return Providers{}, err
	}

	providers := Providers{}
	err = json.NewDecoder(resp.Body).Decode(&providers)
	return providers, err
}
