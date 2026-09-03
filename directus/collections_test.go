package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// badHostClient returns a client whose HostURL cannot be parsed into a request
// URL, forcing http.NewRequest to fail.
func badHostClient(t *testing.T) *Client {
	t.Helper()
	token := "test_token"
	host := "://bad host"
	client, _ := NewClient(&host, &token)
	return client
}

func TestCollectionsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetCollections(); err == nil {
		t.Error("GetCollections: expected request build error, got nil")
	}
	if _, err := client.GetCollectionByName("sites"); err == nil {
		t.Error("GetCollectionByName: expected request build error, got nil")
	}
	if _, err := client.CreateCollection(&CollectionRequest{Collection: "sites"}, nil); err == nil {
		t.Error("CreateCollection: expected request build error, got nil")
	}
}

func TestGetCollections(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"collection":"sites","meta":{"icon":"apartment"}},{"collection":"posts"}]}`)
	})

	cols, err := client.GetCollections()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/collections" {
		t.Errorf("expected path /collections, got %s", gotPath)
	}
	if len(cols) != 2 {
		t.Fatalf("expected 2 collections, got %d", len(cols))
	}
	if cols[0].Collection != "sites" {
		t.Errorf("expected first collection 'sites', got %q", cols[0].Collection)
	}
	if cols[0].Meta == nil || cols[0].Meta.Icon != "apartment" {
		t.Errorf("expected meta icon 'apartment', got %+v", cols[0].Meta)
	}
	if cols[1].Meta != nil {
		t.Errorf("expected nil meta for second collection, got %+v", cols[1].Meta)
	}
}

func TestGetCollectionsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.GetCollections(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetCollectionsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetCollections(); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetCollectionByName(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"sites","meta":{"collapse":"open","singleton":true}}}`)
	})

	col, err := client.GetCollectionByName("sites")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/collections/sites" {
		t.Errorf("expected path /collections/sites, got %s", gotPath)
	}
	if col.Collection != "sites" {
		t.Errorf("expected 'sites', got %q", col.Collection)
	}
	if col.Meta == nil || !col.Meta.Singleton {
		t.Errorf("expected singleton meta, got %+v", col.Meta)
	}
}

func TestGetCollectionByNameHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetCollectionByName("missing"); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestCreateCollection(t *testing.T) {
	var gotMethod, gotPath string
	var gotReq CollectionRequest
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotReq)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"sites"}}`)
	})

	req := &CollectionRequest{
		Collection: "sites",
		Meta:       &CollectionMeta{Icon: "apartment", Collapse: "open"},
		Fields:     []Field{{Field: "id", Type: "uuid"}},
	}
	col, err := client.CreateCollection(req, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST for create, got %s", gotMethod)
	}
	if gotPath != "/collections" {
		t.Errorf("expected path /collections, got %s", gotPath)
	}
	if gotReq.Collection != "sites" {
		t.Errorf("expected request collection 'sites', got %q", gotReq.Collection)
	}
	if len(gotReq.Fields) != 1 || gotReq.Fields[0].Field != "id" {
		t.Errorf("expected request to carry 1 field 'id', got %+v", gotReq.Fields)
	}
	if col.Collection != "sites" {
		t.Errorf("expected response 'sites', got %q", col.Collection)
	}
}

func TestPatchCollection(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"collection":"sites"}}`)
	})

	if _, err := client.PatchCollection("sites", &CollectionRequest{Collection: "sites"}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("expected PATCH for update, got %s", gotMethod)
	}
	if gotPath != "/collections/sites" {
		t.Errorf("expected path /collections/sites, got %s", gotPath)
	}
}

func TestProcessCollectionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.CreateCollection(&CollectionRequest{Collection: "x"}, nil); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestProcessCollectionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{`)
	})
	if _, err := client.CreateCollection(&CollectionRequest{Collection: "x"}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetCollectionByNameBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":`)
	})
	if _, err := client.GetCollectionByName("sites"); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

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
