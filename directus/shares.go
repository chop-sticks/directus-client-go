package directus

import "fmt"

// Share represents a Directus share (a public link granting access to an item).
type Share struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Collection  string `json:"collection,omitempty"`
	Item        string `json:"item,omitempty"`
	Role        any    `json:"role,omitempty"`
	Password    string `json:"password,omitempty"`
	UserCreated any    `json:"user_created,omitempty"`
	DateCreated string `json:"date_created,omitempty"`
	DateStart   string `json:"date_start,omitempty"`
	DateEnd     string `json:"date_end,omitempty"`
	TimesUsed   int    `json:"times_used,omitempty"`
	MaxUses     int    `json:"max_uses,omitempty"`
}

// ShareInfo is the public metadata returned by GET /shares/info/{id}.
type ShareInfo struct {
	ID         string  `json:"id"`
	Collection string  `json:"collection"`
	Item       string  `json:"item"`
	Password   *string `json:"password"`
	DateStart  *string `json:"date_start"`
	DateEnd    *string `json:"date_end"`
	TimesUsed  *int    `json:"times_used"`
	MaxUses    *int    `json:"max_uses"`
}

// GetShares lists all shares.
func (c *Client) GetShares(q *Query) ([]Share, error) {
	return request[[]Share](c, "GET", "/shares", q, nil)
}

// GetShare retrieves a single share by ID.
func (c *Client) GetShare(id string, q *Query) (*Share, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Share](c, "GET", fmt.Sprintf("/shares/%s", id), q, nil)
}

// CreateShare creates a single share.
func (c *Client) CreateShare(item *Share, q *Query) (*Share, error) {
	return request[*Share](c, "POST", "/shares", q, item)
}

// CreateShares creates multiple shares.
func (c *Client) CreateShares(items []Share, q *Query) ([]Share, error) {
	return request[[]Share](c, "POST", "/shares", q, items)
}

// PatchShare updates a single share by ID.
func (c *Client) PatchShare(id string, item *Share, q *Query) (*Share, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Share](c, "PATCH", fmt.Sprintf("/shares/%s", id), q, item)
}

// PatchShares updates multiple shares identified by keys with the same data.
func (c *Client) PatchShares(keys []string, item *Share, q *Query) ([]Share, error) {
	body := struct {
		Keys []string `json:"keys"`
		Data *Share   `json:"data"`
	}{Keys: keys, Data: item}
	return request[[]Share](c, "PATCH", "/shares", q, body)
}

// PatchSharesBatch updates multiple shares from a batch of items.
func (c *Client) PatchSharesBatch(items []Share, q *Query) ([]Share, error) {
	return request[[]Share](c, "PATCH", "/shares", q, items)
}

// DeleteShare deletes a single share by ID.
func (c *Client) DeleteShare(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/shares/%s", id), nil, nil)
}

// DeleteShares deletes multiple shares by keys.
func (c *Client) DeleteShares(keys []string) error {
	return c.execute("DELETE", "/shares", nil, keys)
}

// AuthenticateShare exchanges share credentials for authentication data.
// AuthenticationData is defined in auth.go.
func (c *Client) AuthenticateShare(share, password, mode string) (*AuthenticationData, error) {
	if share == "" {
		return nil, fmt.Errorf("share must be provided")
	}
	body := struct {
		Share    string `json:"share"`
		Password string `json:"password"`
		Mode     string `json:"mode"`
	}{Share: share, Password: password, Mode: mode}
	return request[*AuthenticationData](c, "POST", "/shares/auth", nil, body)
}

// InviteShare sends the given share to the provided email addresses.
func (c *Client) InviteShare(share string, emails []string) error {
	if share == "" {
		return fmt.Errorf("share must be provided")
	}
	body := struct {
		Share  string   `json:"share"`
		Emails []string `json:"emails"`
	}{Share: share, Emails: emails}
	return c.execute("POST", "/shares/invite", nil, body)
}

// ReadShareInfo retrieves the public metadata for a share by ID.
func (c *Client) ReadShareInfo(id string) (*ShareInfo, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*ShareInfo](c, "GET", fmt.Sprintf("/shares/info/%s", id), nil, nil)
}
