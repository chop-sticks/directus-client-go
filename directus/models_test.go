package directus

import (
	"encoding/json"
	"testing"
)

// marshalToMap marshals v and unmarshals into a generic map so tests can assert
// on JSON key presence and values regardless of struct field order.
func marshalToMap(t *testing.T, v interface{}) map[string]json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}
	return m
}

func TestCollectionMarshalOmitsNilMetaSchema(t *testing.T) {
	m := marshalToMap(t, Collection{Collection: "sites"})
	if _, ok := m["meta"]; ok {
		t.Error("expected meta to be omitted when nil")
	}
	if _, ok := m["schema"]; ok {
		t.Error("expected schema to be omitted when nil")
	}
	if _, ok := m["collection"]; !ok {
		t.Error("expected collection key to be present")
	}
}

func TestCollectionRequestMarshalIncludesFields(t *testing.T) {
	req := CollectionRequest{
		Collection: "sites",
		Fields:     []Field{{Field: "id", Type: "uuid"}},
	}
	m := marshalToMap(t, req)
	if _, ok := m["fields"]; !ok {
		t.Error("expected fields key to be present when set")
	}

	// Absent when empty.
	m2 := marshalToMap(t, CollectionRequest{Collection: "sites"})
	if _, ok := m2["fields"]; ok {
		t.Error("expected fields key to be omitted when empty")
	}
}

func TestCollectionMetaUnmarshal(t *testing.T) {
	data := `{
		"collection":"sites",
		"note":"desc",
		"hidden":true,
		"singleton":false,
		"icon":"apartment",
		"translations":null,
		"item_duplication_fields":["a","b"],
		"system":true,
		"versioning":true,
		"autosave_revision_interval":5,
		"accountability":"all",
		"collapse":"open",
		"status":"active"
	}`
	var meta CollectionMeta
	if err := json.Unmarshal([]byte(data), &meta); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if meta.Collection != "sites" || !meta.Hidden || meta.Singleton {
		t.Errorf("unexpected scalar fields: %+v", meta)
	}
	if len(meta.ItemDuplicationFields) != 2 || meta.ItemDuplicationFields[0] != "a" {
		t.Errorf("expected item_duplication_fields [a b], got %v", meta.ItemDuplicationFields)
	}
	if !meta.System || !meta.Versioning || meta.AutosaveRevisionInterval != 5 {
		t.Errorf("unexpected new fields: system=%v versioning=%v interval=%d", meta.System, meta.Versioning, meta.AutosaveRevisionInterval)
	}
	if meta.Accountability != "all" || meta.Collapse != "open" || meta.Status != "active" {
		t.Errorf("unexpected enum fields: %+v", meta)
	}
	if meta.Translations != nil {
		t.Errorf("expected nil translations, got %v", meta.Translations)
	}
}

func TestFieldSchemaOmitempty(t *testing.T) {
	m := marshalToMap(t, FieldSchema{Name: "id", Table: "sites", DataType: "varchar", IsNullable: true})

	// Non-omitempty keys must always be present.
	for _, key := range []string{"name", "table", "data_type", "default_value", "is_generated", "is_nullable", "is_unique", "is_indexed", "is_primary_key", "has_auto_increment"} {
		if _, ok := m[key]; !ok {
			t.Errorf("expected key %q to be present", key)
		}
	}
	// Nullable value fields must be omitted when zero.
	for _, key := range []string{"schema", "generation_expression", "max_length", "numeric_precision", "numeric_scale", "foreign_key_schema", "foreign_key_table", "foreign_key_column", "comment"} {
		if _, ok := m[key]; ok {
			t.Errorf("expected key %q to be omitted when zero", key)
		}
	}
	// default_value with no omitempty serializes as JSON null.
	if string(m["default_value"]) != "null" {
		t.Errorf("expected default_value to be null, got %s", m["default_value"])
	}
}

func TestFieldUnmarshalNested(t *testing.T) {
	data := `{
		"collection":"sites",
		"field":"id",
		"type":"string",
		"meta":{"id":1,"interface":"input","sort":1,"searchable":true,"special":["uuid"]},
		"schema":{"name":"id","table":"sites","data_type":"character varying","max_length":255,"is_primary_key":true}
	}`
	var field Field
	if err := json.Unmarshal([]byte(data), &field); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if field.Meta == nil || field.Meta.Interface != "input" || field.Meta.Sort != 1 || !field.Meta.Searchable {
		t.Errorf("unexpected meta: %+v", field.Meta)
	}
	if len(field.Meta.Special) != 1 || field.Meta.Special[0] != "uuid" {
		t.Errorf("expected special [uuid], got %v", field.Meta.Special)
	}
	if field.Schema == nil || field.Schema.MaxLength != 255 || !field.Schema.IsPrimaryKey {
		t.Errorf("unexpected schema: %+v", field.Schema)
	}
}

func TestFieldMetaRequiredNotOmitted(t *testing.T) {
	// required/hidden/readonly/searchable have no omitempty: false must serialize
	// so a write can explicitly clear them.
	m := marshalToMap(t, FieldMeta{Collection: "sites", Field: "id"})
	for _, key := range []string{"required", "hidden", "readonly", "searchable"} {
		v, ok := m[key]
		if !ok {
			t.Errorf("expected key %q to be present even when false", key)
			continue
		}
		if string(v) != "false" {
			t.Errorf("expected %q to be false, got %s", key, v)
		}
	}
}
