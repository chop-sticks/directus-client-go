package directus

import "fmt"

// CollectionTranslation is a per-language display name for a collection.
type CollectionTranslation struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
	Singular    string `json:"singular"`
	Plural      string `json:"plural"`
}

// CollectionMeta holds the Directus-managed metadata for a collection
// (directus_collections), covering presentation, archiving, and sorting.
type CollectionMeta struct {
	Collection               string                  `json:"collection,omitempty"`
	Note                     string                  `json:"note,omitempty"`
	Hidden                   bool                    `json:"hidden"`
	Singleton                bool                    `json:"singleton"`
	Icon                     string                  `json:"icon,omitempty"`
	Color                    string                  `json:"color,omitempty"`
	Translations             []CollectionTranslation `json:"translations"`
	DisplayTemplate          string                  `json:"display_template,omitempty"`
	PreviewURL               string                  `json:"preview_url,omitempty"`
	Versioning               bool                    `json:"versioning"`
	AutosaveRevisionInterval int                     `json:"autosave_revision_interval,omitempty"`
	SortField                string                  `json:"sort_field,omitempty"`
	ArchiveField             string                  `json:"archive_field,omitempty"`
	ArchiveValue             string                  `json:"archive_value,omitempty"`
	UnarchiveValue           string                  `json:"unarchive_value,omitempty"`
	ArchiveAppFilter         bool                    `json:"archive_app_filter"`
	ItemDuplicationFields    []string                `json:"item_duplication_fields,omitempty"`
	Accountability           string                  `json:"accountability,omitempty"`
	System                   bool                    `json:"system,omitempty"`
	Sort                     int                     `json:"sort,omitempty"`
	Group                    string                  `json:"group,omitempty"`
	Collapse                 string                  `json:"collapse"`
	Status                   string                  `json:"status,omitempty"`
}

// CollectionSchema mirrors the @directus/schema Table type.
type CollectionSchema struct {
	Name    string `json:"name"`
	Schema  string `json:"schema,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// Collection is the read model for a Directus collection, pairing its
// Directus metadata (Meta) with the underlying database table info (Schema).
type Collection struct {
	Collection string            `json:"collection"`
	Meta       *CollectionMeta   `json:"meta,omitempty"`
	Schema     *CollectionSchema `json:"schema,omitempty"`
}

// CollectionRequest is the write model for POST /collections. Unlike the
// read model (Collection), it accepts a Fields array of initial fields to
// create alongside the collection; the API never returns that key on reads.
type CollectionRequest struct {
	Collection string            `json:"collection"`
	Meta       *CollectionMeta   `json:"meta,omitempty"`
	Schema     *CollectionSchema `json:"schema,omitempty"`
	Fields     []Field           `json:"fields,omitempty"`
}

// GetCollections lists all collections. This endpoint does not accept query parameters.
func (c *Client) GetCollections() ([]Collection, error) {
	return request[[]Collection](c, "GET", "/collections", nil, nil)
}

// GetCollectionByName retrieves a single collection by table name.
func (c *Client) GetCollectionByName(name string) (*Collection, error) {
	return request[*Collection](c, "GET", fmt.Sprintf("/collections/%s", name), nil, nil)
}

// CreateCollection creates a new collection (and its underlying table).
func (c *Client) CreateCollection(req *CollectionRequest, q *Query) (*Collection, error) {
	return request[*Collection](c, "POST", "/collections", q, req)
}

// PatchCollection updates the metadata of an existing collection.
func (c *Client) PatchCollection(name string, req *CollectionRequest, q *Query) (*Collection, error) {
	return request[*Collection](c, "PATCH", fmt.Sprintf("/collections/%s", name), q, req)
}

// PatchCollectionsBatch updates multiple collections in a single request.
func (c *Client) PatchCollectionsBatch(items []CollectionRequest, q *Query) ([]Collection, error) {
	return request[[]Collection](c, "PATCH", "/collections", q, items)
}

// DeleteCollection deletes a collection and its underlying table.
func (c *Client) DeleteCollection(name string) error {
	return c.execute("DELETE", fmt.Sprintf("/collections/%s", name), nil, nil)
}
