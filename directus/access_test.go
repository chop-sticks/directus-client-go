package directus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetAccesses(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"a1","policy":"p1"},{"id":"a2","policy":"p2"}]}`)
	})
	accesses, err := client.GetAccesses(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access" || gotMethod != "GET" {
		t.Errorf("expected GET /access, got %s %s", gotMethod, gotPath)
	}
	if len(accesses) != 2 || accesses[0].ID != "a1" {
		t.Errorf("unexpected accesses: %+v", accesses)
	}
}

func TestGetAccessesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetAccesses(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetAccessesBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetAccesses(nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetAccess(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"a1","policy":"p1","user":"u1"}}`)
	})
	access, err := client.GetAccess("a1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access/a1" || gotMethod != "GET" {
		t.Errorf("expected GET /access/a1, got %s %s", gotMethod, gotPath)
	}
	if access == nil || access.ID != "a1" {
		t.Errorf("unexpected access: %+v", access)
	}
}

func TestGetAccessValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetAccess("", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetAccessBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetAccess("a1", nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestCreateAccess(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"data":{"id":"a1","policy":"p1","user":"u1"}}`)
	})
	access, err := client.CreateAccess(&Access{Policy: "p1", User: "u1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access" || gotMethod != "POST" {
		t.Errorf("expected POST /access, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"policy":"p1"`) {
		t.Errorf("expected policy in body, got %s", gotBody)
	}
	if access == nil || access.ID != "a1" {
		t.Errorf("unexpected access: %+v", access)
	}
}

func TestCreateAccessBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateAccess(&Access{}, nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestCreateAccesses(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"a1"},{"id":"a2"}]}`)
	})
	accesses, err := client.CreateAccesses([]Access{{Policy: "p1"}, {Policy: "p2"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access" || gotMethod != "POST" {
		t.Errorf("expected POST /access, got %s %s", gotMethod, gotPath)
	}
	if len(accesses) != 2 {
		t.Errorf("unexpected accesses: %+v", accesses)
	}
}

func TestPatchAccess(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"a1","policy":"p2"}}`)
	})
	access, err := client.PatchAccess("a1", &Access{Policy: "p2"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access/a1" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /access/a1, got %s %s", gotMethod, gotPath)
	}
	if access == nil || access.ID != "a1" {
		t.Errorf("unexpected access: %+v", access)
	}
}

func TestPatchAccessValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchAccess("", &Access{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchAccesses(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"data":[{"id":"a1"},{"id":"a2"}]}`)
	})
	accesses, err := client.PatchAccesses([]string{"a1", "a2"}, &Access{Policy: "p1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /access, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"keys":["a1","a2"]`) {
		t.Errorf("expected keys in body, got %s", gotBody)
	}
	if len(accesses) != 2 {
		t.Errorf("unexpected accesses: %+v", accesses)
	}
}

func TestPatchAccessesBatch(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"a1"},{"id":"a2"}]}`)
	})
	accesses, err := client.PatchAccessesBatch([]Access{{ID: "a1"}, {ID: "a2"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /access, got %s %s", gotMethod, gotPath)
	}
	if len(accesses) != 2 {
		t.Errorf("unexpected accesses: %+v", accesses)
	}
}

func TestDeleteAccess(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteAccess("a1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access/a1" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /access/a1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteAccessValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteAccess(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteAccesses(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteAccesses([]string{"a1", "a2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/access" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /access, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `["a1","a2"]`) {
		t.Errorf("expected raw keys body, got %s", gotBody)
	}
}

func TestAccessRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetAccesses(nil); err == nil {
		t.Error("GetAccesses: expected request build error, got nil")
	}
}
