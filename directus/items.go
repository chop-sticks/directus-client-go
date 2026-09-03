package directus

import "fmt"

// Items are schema-defined records in non-core (user) collections. Their shape
// is not known at compile time, so payloads and results are modeled as
// map[string]any.

// GetItems lists the items of a collection.
func (c *Client) GetItems(collection string, q *Query) ([]map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[[]map[string]any](c, "GET", fmt.Sprintf("/items/%s", collection), q, nil)
}

// GetItem retrieves a single item from a collection by its primary key.
func (c *Client) GetItem(collection string, key string, q *Query) (map[string]any, error) {
	if collection == "" || key == "" {
		return nil, fmt.Errorf("collection and key must be provided")
	}
	return request[map[string]any](c, "GET", fmt.Sprintf("/items/%s/%s", collection, key), q, nil)
}

// CreateItem creates a single item in a collection.
func (c *Client) CreateItem(collection string, item map[string]any, q *Query) (map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[map[string]any](c, "POST", fmt.Sprintf("/items/%s", collection), q, item)
}

// CreateItems creates multiple items in a collection.
func (c *Client) CreateItems(collection string, items []map[string]any, q *Query) ([]map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[[]map[string]any](c, "POST", fmt.Sprintf("/items/%s", collection), q, items)
}

// PatchItem updates a single item in a collection by its primary key.
func (c *Client) PatchItem(collection string, key string, item map[string]any, q *Query) (map[string]any, error) {
	if collection == "" || key == "" {
		return nil, fmt.Errorf("collection and key must be provided")
	}
	return request[map[string]any](c, "PATCH", fmt.Sprintf("/items/%s/%s", collection, key), q, item)
}

// PatchItems updates multiple items identified by keys with the same data.
func (c *Client) PatchItems(collection string, keys []any, item map[string]any, q *Query) ([]map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	body := map[string]any{"keys": keys, "data": item}
	return request[[]map[string]any](c, "PATCH", fmt.Sprintf("/items/%s", collection), q, body)
}

// PatchItemsBatch updates multiple items from a batch of full item payloads.
func (c *Client) PatchItemsBatch(collection string, items []map[string]any, q *Query) ([]map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[[]map[string]any](c, "PATCH", fmt.Sprintf("/items/%s", collection), q, items)
}

// DeleteItem deletes a single item from a collection by its primary key.
func (c *Client) DeleteItem(collection string, key string) error {
	if collection == "" || key == "" {
		return fmt.Errorf("collection and key must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/items/%s/%s", collection, key), nil, nil)
}

// DeleteItems deletes multiple items from a collection by their primary keys.
func (c *Client) DeleteItems(collection string, keys []any) error {
	if collection == "" {
		return fmt.Errorf("collection must be provided")
	}
	body := map[string]any{"keys": keys}
	return c.execute("DELETE", fmt.Sprintf("/items/%s", collection), nil, body)
}

// DeleteItemsByQuery deletes the items of a collection that match a query.
func (c *Client) DeleteItemsByQuery(collection string, query map[string]any) error {
	if collection == "" {
		return fmt.Errorf("collection must be provided")
	}
	body := map[string]any{"query": query}
	return c.execute("DELETE", fmt.Sprintf("/items/%s", collection), nil, body)
}

// GetSingleton retrieves the single object of a singleton collection.
func (c *Client) GetSingleton(collection string, q *Query) (map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[map[string]any](c, "GET", fmt.Sprintf("/items/%s", collection), q, nil)
}

// PatchSingleton updates the single object of a singleton collection.
func (c *Client) PatchSingleton(collection string, item map[string]any, q *Query) (map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[map[string]any](c, "PATCH", fmt.Sprintf("/items/%s", collection), q, item)
}

// Aggregate runs an aggregation over a collection. The caller sets the
// aggregation functions and optional grouping via q.Aggregate and q.GroupBy.
func (c *Client) Aggregate(collection string, q *Query) ([]map[string]any, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[[]map[string]any](c, "GET", fmt.Sprintf("/items/%s", collection), q, nil)
}
