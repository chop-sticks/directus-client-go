package directus

import "fmt"

// Permission is a Directus permission (system collection
// `directus_permissions`).
type Permission struct {
	ID          int            `json:"id,omitempty"`
	Policy      any            `json:"policy,omitempty"`
	Collection  string         `json:"collection,omitempty"`
	Action      string         `json:"action,omitempty"`
	Permissions map[string]any `json:"permissions,omitempty"`
	Validation  map[string]any `json:"validation,omitempty"`
	Presets     map[string]any `json:"presets,omitempty"`
	Fields      []string       `json:"fields,omitempty"`
}

// GetPermissions lists permissions.
func (c *Client) GetPermissions(q *Query) ([]Permission, error) {
	return request[[]Permission](c, "GET", "/permissions", q, nil)
}

// GetPermission retrieves a single permission by ID.
func (c *Client) GetPermission(id int, q *Query) (*Permission, error) {
	return request[*Permission](c, "GET", fmt.Sprintf("/permissions/%d", id), q, nil)
}

// GetUserPermissions retrieves the current user's aggregated permissions.
func (c *Client) GetUserPermissions() (map[string]any, error) {
	return request[map[string]any](c, "GET", "/permissions/me", nil, nil)
}

// GetItemPermissions retrieves the current user's permissions for a specific
// collection, optionally scoped to a single item key. An empty key omits the
// key segment.
func (c *Client) GetItemPermissions(collection, key string) (map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	path := fmt.Sprintf("/permissions/me/%s", collection)
	if key != "" {
		path = fmt.Sprintf("%s/%s", path, key)
	}
	return request[map[string]any](c, "GET", path, nil, nil)
}

// CreatePermission creates a single permission.
func (c *Client) CreatePermission(item *Permission, q *Query) (*Permission, error) {
	return request[*Permission](c, "POST", "/permissions", q, item)
}

// CreatePermissions creates multiple permissions.
func (c *Client) CreatePermissions(items []Permission, q *Query) ([]Permission, error) {
	return request[[]Permission](c, "POST", "/permissions", q, items)
}

// PatchPermission updates a single permission by ID.
func (c *Client) PatchPermission(id int, item *Permission, q *Query) (*Permission, error) {
	return request[*Permission](c, "PATCH", fmt.Sprintf("/permissions/%d", id), q, item)
}

// PatchPermissions updates multiple permissions identified by keys with the
// same data.
func (c *Client) PatchPermissions(keys []int, item *Permission, q *Query) ([]Permission, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Permission](c, "PATCH", "/permissions", q, body)
}

// PatchPermissionsBatch updates multiple permissions from a batch of items.
func (c *Client) PatchPermissionsBatch(items []Permission, q *Query) ([]Permission, error) {
	return request[[]Permission](c, "PATCH", "/permissions", q, items)
}

// DeletePermission deletes a single permission by ID.
func (c *Client) DeletePermission(id int) error {
	return c.execute("DELETE", fmt.Sprintf("/permissions/%d", id), nil, nil)
}

// DeletePermissions deletes multiple permissions identified by keys.
func (c *Client) DeletePermissions(keys []int) error {
	return c.execute("DELETE", "/permissions", nil, keys)
}
