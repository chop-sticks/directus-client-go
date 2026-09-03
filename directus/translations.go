package directus

import "fmt"

// Translation represents a Directus custom translation string.
type Translation struct {
	ID       string `json:"id,omitempty"`
	Language string `json:"language,omitempty"`
	Key      string `json:"key,omitempty"`
	Value    string `json:"value,omitempty"`
}

// GetTranslations lists all translations.
func (c *Client) GetTranslations(q *Query) ([]Translation, error) {
	return request[[]Translation](c, "GET", "/translations", q, nil)
}

// GetTranslation retrieves a single translation by ID.
func (c *Client) GetTranslation(id string, q *Query) (*Translation, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Translation](c, "GET", fmt.Sprintf("/translations/%s", id), q, nil)
}

// CreateTranslation creates a single translation.
func (c *Client) CreateTranslation(item *Translation, q *Query) (*Translation, error) {
	return request[*Translation](c, "POST", "/translations", q, item)
}

// CreateTranslations creates multiple translations.
func (c *Client) CreateTranslations(items []Translation, q *Query) ([]Translation, error) {
	return request[[]Translation](c, "POST", "/translations", q, items)
}

// PatchTranslation updates a single translation by ID.
func (c *Client) PatchTranslation(id string, item *Translation, q *Query) (*Translation, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*Translation](c, "PATCH", fmt.Sprintf("/translations/%s", id), q, item)
}

// PatchTranslations updates multiple translations identified by keys with the same data.
func (c *Client) PatchTranslations(keys []string, item *Translation, q *Query) ([]Translation, error) {
	body := struct {
		Keys []string     `json:"keys"`
		Data *Translation `json:"data"`
	}{Keys: keys, Data: item}
	return request[[]Translation](c, "PATCH", "/translations", q, body)
}

// PatchTranslationsBatch updates multiple translations from a batch of items.
func (c *Client) PatchTranslationsBatch(items []Translation, q *Query) ([]Translation, error) {
	return request[[]Translation](c, "PATCH", "/translations", q, items)
}

// DeleteTranslation deletes a single translation by ID.
func (c *Client) DeleteTranslation(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/translations/%s", id), nil, nil)
}

// DeleteTranslations deletes multiple translations by keys.
func (c *Client) DeleteTranslations(keys []string) error {
	return c.execute("DELETE", "/translations", nil, keys)
}
