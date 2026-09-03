package directus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetPolicies(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"p1","name":"Base"},{"id":"p2","name":"Extra"}]}`)
	})
	policies, err := client.GetPolicies(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies" || gotMethod != "GET" {
		t.Errorf("expected GET /policies, got %s %s", gotMethod, gotPath)
	}
	if len(policies) != 2 || policies[0].Name != "Base" {
		t.Errorf("unexpected policies: %+v", policies)
	}
}

func TestGetPoliciesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetPolicies(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetPoliciesBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetPolicies(nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetPolicy(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"p1","name":"Base","admin_access":true}}`)
	})
	policy, err := client.GetPolicy("p1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies/p1" || gotMethod != "GET" {
		t.Errorf("expected GET /policies/p1, got %s %s", gotMethod, gotPath)
	}
	if policy == nil || !policy.AdminAccess {
		t.Errorf("unexpected policy: %+v", policy)
	}
}

func TestGetPolicyValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetPolicy("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetPolicyBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetPolicy("p1", nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetPolicyGlobals(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"app_access":true,"admin_access":false,"enforce_tfa":true}}`)
	})
	globals, err := client.GetPolicyGlobals()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies/me/globals" || gotMethod != "GET" {
		t.Errorf("expected GET /policies/me/globals, got %s %s", gotMethod, gotPath)
	}
	if globals == nil || !globals.AppAccess || globals.AdminAccess || !globals.EnforceTFA {
		t.Errorf("unexpected globals: %+v", globals)
	}
}

func TestGetPolicyGlobalsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetPolicyGlobals(); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestCreatePolicy(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"p1","name":"Base"}}`)
	})
	policy, err := client.CreatePolicy(&Policy{Name: "Base"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies" || gotMethod != "POST" {
		t.Errorf("expected POST /policies, got %s %s", gotMethod, gotPath)
	}
	if policy == nil || policy.ID != "p1" {
		t.Errorf("unexpected policy: %+v", policy)
	}
}

func TestCreatePolicyBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreatePolicy(&Policy{}, nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestCreatePolicies(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"p1"},{"id":"p2"}]}`)
	})
	policies, err := client.CreatePolicies([]Policy{{Name: "A"}, {Name: "B"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies" || gotMethod != "POST" {
		t.Errorf("expected POST /policies, got %s %s", gotMethod, gotPath)
	}
	if len(policies) != 2 {
		t.Errorf("unexpected policies: %+v", policies)
	}
}

func TestPatchPolicy(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"p1","name":"Renamed"}}`)
	})
	policy, err := client.PatchPolicy("p1", &Policy{Name: "Renamed"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies/p1" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /policies/p1, got %s %s", gotMethod, gotPath)
	}
	if policy == nil || policy.Name != "Renamed" {
		t.Errorf("unexpected policy: %+v", policy)
	}
}

func TestPatchPolicyValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchPolicy("", &Policy{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchPolicies(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"data":[{"id":"p1"},{"id":"p2"}]}`)
	})
	policies, err := client.PatchPolicies([]string{"p1", "p2"}, &Policy{EnforceTFA: true}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /policies, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"keys"`) || !strings.Contains(gotBody, `"data"`) {
		t.Errorf("expected keys+data body, got %s", gotBody)
	}
	if len(policies) != 2 {
		t.Errorf("unexpected policies: %+v", policies)
	}
}

func TestPatchPoliciesBatch(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"p1"}]}`)
	})
	policies, err := client.PatchPoliciesBatch([]Policy{{ID: "p1", Name: "A"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /policies, got %s %s", gotMethod, gotPath)
	}
	if len(policies) != 1 {
		t.Errorf("unexpected policies: %+v", policies)
	}
}

func TestDeletePolicy(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeletePolicy("p1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies/p1" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /policies/p1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeletePolicyValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeletePolicy(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeletePolicies(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeletePolicies([]string{"p1", "p2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/policies" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /policies, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `["p1","p2"]`) {
		t.Errorf("expected raw keys body, got %s", gotBody)
	}
}

func TestPoliciesRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetPolicies(nil); err == nil {
		t.Error("GetPolicies: expected request build error, got nil")
	}
}
