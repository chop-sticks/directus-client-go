package directus

import "fmt"

// Notification represents a Directus notification (directus_notifications).
type Notification struct {
	ID         int    `json:"id,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Status     string `json:"status,omitempty"`
	Recipient  any    `json:"recipient,omitempty"`
	Sender     any    `json:"sender,omitempty"`
	Subject    string `json:"subject,omitempty"`
	Message    string `json:"message,omitempty"`
	Collection string `json:"collection,omitempty"`
	Item       string `json:"item,omitempty"`
}

// GetNotifications lists notifications.
func (c *Client) GetNotifications(q *Query) ([]Notification, error) {
	return request[[]Notification](c, "GET", "/notifications", q, nil)
}

// GetNotification retrieves a single notification by ID.
func (c *Client) GetNotification(id int, q *Query) (*Notification, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Notification](c, "GET", fmt.Sprintf("/notifications/%d", id), q, nil)
}

// CreateNotification creates a single notification.
func (c *Client) CreateNotification(item *Notification, q *Query) (*Notification, error) {
	return request[*Notification](c, "POST", "/notifications", q, item)
}

// CreateNotifications creates multiple notifications.
func (c *Client) CreateNotifications(items []Notification, q *Query) ([]Notification, error) {
	return request[[]Notification](c, "POST", "/notifications", q, items)
}

// PatchNotification updates a single notification by ID.
func (c *Client) PatchNotification(id int, item *Notification, q *Query) (*Notification, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Notification](c, "PATCH", fmt.Sprintf("/notifications/%d", id), q, item)
}

// PatchNotifications updates multiple notifications selected by keys, applying
// the same data to each.
func (c *Client) PatchNotifications(keys []int, item *Notification, q *Query) ([]Notification, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]Notification](c, "PATCH", "/notifications", q, body)
}

// PatchNotificationsBatch updates multiple notifications from an array of items.
func (c *Client) PatchNotificationsBatch(items []Notification, q *Query) ([]Notification, error) {
	return request[[]Notification](c, "PATCH", "/notifications", q, items)
}

// DeleteNotification deletes a single notification by ID.
func (c *Client) DeleteNotification(id int) error {
	if id <= 0 {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/notifications/%d", id), nil, nil)
}

// DeleteNotifications deletes multiple notifications by keys.
func (c *Client) DeleteNotifications(keys []int) error {
	return c.execute("DELETE", "/notifications", nil, keys)
}
