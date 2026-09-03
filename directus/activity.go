package directus

import "fmt"

// Activity represents a Directus activity record (directus_activity). Read-only.
type Activity struct {
	ID         int    `json:"id,omitempty"`
	Action     string `json:"action,omitempty"`
	User       any    `json:"user,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	IP         string `json:"ip,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	Collection string `json:"collection,omitempty"`
	Item       string `json:"item,omitempty"`
	Origin     string `json:"origin,omitempty"`
	Revisions  []any  `json:"revisions,omitempty"`
}

// GetActivities lists activity records.
func (c *Client) GetActivities(q *Query) ([]Activity, error) {
	return request[[]Activity](c, "GET", "/activity", q, nil)
}

// GetActivity retrieves a single activity record by ID.
func (c *Client) GetActivity(id int, q *Query) (*Activity, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Activity](c, "GET", fmt.Sprintf("/activity/%d", id), q, nil)
}
