package directus

import "fmt"

// Version represents a Directus content version (directus_versions).
type Version struct {
	ID          string         `json:"id,omitempty"`
	Key         string         `json:"key,omitempty"`
	Name        string         `json:"name,omitempty"`
	Collection  any            `json:"collection,omitempty"`
	Item        string         `json:"item,omitempty"`
	Hash        string         `json:"hash,omitempty"`
	DateCreated string         `json:"date_created,omitempty"`
	DateUpdated string         `json:"date_updated,omitempty"`
	UserCreated any            `json:"user_created,omitempty"`
	UserUpdated any            `json:"user_updated,omitempty"`
	Delta       map[string]any `json:"delta,omitempty"`
}

// VersionCompare is the result of comparing a content version against main.
type VersionCompare struct {
	Outdated bool           `json:"outdated"`
	MainHash string         `json:"mainHash,omitempty"`
	Current  map[string]any `json:"current,omitempty"`
	Main     map[string]any `json:"main,omitempty"`
}

// GetContentVersions lists content versions.
func (c *Client) GetContentVersions(q *Query) ([]Version, error) {
	return request[[]Version](c, "GET", "/versions", q, nil)
}

// GetContentVersion retrieves a single content version by ID.
func (c *Client) GetContentVersion(id string, q *Query) (*Version, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Version](c, "GET", fmt.Sprintf("/versions/%s", id), q, nil)
}

// CreateContentVersion creates a single content version.
func (c *Client) CreateContentVersion(item *Version, q *Query) (*Version, error) {
	return request[*Version](c, "POST", "/versions", q, item)
}

// CreateContentVersions creates multiple content versions.
func (c *Client) CreateContentVersions(items []Version, q *Query) ([]Version, error) {
	return request[[]Version](c, "POST", "/versions", q, items)
}

// PatchContentVersion updates a single content version by ID.
func (c *Client) PatchContentVersion(id string, item *Version, q *Query) (*Version, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Version](c, "PATCH", fmt.Sprintf("/versions/%s", id), q, item)
}

// PatchContentVersions updates multiple content versions selected by keys,
// applying the same data to each.
func (c *Client) PatchContentVersions(keys []string, item *Version, q *Query) ([]Version, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Version](c, "PATCH", "/versions", q, body)
}

// PatchContentVersionsBatch updates multiple content versions from an array of items.
func (c *Client) PatchContentVersionsBatch(items []Version, q *Query) ([]Version, error) {
	return request[[]Version](c, "PATCH", "/versions", q, items)
}

// DeleteContentVersion deletes a single content version by ID.
func (c *Client) DeleteContentVersion(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/versions/%s", id), nil, nil)
}

// DeleteContentVersions deletes multiple content versions by keys.
func (c *Client) DeleteContentVersions(keys []string) error {
	return c.execute("DELETE", "/versions", nil, keys)
}

// SaveToContentVersion saves item state into a content version and returns the
// resulting item state.
func (c *Client) SaveToContentVersion(id string, item map[string]any) (map[string]any, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[map[string]any](c, "POST", fmt.Sprintf("/versions/%s/save", id), nil, item)
}

// CompareContentVersion compares a content version against the main revision.
func (c *Client) CompareContentVersion(id string) (*VersionCompare, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*VersionCompare](c, "GET", fmt.Sprintf("/versions/%s/compare", id), nil, nil)
}

// PromoteContentVersion promotes a content version to the main revision. When
// fields is nil it is omitted from the request body.
func (c *Client) PromoteContentVersion(id, mainHash string, fields []string) (any, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	body := map[string]any{"mainHash": mainHash}
	if fields != nil {
		body["fields"] = fields
	}
	return request[any](c, "POST", fmt.Sprintf("/versions/%s/promote", id), nil, body)
}
