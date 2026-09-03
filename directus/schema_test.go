package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestSchemaRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetSchemaSnapshot(nil, nil); err == nil {
		t.Error("GetSchemaSnapshot: expected request build error, got nil")
	}
	if _, err := client.SchemaDiffSnapshot(&SchemaSnapshot{}, false, ""); err == nil {
		t.Error("SchemaDiffSnapshot: expected request build error, got nil")
	}
	if err := client.SchemaApply(&SchemaDiff{}, false); err == nil {
		t.Error("SchemaApply: expected request build error, got nil")
	}
}

func TestGetSchemaSnapshot(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"version":1,"directus":"11.0.0","vendor":"postgres","collections":[{"collection":"articles"}],"fields":[],"systemFields":[],"relations":[]}}`)
	})

	snap, err := client.GetSchemaSnapshot(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/schema/snapshot" {
		t.Errorf("expected path /schema/snapshot, got %s", gotPath)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if snap.Version != 1 {
		t.Errorf("expected version 1, got %d", snap.Version)
	}
	if snap.Directus != "11.0.0" {
		t.Errorf("expected directus '11.0.0', got %q", snap.Directus)
	}
	if len(snap.Collections) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(snap.Collections))
	}
	if snap.Collections[0]["collection"] != "articles" {
		t.Errorf("expected collection 'articles', got %v", snap.Collections[0]["collection"])
	}
}

func TestGetSchemaSnapshotIncludeCollections(t *testing.T) {
	var gotQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("includeCollections")
		if r.URL.Query().Get("excludeCollections") != "" {
			t.Error("expected no excludeCollections when include is set")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"version":1}}`)
	})

	if _, err := client.GetSchemaSnapshot([]string{"articles", "authors"}, []string{"secret"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery != "articles,authors" {
		t.Errorf("expected includeCollections 'articles,authors', got %q", gotQuery)
	}
}

func TestGetSchemaSnapshotExcludeCollections(t *testing.T) {
	var gotQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("excludeCollections")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"version":1}}`)
	})

	if _, err := client.GetSchemaSnapshot(nil, []string{"secret", "internal"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery != "secret,internal" {
		t.Errorf("expected excludeCollections 'secret,internal', got %q", gotQuery)
	}
}

func TestGetSchemaSnapshotHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.GetSchemaSnapshot(nil, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetSchemaSnapshotBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{not json`)
	})
	if _, err := client.GetSchemaSnapshot(nil, nil); err == nil {
		t.Error("expected JSON decode error, got nil")
	}
}

func TestSchemaDiffSnapshot(t *testing.T) {
	var gotPath, gotMethod, gotForce, gotMode string
	var gotBody SchemaSnapshot
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotForce = r.URL.Query().Get("force")
		gotMode = r.URL.Query().Get("mode")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"hash":"abc123","diff":{"collections":[]}}}`)
	})

	diff, err := client.SchemaDiffSnapshot(&SchemaSnapshot{Version: 1, Directus: "11.0.0"}, true, "diff")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/schema/diff" {
		t.Errorf("expected path /schema/diff, got %s", gotPath)
	}
	if gotForce != "true" {
		t.Errorf("expected force 'true', got %q", gotForce)
	}
	if gotMode != "diff" {
		t.Errorf("expected mode 'diff', got %q", gotMode)
	}
	if gotBody.Directus != "11.0.0" {
		t.Errorf("expected body directus '11.0.0', got %q", gotBody.Directus)
	}
	if diff == nil {
		t.Fatal("expected diff, got nil")
	}
	if diff.Hash != "abc123" {
		t.Errorf("expected hash 'abc123', got %q", diff.Hash)
	}
}

func TestSchemaDiffSnapshotNoDiff(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		// No difference: empty body / 204.
		w.WriteHeader(http.StatusNoContent)
	})

	diff, err := client.SchemaDiffSnapshot(&SchemaSnapshot{Version: 1}, false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff != nil {
		t.Errorf("expected nil diff for no-difference response, got %+v", diff)
	}
}

func TestSchemaDiffSnapshotValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.SchemaDiffSnapshot(nil, false, ""); err == nil {
		t.Error("expected validation error for nil snapshot, got nil")
	}
	if called {
		t.Error("expected no HTTP call when snapshot is nil")
	}
}

func TestSchemaDiffSnapshotHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.SchemaDiffSnapshot(&SchemaSnapshot{}, false, ""); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestSchemaDiffSnapshotBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{not json`)
	})
	if _, err := client.SchemaDiffSnapshot(&SchemaSnapshot{}, false, ""); err == nil {
		t.Error("expected JSON decode error, got nil")
	}
}

func TestSchemaApply(t *testing.T) {
	var gotPath, gotMethod, gotForce string
	var gotBody SchemaDiff
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotForce = r.URL.Query().Get("force")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.SchemaApply(&SchemaDiff{Hash: "abc123", Diff: map[string]any{"collections": []any{}}}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/schema/apply" {
		t.Errorf("expected path /schema/apply, got %s", gotPath)
	}
	if gotForce != "true" {
		t.Errorf("expected force 'true', got %q", gotForce)
	}
	if gotBody.Hash != "abc123" {
		t.Errorf("expected body hash 'abc123', got %q", gotBody.Hash)
	}
}

func TestSchemaApplyValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.SchemaApply(nil, false); err == nil {
		t.Error("expected validation error for nil diff, got nil")
	}
	if called {
		t.Error("expected no HTTP call when diff is nil")
	}
}

func TestSchemaApplyHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if err := client.SchemaApply(&SchemaDiff{}, false); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}
