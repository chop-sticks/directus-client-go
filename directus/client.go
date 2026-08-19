package directus

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(host *string, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    *host,
		Token:      *token,
	}
	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token := c.Token
	req.Header.Set("Authorization", "Bearer "+token)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending %s request: %w", req.Header, err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading %s response body: %w", req.Header, err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("status: %d, body %s", res.StatusCode, body)
	}

	return body, nil
}
