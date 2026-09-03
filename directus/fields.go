package directus

import "fmt"

// FieldTranslation is a per-language display name for a field.
type FieldTranslation struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

// FieldCondition is a conditional-formatting rule applied to a field in the
// admin app based on the evaluated Rule.
type FieldCondition struct {
	Name                   string         `json:"name"`
	Rule                   map[string]any `json:"rule"`
	Readonly               bool           `json:"readonly,omitempty"`
	Hidden                 bool           `json:"hidden,omitempty"`
	Options                map[string]any `json:"options,omitempty"`
	Required               bool           `json:"required,omitempty"`
	ClearHiddenValueOnSave bool           `json:"clear_hidden_value_on_save,omitempty"`
}

// FieldMeta holds the Directus-managed metadata for a field
// (directus_fields): interface, display, validation, and layout options.
type FieldMeta struct {
	ID                     int                `json:"id,omitempty"`
	Collection             string             `json:"collection"`
	Field                  string             `json:"field"`
	Group                  string             `json:"group,omitempty"`
	Hidden                 bool               `json:"hidden"`
	Interface              string             `json:"interface,omitempty"`
	Display                string             `json:"display,omitempty"`
	Options                map[string]any     `json:"options,omitempty"`
	DisplayOptions         map[string]any     `json:"display_options,omitempty"`
	Readonly               bool               `json:"readonly"`
	Required               bool               `json:"required"`
	Sort                   int                `json:"sort,omitempty"`
	Special                []string           `json:"special"`
	Translations           []FieldTranslation `json:"translations"`
	Width                  string             `json:"width,omitempty"`
	Note                   string             `json:"note,omitempty"`
	Conditions             []FieldCondition   `json:"conditions"`
	Validation             map[string]any     `json:"validation"`
	ValidationMessage      string             `json:"validation_message,omitempty"`
	Searchable             bool               `json:"searchable"`
	System                 bool               `json:"system,omitempty"`
	ClearHiddenValueOnSave bool               `json:"clear_hidden_value_on_save,omitempty"`
}

// FieldSchema describes the underlying database column backing a field.
type FieldSchema struct {
	Name                 string `json:"name"`
	Table                string `json:"table"`
	Schema               string `json:"schema,omitempty"`
	DataType             string `json:"data_type"`
	DefaultValue         any    `json:"default_value"`
	GenerationExpression string `json:"generation_expression,omitempty"`
	MaxLength            int    `json:"max_length,omitempty"`
	NumericPrecision     int    `json:"numeric_precision,omitempty"`
	NumericScale         int    `json:"numeric_scale,omitempty"`
	IsGenerated          bool   `json:"is_generated"`
	IsNullable           bool   `json:"is_nullable"`
	IsUnique             bool   `json:"is_unique"`
	IsIndexed            bool   `json:"is_indexed"`
	IsPrimaryKey         bool   `json:"is_primary_key"`
	HasAutoIncrement     bool   `json:"has_auto_increment"`
	ForeignKeySchema     string `json:"foreign_key_schema,omitempty"`
	ForeignKeyTable      string `json:"foreign_key_table,omitempty"`
	ForeignKeyColumn     string `json:"foreign_key_column,omitempty"`
	Comment              string `json:"comment,omitempty"`
}

// Field is the read/write model for a Directus field, pairing its Directus
// metadata (Meta) with the underlying column definition (Schema).
type Field struct {
	Collection string       `json:"collection"`
	Field      string       `json:"field"`
	Type       string       `json:"type,omitempty"`
	Meta       *FieldMeta   `json:"meta,omitempty"`
	Schema     *FieldSchema `json:"schema,omitempty"`
}

// GetFields lists all fields across every collection.
func (c *Client) GetFields() ([]Field, error) {
	return c.GetFieldsByCollection("")
}

// GetFieldsByCollection lists the fields of a single collection. An empty
// collection lists fields across every collection.
func (c *Client) GetFieldsByCollection(collection string) ([]Field, error) {
	return request[[]Field](c, "GET", fmt.Sprintf("/fields/%s", collection), nil, nil)
}

// GetFieldByCollectionAndName retrieves a single field by collection and name.
func (c *Client) GetFieldByCollectionAndName(collection string, name string) (*Field, error) {
	if collection == "" || name == "" {
		return nil, fmt.Errorf("collection and name must be provided")
	}
	return request[*Field](c, "GET", fmt.Sprintf("/fields/%s/%s", collection, name), nil, nil)
}

// withFieldLocation returns a copy of field with its collection (and, when name
// is non-empty, its field name) populated so the JSON payload is consistent
// with the request URL. Directus rejects an empty "collection" on field writes.
func withFieldLocation(field *Field, collection, name string) *Field {
	f := Field{Collection: collection}
	if field != nil {
		f = *field
		f.Collection = collection
	}
	if name != "" {
		f.Field = name
	}
	if f.Meta != nil {
		m := *f.Meta
		m.Collection = collection
		if name != "" {
			m.Field = name
		}
		f.Meta = &m
	}
	return &f
}

// CreateField creates a new field in the given collection.
func (c *Client) CreateField(collection string, field *Field, q *Query) (*Field, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	return request[*Field](c, "POST", fmt.Sprintf("/fields/%s", collection), q, withFieldLocation(field, collection, ""))
}

// PatchField updates a single field in the given collection.
func (c *Client) PatchField(collection string, name string, field *Field, q *Query) (*Field, error) {
	if collection == "" || name == "" {
		return nil, fmt.Errorf("collection and name must be provided")
	}
	return request[*Field](c, "PATCH", fmt.Sprintf("/fields/%s/%s", collection, name), q, withFieldLocation(field, collection, name))
}

// PatchFields updates multiple fields in the given collection.
func (c *Client) PatchFields(collection string, fields []Field, q *Query) ([]Field, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}
	payload := make([]Field, len(fields))
	for i := range fields {
		payload[i] = *withFieldLocation(&fields[i], collection, "")
	}
	return request[[]Field](c, "PATCH", fmt.Sprintf("/fields/%s", collection), q, payload)
}

// DeleteField deletes a field from the given collection.
func (c *Client) DeleteField(collection string, name string) error {
	if collection == "" || name == "" {
		return fmt.Errorf("collection and name must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/fields/%s/%s", collection, name), nil, nil)
}
