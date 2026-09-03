package directus

import "fmt"

// Preset represents a Directus preset (bookmark or default view state).
type Preset struct {
	ID              int            `json:"id,omitempty"`
	Bookmark        string         `json:"bookmark,omitempty"`
	User            any            `json:"user,omitempty"`
	Role            any            `json:"role,omitempty"`
	Collection      string         `json:"collection,omitempty"`
	Search          string         `json:"search,omitempty"`
	Layout          string         `json:"layout,omitempty"`
	LayoutQuery     map[string]any `json:"layout_query,omitempty"`
	LayoutOptions   map[string]any `json:"layout_options,omitempty"`
	RefreshInterval int            `json:"refresh_interval,omitempty"`
	Filter          map[string]any `json:"filter,omitempty"`
	Icon            string         `json:"icon,omitempty"`
	Color           string         `json:"color,omitempty"`
}

// GetPresets lists all presets.
func (c *Client) GetPresets(q *Query) ([]Preset, error) {
	return request[[]Preset](c, "GET", "/presets", q, nil)
}

// GetPreset retrieves a single preset by ID.
func (c *Client) GetPreset(id int, q *Query) (*Preset, error) {
	return request[*Preset](c, "GET", fmt.Sprintf("/presets/%d", id), q, nil)
}

// CreatePreset creates a single preset.
func (c *Client) CreatePreset(item *Preset, q *Query) (*Preset, error) {
	return request[*Preset](c, "POST", "/presets", q, item)
}

// CreatePresets creates multiple presets.
func (c *Client) CreatePresets(items []Preset, q *Query) ([]Preset, error) {
	return request[[]Preset](c, "POST", "/presets", q, items)
}

// PatchPreset updates a single preset by ID.
func (c *Client) PatchPreset(id int, item *Preset, q *Query) (*Preset, error) {
	return request[*Preset](c, "PATCH", fmt.Sprintf("/presets/%d", id), q, item)
}

// PatchPresets updates multiple presets identified by keys with the same data.
func (c *Client) PatchPresets(keys []int, item *Preset, q *Query) ([]Preset, error) {
	body := struct {
		Keys []int   `json:"keys"`
		Data *Preset `json:"data"`
	}{Keys: keys, Data: item}
	return request[[]Preset](c, "PATCH", "/presets", q, body)
}

// PatchPresetsBatch updates multiple presets from a batch of items.
func (c *Client) PatchPresetsBatch(items []Preset, q *Query) ([]Preset, error) {
	return request[[]Preset](c, "PATCH", "/presets", q, items)
}

// DeletePreset deletes a single preset by ID.
func (c *Client) DeletePreset(id int) error {
	return c.execute("DELETE", fmt.Sprintf("/presets/%d", id), nil, nil)
}

// DeletePresets deletes multiple presets by keys.
func (c *Client) DeletePresets(keys []int) error {
	return c.execute("DELETE", "/presets", nil, keys)
}
