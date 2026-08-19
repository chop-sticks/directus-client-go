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
	if _, err := client.CreateCollection(&CollectionRequest{Collection: "sites"}); err == nil {
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
	col, err := client.CreateCollection(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST for create, got %s", gotMethod)
	}
	if gotPath != "/collections/" {
		t.Errorf("expected path /collections/, got %s", gotPath)
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

	if _, err := client.PatchCollection("sites", &CollectionRequest{Collection: "sites"}); err != nil {
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
	if _, err := client.CreateCollection(&CollectionRequest{Collection: "x"}); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestProcessCollectionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{`)
	})
	if _, err := client.CreateCollection(&CollectionRequest{Collection: "x"}); err == nil {
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
