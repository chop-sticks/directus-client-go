package directus

import (
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
