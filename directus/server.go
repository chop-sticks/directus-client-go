package directus

import "encoding/json"

// ServerHealth reports the health of the Directus instance and its subsystems.
type ServerHealth struct {
	Status    string         `json:"status"`
	ReleaseID string         `json:"releaseId,omitempty"`
	ServiceID string         `json:"serviceId,omitempty"`
	Checks    map[string]any `json:"checks,omitempty"`
}

// ServerHealth retrieves the health status of the server. The response is not
// wrapped in a `data` envelope, so it is decoded directly.
func (c *Client) ServerHealth() (*ServerHealth, error) {
	raw, err := c.requestRaw("GET", "/server/health", nil, nil)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}
	var health ServerHealth
	if err := json.Unmarshal(raw, &health); err != nil {
		return nil, err
	}
	return &health, nil
}

// ServerPing returns the server's ping response ("pong"). The endpoint returns
// plain text, not a `data` envelope.
func (c *Client) ServerPing() (string, error) {
	raw, err := c.requestRaw("GET", "/server/ping", nil, nil)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// ServerInfo retrieves information about the server. The response is wrapped in
// the standard `data` envelope.
func (c *Client) ServerInfo() (map[string]any, error) {
	return request[map[string]any](c, "GET", "/server/info", nil, nil)
}

// ReadOpenAPISpec retrieves the OpenAPI (OAS) specification of the instance. The
// response is a raw JSON object, not wrapped in a `data` envelope.
func (c *Client) ReadOpenAPISpec() (map[string]any, error) {
	raw, err := c.requestRaw("GET", "/server/specs/oas", nil, nil)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}
	var spec map[string]any
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}
	return spec, nil
}

// ReadGraphqlSDL retrieves the GraphQL SDL of the instance as plain text. Pass
// scope "system" for the system GraphQL schema; any other value returns the
// main GraphQL schema.
func (c *Client) ReadGraphqlSDL(scope string) (string, error) {
	path := "/server/specs/graphql"
	if scope == "system" {
		path = "/server/specs/graphql/system"
	}
	raw, err := c.requestRaw("GET", path, nil, nil)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
