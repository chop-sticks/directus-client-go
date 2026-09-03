package directus

import "fmt"

// Folder models a Directus folder (directus_folders). The id is a UUID string.
type Folder struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Parent any    `json:"parent,omitempty"`
}

// GetFolders lists folders.
func (c *Client) GetFolders(q *Query) ([]Folder, error) {
	return request[[]Folder](c, "GET", "/folders", q, nil)
}

// GetFolder retrieves a single folder by id.
func (c *Client) GetFolder(id string, q *Query) (*Folder, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Folder](c, "GET", fmt.Sprintf("/folders/%s", id), q, nil)
}

// CreateFolder creates a single folder.
func (c *Client) CreateFolder(item *Folder, q *Query) (*Folder, error) {
	return request[*Folder](c, "POST", "/folders", q, item)
}

// CreateFolders creates multiple folders.
func (c *Client) CreateFolders(items []Folder, q *Query) ([]Folder, error) {
	return request[[]Folder](c, "POST", "/folders", q, items)
}

// DeleteFolder deletes a single folder by id.
func (c *Client) DeleteFolder(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/folders/%s", id), nil, nil)
}

// DeleteFolders deletes multiple folders identified by keys.
func (c *Client) DeleteFolders(keys []string) error {
	return c.execute("DELETE", "/folders", nil, keys)
}
