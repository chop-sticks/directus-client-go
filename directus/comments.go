package directus

import "fmt"

// Comment represents a Directus comment on an item.
type Comment struct {
	ID          string `json:"id,omitempty"`
	Collection  any    `json:"collection,omitempty"`
	Item        string `json:"item,omitempty"`
	Comment     string `json:"comment,omitempty"`
	DateCreated string `json:"date_created,omitempty"`
	DateUpdated string `json:"date_updated,omitempty"`
	UserCreated any    `json:"user_created,omitempty"`
	UserUpdated any    `json:"user_updated,omitempty"`
}

// GetComments lists all comments.
func (c *Client) GetComments(q *Query) ([]Comment, error) {
	return request[[]Comment](c, "GET", "/comments", q, nil)
}

// GetComment retrieves a single comment by ID.
func (c *Client) GetComment(id string, q *Query) (*Comment, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Comment](c, "GET", fmt.Sprintf("/comments/%s", id), q, nil)
}

// CreateComment creates a single comment.
func (c *Client) CreateComment(item *Comment, q *Query) (*Comment, error) {
	return request[*Comment](c, "POST", "/comments", q, item)
}

// CreateComments creates multiple comments.
func (c *Client) CreateComments(items []Comment, q *Query) ([]Comment, error) {
	return request[[]Comment](c, "POST", "/comments", q, items)
}

// PatchComment updates a single comment by ID.
func (c *Client) PatchComment(id string, item *Comment, q *Query) (*Comment, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Comment](c, "PATCH", fmt.Sprintf("/comments/%s", id), q, item)
}

// DeleteComment deletes a single comment by ID.
func (c *Client) DeleteComment(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/comments/%s", id), nil, nil)
}

// DeleteComments deletes multiple comments by keys. Unlike most resources,
// the comments endpoint expects the keys wrapped in an object.
func (c *Client) DeleteComments(keys []string) error {
	body := struct {
		Keys []string `json:"keys"`
	}{Keys: keys}
	return c.execute("DELETE", "/comments", nil, body)
}
