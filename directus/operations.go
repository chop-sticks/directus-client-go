package directus

import "fmt"

// Operation represents a Directus flow operation (directus_operations).
type Operation struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Key         string         `json:"key,omitempty"`
	Type        string         `json:"type,omitempty"`
	PositionX   int            `json:"position_x,omitempty"`
	PositionY   int            `json:"position_y,omitempty"`
	Options     map[string]any `json:"options,omitempty"`
	Resolve     any            `json:"resolve,omitempty"`
	Reject      any            `json:"reject,omitempty"`
	Flow        any            `json:"flow,omitempty"`
	DateCreated string         `json:"date_created,omitempty"`
	UserCreated any            `json:"user_created,omitempty"`
}

// GetOperations lists all operations.
func (c *Client) GetOperations(q *Query) ([]Operation, error) {
	return request[[]Operation](c, "GET", "/operations", q, nil)
}

// GetOperation retrieves a single operation by ID.
func (c *Client) GetOperation(id string, q *Query) (*Operation, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Operation](c, "GET", fmt.Sprintf("/operations/%s", id), q, nil)
}

// CreateOperation creates a single operation.
func (c *Client) CreateOperation(item *Operation, q *Query) (*Operation, error) {
	return request[*Operation](c, "POST", "/operations", q, item)
}

// CreateOperations creates multiple operations.
func (c *Client) CreateOperations(items []Operation, q *Query) ([]Operation, error) {
	return request[[]Operation](c, "POST", "/operations", q, items)
}

// PatchOperation updates a single operation by ID.
func (c *Client) PatchOperation(id string, item *Operation, q *Query) (*Operation, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Operation](c, "PATCH", fmt.Sprintf("/operations/%s", id), q, item)
}

// PatchOperations updates multiple operations identified by keys with a shared payload.
func (c *Client) PatchOperations(keys []string, item *Operation, q *Query) ([]Operation, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Operation](c, "PATCH", "/operations", q, body)
}

// PatchOperationsBatch updates multiple operations from an array of items.
func (c *Client) PatchOperationsBatch(items []Operation, q *Query) ([]Operation, error) {
	return request[[]Operation](c, "PATCH", "/operations", q, items)
}

// DeleteOperation deletes a single operation by ID.
func (c *Client) DeleteOperation(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/operations/%s", id), nil, nil)
}

// DeleteOperations deletes multiple operations by their keys.
func (c *Client) DeleteOperations(keys []string) error {
	return c.execute("DELETE", "/operations", nil, keys)
}
