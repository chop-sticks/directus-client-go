package directus

import "fmt"

// Policy is a Directus access policy (system collection `directus_policies`).
type Policy struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	IPAccess    string `json:"ip_access,omitempty"`
	EnforceTFA  bool   `json:"enforce_tfa,omitempty"`
	AdminAccess bool   `json:"admin_access,omitempty"`
	AppAccess   bool   `json:"app_access,omitempty"`
	Permissions []any  `json:"permissions,omitempty"`
	Users       []any  `json:"users,omitempty"`
	Roles       []any  `json:"roles,omitempty"`
}

// PolicyGlobals reports the aggregate access flags for the current user's
// policies.
type PolicyGlobals struct {
	AppAccess   bool `json:"app_access"`
	AdminAccess bool `json:"admin_access"`
	EnforceTFA  bool `json:"enforce_tfa"`
}

// GetPolicies lists policies.
func (c *Client) GetPolicies(q *Query) ([]Policy, error) {
	return request[[]Policy](c, "GET", "/policies", q, nil)
}

// GetPolicy retrieves a single policy by ID.
func (c *Client) GetPolicy(id string, q *Query) (*Policy, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Policy](c, "GET", fmt.Sprintf("/policies/%s", id), q, nil)
}

// GetPolicyGlobals retrieves the aggregate access flags for the current user.
func (c *Client) GetPolicyGlobals() (*PolicyGlobals, error) {
	return request[*PolicyGlobals](c, "GET", "/policies/me/globals", nil, nil)
}

// CreatePolicy creates a single policy.
func (c *Client) CreatePolicy(item *Policy, q *Query) (*Policy, error) {
	return request[*Policy](c, "POST", "/policies", q, item)
}

// CreatePolicies creates multiple policies.
func (c *Client) CreatePolicies(items []Policy, q *Query) ([]Policy, error) {
	return request[[]Policy](c, "POST", "/policies", q, items)
}

// PatchPolicy updates a single policy by ID.
func (c *Client) PatchPolicy(id string, item *Policy, q *Query) (*Policy, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Policy](c, "PATCH", fmt.Sprintf("/policies/%s", id), q, item)
}

// PatchPolicies updates multiple policies identified by keys with the same data.
func (c *Client) PatchPolicies(keys []string, item *Policy, q *Query) ([]Policy, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Policy](c, "PATCH", "/policies", q, body)
}

// PatchPoliciesBatch updates multiple policies from a batch of items.
func (c *Client) PatchPoliciesBatch(items []Policy, q *Query) ([]Policy, error) {
	return request[[]Policy](c, "PATCH", "/policies", q, items)
}

// DeletePolicy deletes a single policy by ID.
func (c *Client) DeletePolicy(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/policies/%s", id), nil, nil)
}

// DeletePolicies deletes multiple policies identified by keys.
func (c *Client) DeletePolicies(keys []string) error {
	return c.execute("DELETE", "/policies", nil, keys)
}
