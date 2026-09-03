package directus

import "fmt"

// Role is a Directus role (system collection `directus_roles`).
type Role struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	Parent      any    `json:"parent,omitempty"`
	Children    []any  `json:"children,omitempty"`
	Policies    []any  `json:"policies,omitempty"`
	Users       []any  `json:"users,omitempty"`
}

// GetRoles lists roles.
func (c *Client) GetRoles(q *Query) ([]Role, error) {
	return request[[]Role](c, "GET", "/roles", q, nil)
}

// GetRole retrieves a single role by ID.
func (c *Client) GetRole(id string, q *Query) (*Role, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Role](c, "GET", fmt.Sprintf("/roles/%s", id), q, nil)
}

// GetRolesMe lists the roles of the currently authenticated user.
func (c *Client) GetRolesMe(q *Query) ([]Role, error) {
	return request[[]Role](c, "GET", "/roles/me", q, nil)
}

// CreateRole creates a single role.
func (c *Client) CreateRole(item *Role, q *Query) (*Role, error) {
	return request[*Role](c, "POST", "/roles", q, item)
}

// CreateRoles creates multiple roles.
func (c *Client) CreateRoles(items []Role, q *Query) ([]Role, error) {
	return request[[]Role](c, "POST", "/roles", q, items)
}

// PatchRole updates a single role by ID.
func (c *Client) PatchRole(id string, item *Role, q *Query) (*Role, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Role](c, "PATCH", fmt.Sprintf("/roles/%s", id), q, item)
}

// PatchRoles updates multiple roles identified by keys with the same data.
func (c *Client) PatchRoles(keys []string, item *Role, q *Query) ([]Role, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Role](c, "PATCH", "/roles", q, body)
}

// PatchRolesBatch updates multiple roles from a batch of items.
func (c *Client) PatchRolesBatch(items []Role, q *Query) ([]Role, error) {
	return request[[]Role](c, "PATCH", "/roles", q, items)
}

// DeleteRole deletes a single role by ID.
func (c *Client) DeleteRole(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/roles/%s", id), nil, nil)
}

// DeleteRoles deletes multiple roles identified by keys.
func (c *Client) DeleteRoles(keys []string) error {
	return c.execute("DELETE", "/roles", nil, keys)
}
