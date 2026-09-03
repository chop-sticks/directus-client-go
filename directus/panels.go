package directus

import "fmt"

// Panel represents a Directus dashboard panel (directus_panels).
type Panel struct {
	ID          string         `json:"id,omitempty"`
	Dashboard   any            `json:"dashboard,omitempty"`
	Name        string         `json:"name,omitempty"`
	Icon        string         `json:"icon,omitempty"`
	Color       string         `json:"color,omitempty"`
	ShowHeader  bool           `json:"show_header"`
	Note        string         `json:"note,omitempty"`
	Type        string         `json:"type,omitempty"`
	PositionX   int            `json:"position_x,omitempty"`
	PositionY   int            `json:"position_y,omitempty"`
	Width       int            `json:"width,omitempty"`
	Height      int            `json:"height,omitempty"`
	Options     map[string]any `json:"options,omitempty"`
	DateCreated string         `json:"date_created,omitempty"`
	UserCreated any            `json:"user_created,omitempty"`
}

// GetPanels lists all panels.
func (c *Client) GetPanels(q *Query) ([]Panel, error) {
	return request[[]Panel](c, "GET", "/panels", q, nil)
}

// GetPanel retrieves a single panel by ID.
func (c *Client) GetPanel(id string, q *Query) (*Panel, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Panel](c, "GET", fmt.Sprintf("/panels/%s", id), q, nil)
}

// CreatePanel creates a single panel.
func (c *Client) CreatePanel(item *Panel, q *Query) (*Panel, error) {
	return request[*Panel](c, "POST", "/panels", q, item)
}

// CreatePanels creates multiple panels.
func (c *Client) CreatePanels(items []Panel, q *Query) ([]Panel, error) {
	return request[[]Panel](c, "POST", "/panels", q, items)
}

// PatchPanel updates a single panel by ID.
func (c *Client) PatchPanel(id string, item *Panel, q *Query) (*Panel, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Panel](c, "PATCH", fmt.Sprintf("/panels/%s", id), q, item)
}

// PatchPanels updates multiple panels identified by keys with a shared payload.
func (c *Client) PatchPanels(keys []string, item *Panel, q *Query) ([]Panel, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Panel](c, "PATCH", "/panels", q, body)
}

// PatchPanelsBatch updates multiple panels from an array of items.
func (c *Client) PatchPanelsBatch(items []Panel, q *Query) ([]Panel, error) {
	return request[[]Panel](c, "PATCH", "/panels", q, items)
}

// DeletePanel deletes a single panel by ID.
func (c *Client) DeletePanel(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/panels/%s", id), nil, nil)
}

// DeletePanels deletes multiple panels by their keys.
func (c *Client) DeletePanels(keys []string) error {
	return c.execute("DELETE", "/panels", nil, keys)
}
