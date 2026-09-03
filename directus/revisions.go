package directus

import "fmt"

// Revision represents a Directus revision record (directus_revisions). Read-only.
type Revision struct {
	ID         int            `json:"id,omitempty"`
	Activity   any            `json:"activity,omitempty"`
	Collection string         `json:"collection,omitempty"`
	Item       string         `json:"item,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
	Delta      map[string]any `json:"delta,omitempty"`
	Parent     any            `json:"parent,omitempty"`
	Version    any            `json:"version,omitempty"`
}

// GetRevisions lists revision records.
func (c *Client) GetRevisions(q *Query) ([]Revision, error) {
	return request[[]Revision](c, "GET", "/revisions", q, nil)
}

// GetRevision retrieves a single revision record by ID.
func (c *Client) GetRevision(id int, q *Query) (*Revision, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Revision](c, "GET", fmt.Sprintf("/revisions/%d", id), q, nil)
}
