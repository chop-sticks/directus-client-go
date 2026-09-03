package directus

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Flow represents a Directus flow (directus_flows).
type Flow struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Icon           string         `json:"icon,omitempty"`
	Color          string         `json:"color,omitempty"`
	Description    string         `json:"description,omitempty"`
	Status         string         `json:"status,omitempty"`
	Trigger        string         `json:"trigger,omitempty"`
	Accountability string         `json:"accountability,omitempty"`
	Options        map[string]any `json:"options,omitempty"`
	Operation      any            `json:"operation,omitempty"`
	DateCreated    string         `json:"date_created,omitempty"`
	UserCreated    any            `json:"user_created,omitempty"`
}

// GetFlows lists all flows.
func (c *Client) GetFlows(q *Query) ([]Flow, error) {
	return request[[]Flow](c, "GET", "/flows", q, nil)
}

// GetFlow retrieves a single flow by ID.
func (c *Client) GetFlow(id string, q *Query) (*Flow, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Flow](c, "GET", fmt.Sprintf("/flows/%s", id), q, nil)
}

// CreateFlow creates a single flow.
func (c *Client) CreateFlow(item *Flow, q *Query) (*Flow, error) {
	return request[*Flow](c, "POST", "/flows", q, item)
}

// CreateFlows creates multiple flows.
func (c *Client) CreateFlows(items []Flow, q *Query) ([]Flow, error) {
	return request[[]Flow](c, "POST", "/flows", q, items)
}

// PatchFlow updates a single flow by ID.
func (c *Client) PatchFlow(id string, item *Flow, q *Query) (*Flow, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Flow](c, "PATCH", fmt.Sprintf("/flows/%s", id), q, item)
}

// PatchFlows updates multiple flows identified by keys with a shared payload.
func (c *Client) PatchFlows(keys []string, item *Flow, q *Query) ([]Flow, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Flow](c, "PATCH", "/flows", q, body)
}

// PatchFlowsBatch updates multiple flows from an array of items.
func (c *Client) PatchFlowsBatch(items []Flow, q *Query) ([]Flow, error) {
	return request[[]Flow](c, "PATCH", "/flows", q, items)
}

// DeleteFlow deletes a single flow by ID.
func (c *Client) DeleteFlow(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/flows/%s", id), nil, nil)
}

// DeleteFlows deletes multiple flows by their keys.
func (c *Client) DeleteFlows(keys []string) error {
	return c.execute("DELETE", "/flows", nil, keys)
}

// TriggerFlow triggers a flow with a webhook or manual trigger. For a GET
// trigger, data is encoded as query parameters; for a POST trigger, data is
// sent as the JSON body. The raw response body is returned.
func (c *Client) TriggerFlow(method, id string, data map[string]string) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	path := fmt.Sprintf("/flows/trigger/%s", id)
	switch strings.ToUpper(method) {
	case "", http.MethodGet:
		if len(data) > 0 {
			values := url.Values{}
			for k, v := range data {
				values.Set(k, v)
			}
			path += "?" + values.Encode()
		}
		return c.requestRaw(http.MethodGet, path, nil, nil)
	case http.MethodPost:
		return c.requestRaw(http.MethodPost, path, nil, data)
	default:
		return nil, fmt.Errorf("unsupported trigger method %q", method)
	}
}
