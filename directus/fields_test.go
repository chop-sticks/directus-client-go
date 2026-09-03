package directus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestFieldsRequestBuildError(t *testing.T) {
	token := "test_token"
	host := "://bad host"
	client, _ := NewClient(&host, &token)

	if _, err := client.GetFields(); err == nil {
		t.Error("GetFields: expected request build error, got nil")
	}
	if _, err := client.GetFieldsByCollection("sites"); err == nil {
		t.Error("GetFieldsByCollection: expected request build error, got nil")
	}
	if _, err := client.GetFieldByCollectionAndName("sites", "id"); err == nil {
		t.Error("GetFieldByCollectionAndName: expected request build error, got nil")
	}
}

func TestGetFields(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"collection":"sites","field":"id","type":"uuid"},{"collection":"sites","field":"name","type":"string"}]}`)
	})

	fields, err := client.GetFields()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/fields/" {
		t.Errorf("expected path /fields/, got %s", gotPath)
	}
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
	if fields[0].Field != "id" || fields[0].Type != "uuid" {
		t.Errorf("unexpected first field: %+v", fields[0])
	}
}

func TestGetFieldsByCollection(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[]}`)
	})

	if _, err := client.GetFieldsByCollection("sites"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/fields/sites" {
		t.Errorf("expected path /fields/sites, got %s", gotPath)
	}
}

func TestGetFieldsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetFields(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetFieldsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.GetFields(); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetFieldByCollectionAndName(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"sites","field":"id","type":"string","meta":{"id":1,"collection":"sites","field":"id","interface":"input","sort":1,"width":"full","searchable":true},"schema":{"name":"id","table":"sites","data_type":"character varying","max_length":255,"is_nullable":false,"is_unique":true,"is_primary_key":true}}}`)
	})

	field, err := client.GetFieldByCollectionAndName("sites", "id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/fields/sites/id" {
		t.Errorf("expected path /fields/sites/id, got %s", gotPath)
	}
	if field.Collection != "sites" || field.Field != "id" || field.Type != "string" {
		t.Errorf("unexpected field top-level: %+v", field)
	}
	if field.Meta == nil {
		t.Fatal("expected non-nil meta")
	}
	if field.Meta.Interface != "input" || field.Meta.Sort != 1 || !field.Meta.Searchable {
		t.Errorf("unexpected meta: %+v", field.Meta)
	}
	if field.Schema == nil {
		t.Fatal("expected non-nil schema")
	}
	if field.Schema.MaxLength != 255 || !field.Schema.IsPrimaryKey || field.Schema.IsNullable {
		t.Errorf("unexpected schema: %+v", field.Schema)
	}
}

func TestGetFieldByCollectionAndNameValidation(t *testing.T) {
	// Must fail before any HTTP call is made.
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	cases := []struct{ collection, name string }{
		{"", "id"},
		{"sites", ""},
		{"", ""},
	}
	for _, tc := range cases {
		_, err := client.GetFieldByCollectionAndName(tc.collection, tc.name)
		if err == nil {
			t.Errorf("collection=%q name=%q: expected validation error, got nil", tc.collection, tc.name)
		} else if !strings.Contains(err.Error(), "collection and name must be provided") {
			t.Errorf("unexpected error message: %v", err)
		}
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetFieldByCollectionAndNameHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetFieldByCollectionAndName("sites", "missing"); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetFieldByCollectionAndNameBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetFieldByCollectionAndName("sites", "id"); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
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

func TestCreateFieldPopulatesCollection(t *testing.T) {
	var gotPath, gotMethod string
	var got Field
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"sites","field":"title","type":"string"}}`)
	})

	f, err := client.CreateField("sites", &Field{Field: "title", Type: "string", Meta: &FieldMeta{Interface: "input"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/fields/sites" {
		t.Errorf("expected path /fields/sites, got %s", gotPath)
	}
	if got.Collection != "sites" {
		t.Errorf("expected body collection 'sites', got %q", got.Collection)
	}
	if got.Meta == nil || got.Meta.Collection != "sites" {
		t.Errorf("expected body meta.collection 'sites', got %+v", got.Meta)
	}
	if f.Field != "title" {
		t.Errorf("expected returned field 'title', got %q", f.Field)
	}
}

func TestPatchFieldPopulatesLocation(t *testing.T) {
	var gotPath string
	var got Field
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"sites","field":"title","type":"string"}}`)
	})

	if _, err := client.PatchField("sites", "title", &Field{Type: "string"}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/fields/sites/title" {
		t.Errorf("expected path /fields/sites/title, got %s", gotPath)
	}
	if got.Collection != "sites" || got.Field != "title" {
		t.Errorf("expected body collection/field sites/title, got %q/%q", got.Collection, got.Field)
	}
}

func TestFieldWriteValidation(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.CreateField("", &Field{}, nil); err == nil {
		t.Error("CreateField: expected error for empty collection")
	}
	if _, err := client.PatchField("sites", "", &Field{}, nil); err == nil {
		t.Error("PatchField: expected error for empty field name")
	}
	if err := client.DeleteField("", "x"); err == nil {
		t.Error("DeleteField: expected error for empty collection")
	}
}
