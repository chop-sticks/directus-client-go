package directus

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestRevisionsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetRevisions(nil); err == nil {
		t.Error("GetRevisions: expected request build error, got nil")
	}
	if _, err := client.GetRevision(1, nil); err == nil {
		t.Error("GetRevision: expected request build error, got nil")
	}
}

func TestGetRevisions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1,"collection":"sites"},{"id":2,"collection":"pages"}]}`)
	})

	items, err := client.GetRevisions(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/revisions" {
		t.Errorf("expected GET /revisions, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 2 || items[0].ID != 1 || items[0].Collection != "sites" {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestGetRevisionsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetRevisions(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetRevisionsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetRevisions(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetRevision(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":5,"collection":"sites","data":{"name":"x"}}}`)
	})

	item, err := client.GetRevision(5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/revisions/5" {
		t.Errorf("expected GET /revisions/5, got %s %s", gotMethod, gotPath)
	}
	if item == nil || item.ID != 5 || item.Data["name"] != "x" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestGetRevisionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	for _, id := range []int{0, -1} {
		if _, err := client.GetRevision(id, nil); err == nil {
			t.Errorf("id=%d: expected validation error, got nil", id)
		} else if !strings.Contains(err.Error(), "id must be provided") {
			t.Errorf("id=%d: unexpected error message: %v", id, err)
		}
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetRevisionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetRevision(999, nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetRevisionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetRevision(1, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}
