package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a Directus REST API client. HostURL is the instance base URL
// (e.g. "https://example.directus.app"); Token is the bearer token sent as the
// Authorization header on every request; HTTPClient is the underlying HTTP
// client and may be replaced to customize timeouts or transport.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient returns a Client for the given host and token. Both are required;
// the client is configured with a default 10s request timeout.
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

	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading %s response body: %w", req.Header, err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("status: %d, body %s", res.StatusCode, body)
	}

	return body, nil
}

// apiResponse is the standard Directus envelope that wraps every payload in a
// top-level "data" key.
type apiResponse[T any] struct {
	Data T `json:"data"`
}

// request performs an HTTP request with an optional JSON body and query, then
// decodes the Directus `{"data": ...}` envelope into T. An empty response body
// (e.g. 204 No Content) yields the zero value of T.
//
// Use request[[]T] for list endpoints, request[*T] for single-item endpoints
// (nil is returned when data is null), and request[T] for scalar payloads.
func request[T any](c *Client, method, path string, q *Query, body any) (T, error) {
	var zero T
	raw, err := c.requestRaw(method, path, q, body)
	if err != nil {
		return zero, err
	}
	if len(raw) == 0 {
		return zero, nil
	}
	var resp apiResponse[T]
	if err := json.Unmarshal(raw, &resp); err != nil {
		return zero, err
	}
	return resp.Data, nil
}

// requestRaw performs an HTTP request and returns the raw response body. It is
// used for endpoints that do not wrap their payload in a `data` envelope (such
// as /server/health, /server/ping, and the spec endpoints) or return plain
// text. A body of type []byte or io.Reader is sent verbatim; any other value is
// JSON-encoded.
func (c *Client) requestRaw(method, path string, q *Query, body any) ([]byte, error) {
	var reader io.Reader
	switch b := body.(type) {
	case nil:
		// no body
	case io.Reader:
		reader = b
	case []byte:
		reader = bytes.NewReader(b)
	default:
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, buildURL(c.HostURL, path, q), reader)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// execute performs a request whose response body is not needed (DELETE and
// action endpoints), returning only the error.
func (c *Client) execute(method, path string, q *Query, body any) error {
	_, err := c.requestRaw(method, path, q, body)
	return err
}
