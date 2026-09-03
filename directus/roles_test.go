package directus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetRoles(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"r1","name":"Admin"},{"id":"r2","name":"Editor"}]}`)
	})
	roles, err := client.GetRoles(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles" || gotMethod != "GET" {
		t.Errorf("expected GET /roles, got %s %s", gotMethod, gotPath)
	}
	if len(roles) != 2 || roles[0].Name != "Admin" {
		t.Errorf("unexpected roles: %+v", roles)
	}
}

func TestGetRolesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetRoles(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetRolesBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetRoles(nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetRole(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"r1","name":"Admin"}}`)
	})
	role, err := client.GetRole("r1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles/r1" || gotMethod != "GET" {
		t.Errorf("expected GET /roles/r1, got %s %s", gotMethod, gotPath)
	}
	if role == nil || role.Name != "Admin" {
		t.Errorf("unexpected role: %+v", role)
	}
}

func TestGetRoleValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetRole("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetRoleBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetRole("r1", nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetRolesMe(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		fmt.Fprint(w, `{"data":[{"id":"r1","name":"Admin"}]}`)
	})
	roles, err := client.GetRolesMe(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles/me" {
		t.Errorf("expected /roles/me, got %s", gotPath)
	}
	if len(roles) != 1 {
		t.Errorf("unexpected roles: %+v", roles)
	}
}

func TestCreateRole(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"r1","name":"Admin"}}`)
	})
	role, err := client.CreateRole(&Role{Name: "Admin"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles" || gotMethod != "POST" {
		t.Errorf("expected POST /roles, got %s %s", gotMethod, gotPath)
	}
	if role == nil || role.ID != "r1" {
		t.Errorf("unexpected role: %+v", role)
	}
}

func TestCreateRoleBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateRole(&Role{}, nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestCreateRoles(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"r1"},{"id":"r2"}]}`)
	})
	roles, err := client.CreateRoles([]Role{{Name: "A"}, {Name: "B"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles" || gotMethod != "POST" {
		t.Errorf("expected POST /roles, got %s %s", gotMethod, gotPath)
	}
	if len(roles) != 2 {
		t.Errorf("unexpected roles: %+v", roles)
	}
}

func TestPatchRole(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"r1","name":"Renamed"}}`)
	})
	role, err := client.PatchRole("r1", &Role{Name: "Renamed"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles/r1" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /roles/r1, got %s %s", gotMethod, gotPath)
	}
	if role == nil || role.Name != "Renamed" {
		t.Errorf("unexpected role: %+v", role)
	}
}

func TestPatchRoleValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchRole("", &Role{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchRoles(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"data":[{"id":"r1"},{"id":"r2"}]}`)
	})
	roles, err := client.PatchRoles([]string{"r1", "r2"}, &Role{Icon: "star"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /roles, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"keys"`) || !strings.Contains(gotBody, `"data"`) {
		t.Errorf("expected keys+data body, got %s", gotBody)
	}
	if len(roles) != 2 {
		t.Errorf("unexpected roles: %+v", roles)
	}
}

func TestPatchRolesBatch(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"r1"}]}`)
	})
	roles, err := client.PatchRolesBatch([]Role{{ID: "r1", Name: "A"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /roles, got %s %s", gotMethod, gotPath)
	}
	if len(roles) != 1 {
		t.Errorf("unexpected roles: %+v", roles)
	}
}

func TestDeleteRole(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteRole("r1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles/r1" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /roles/r1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteRoleValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteRole(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteRoles(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteRoles([]string{"r1", "r2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/roles" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /roles, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `["r1","r2"]`) {
		t.Errorf("expected raw keys body, got %s", gotBody)
	}
}

func TestRolesRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetRoles(nil); err == nil {
		t.Error("GetRoles: expected request build error, got nil")
	}
}
