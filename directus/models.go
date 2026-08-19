package directus

type CollectionTranslation struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
	Singular    string `json:"singular"`
	Plural      string `json:"plural"`
}

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

type FieldTranslation struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

type FieldCondition struct {
	Name                   string         `json:"name"`
	Rule                   map[string]any `json:"rule"`
	Readonly               bool           `json:"readonly,omitempty"`
	Hidden                 bool           `json:"hidden,omitempty"`
	Options                map[string]any `json:"options,omitempty"`
	Required               bool           `json:"required,omitempty"`
	ClearHiddenValueOnSave bool           `json:"clear_hidden_value_on_save,omitempty"`
}

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

type Field struct {
	Collection string       `json:"collection"`
	Field      string       `json:"field"`
	Type       string       `json:"type"`
	Meta       *FieldMeta   `json:"meta"`
	Schema     *FieldSchema `json:"schema"`
}
