package directus

import "fmt"

// Dashboard represents a Directus insights dashboard (directus_dashboards).
type Dashboard struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Note        string `json:"note,omitempty"`
	DateCreated string `json:"date_created,omitempty"`
	UserCreated any    `json:"user_created,omitempty"`
	Color       string `json:"color,omitempty"`
}

// GetDashboards lists all dashboards.
func (c *Client) GetDashboards(q *Query) ([]Dashboard, error) {
	return request[[]Dashboard](c, "GET", "/dashboards", q, nil)
}

// GetDashboard retrieves a single dashboard by ID.
func (c *Client) GetDashboard(id string, q *Query) (*Dashboard, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Dashboard](c, "GET", fmt.Sprintf("/dashboards/%s", id), q, nil)
}

// CreateDashboard creates a single dashboard.
func (c *Client) CreateDashboard(item *Dashboard, q *Query) (*Dashboard, error) {
	return request[*Dashboard](c, "POST", "/dashboards", q, item)
}

// CreateDashboards creates multiple dashboards.
func (c *Client) CreateDashboards(items []Dashboard, q *Query) ([]Dashboard, error) {
	return request[[]Dashboard](c, "POST", "/dashboards", q, items)
}

// PatchDashboard updates a single dashboard by ID.
func (c *Client) PatchDashboard(id string, item *Dashboard, q *Query) (*Dashboard, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Dashboard](c, "PATCH", fmt.Sprintf("/dashboards/%s", id), q, item)
}

// PatchDashboards updates multiple dashboards identified by keys with a shared payload.
func (c *Client) PatchDashboards(keys []string, item *Dashboard, q *Query) ([]Dashboard, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Dashboard](c, "PATCH", "/dashboards", q, body)
}

// PatchDashboardsBatch updates multiple dashboards from an array of items.
func (c *Client) PatchDashboardsBatch(items []Dashboard, q *Query) ([]Dashboard, error) {
	return request[[]Dashboard](c, "PATCH", "/dashboards", q, items)
}

// DeleteDashboard deletes a single dashboard by ID.
func (c *Client) DeleteDashboard(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/dashboards/%s", id), nil, nil)
}

// DeleteDashboards deletes multiple dashboards by their keys.
func (c *Client) DeleteDashboards(keys []string) error {
	return c.execute("DELETE", "/dashboards", nil, keys)
}
