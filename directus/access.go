package directus

import "fmt"

// Access is a Directus access record (system collection `directus_access`), the
// junction that assigns a policy to a user or a role.
type Access struct {
	ID     string `json:"id,omitempty"`
	Role   any    `json:"role,omitempty"`
	User   any    `json:"user,omitempty"`
	Policy any    `json:"policy,omitempty"`
	Sort   *int   `json:"sort,omitempty"`
}

// GetAccesses lists access records.
func (c *Client) GetAccesses(q *Query) ([]Access, error) {
	return request[[]Access](c, "GET", "/access", q, nil)
}

// GetAccess retrieves a single access record by ID.
func (c *Client) GetAccess(id string, q *Query) (*Access, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Access](c, "GET", fmt.Sprintf("/access/%s", id), q, nil)
}

// CreateAccess creates a single access record.
func (c *Client) CreateAccess(item *Access, q *Query) (*Access, error) {
	return request[*Access](c, "POST", "/access", q, item)
}

// CreateAccesses creates multiple access records.
func (c *Client) CreateAccesses(items []Access, q *Query) ([]Access, error) {
	return request[[]Access](c, "POST", "/access", q, items)
}

// PatchAccess updates a single access record by ID.
func (c *Client) PatchAccess(id string, item *Access, q *Query) (*Access, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Access](c, "PATCH", fmt.Sprintf("/access/%s", id), q, item)
}

// PatchAccesses updates multiple access records identified by keys with the same data.
func (c *Client) PatchAccesses(keys []string, item *Access, q *Query) ([]Access, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Access](c, "PATCH", "/access", q, body)
}

// PatchAccessesBatch updates multiple access records from a batch of items.
func (c *Client) PatchAccessesBatch(items []Access, q *Query) ([]Access, error) {
	return request[[]Access](c, "PATCH", "/access", q, items)
}

// DeleteAccess deletes a single access record by ID.
func (c *Client) DeleteAccess(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/access/%s", id), nil, nil)
}

// DeleteAccesses deletes multiple access records identified by keys.
func (c *Client) DeleteAccesses(keys []string) error {
	return c.execute("DELETE", "/access", nil, keys)
}
