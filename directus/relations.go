package directus

import "fmt"

// RelationMeta holds the Directus metadata for a relation.
type RelationMeta struct {
	ID                    int      `json:"id,omitempty"`
	JunctionField         string   `json:"junction_field,omitempty"`
	ManyCollection        string   `json:"many_collection,omitempty"`
	ManyField             string   `json:"many_field,omitempty"`
	OneAllowedCollections []string `json:"one_allowed_collections,omitempty"`
	OneCollection         string   `json:"one_collection,omitempty"`
	OneCollectionField    string   `json:"one_collection_field,omitempty"`
	OneDeselectAction     string   `json:"one_deselect_action,omitempty"`
	OneField              string   `json:"one_field,omitempty"`
	SortField             string   `json:"sort_field,omitempty"`
	System                bool     `json:"system,omitempty"`
}

// RelationSchema mirrors the database constraint backing a relation.
type RelationSchema struct {
	Column           string `json:"column,omitempty"`
	ConstraintName   string `json:"constraint_name,omitempty"`
	ForeignKeyColumn string `json:"foreign_key_column,omitempty"`
	ForeignKeySchema string `json:"foreign_key_schema,omitempty"`
	ForeignKeyTable  string `json:"foreign_key_table,omitempty"`
	OnDelete         string `json:"on_delete,omitempty"`
	OnUpdate         string `json:"on_update,omitempty"`
	Table            string `json:"table,omitempty"`
}

// Relation describes a relationship between collections.
type Relation struct {
	Collection        string          `json:"collection,omitempty"`
	Field             string          `json:"field,omitempty"`
	RelatedCollection string          `json:"related_collection,omitempty"`
	Meta              *RelationMeta   `json:"meta,omitempty"`
	Schema            *RelationSchema `json:"schema,omitempty"`
}

// GetRelations lists every relation.
func (c *Client) GetRelations() ([]Relation, error) {
	return request[[]Relation](c, "GET", "/relations", nil, nil)
}

// GetRelationsByCollection lists the relations of a single collection.
func (c *Client) GetRelationsByCollection(collection string) ([]Relation, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[[]Relation](c, "GET", fmt.Sprintf("/relations/%s", collection), nil, nil)
}

// GetRelation retrieves a single relation by collection and field.
func (c *Client) GetRelation(collection string, field string) (*Relation, error) {
	if collection == "" || field == "" {
		return nil, fmt.Errorf("collection and field must be provided")
	}
	return request[*Relation](c, "GET", fmt.Sprintf("/relations/%s/%s", collection, field), nil, nil)
}

// CreateRelation creates a new relation.
func (c *Client) CreateRelation(item *Relation) (*Relation, error) {
	return request[*Relation](c, "POST", "/relations", nil, item)
}

// PatchRelation updates a single relation by collection and field.
func (c *Client) PatchRelation(collection string, field string, item *Relation, q *Query) (*Relation, error) {
	if collection == "" || field == "" {
		return nil, fmt.Errorf("collection and field must be provided")
	}
	return request[*Relation](c, "PATCH", fmt.Sprintf("/relations/%s/%s", collection, field), q, item)
}

// DeleteRelation deletes a relation by collection and field.
func (c *Client) DeleteRelation(collection string, field string) error {
	if collection == "" || field == "" {
		return fmt.Errorf("collection and field must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/relations/%s/%s", collection, field), nil, nil)
}
