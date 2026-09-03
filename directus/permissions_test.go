package directus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetPermissions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":1,"collection":"sites","action":"read"},{"id":2,"collection":"pages","action":"create"}]}`)
	})
	perms, err := client.GetPermissions(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions" || gotMethod != "GET" {
		t.Errorf("expected GET /permissions, got %s %s", gotMethod, gotPath)
	}
	if len(perms) != 2 || perms[0].ID != 1 || perms[0].Action != "read" {
		t.Errorf("unexpected perms: %+v", perms)
	}
}

func TestGetPermissionsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetPermissions(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetPermissionsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetPermissions(nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetPermission(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":42,"collection":"sites","action":"read"}}`)
	})
	perm, err := client.GetPermission(42, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions/42" || gotMethod != "GET" {
		t.Errorf("expected GET /permissions/42, got %s %s", gotMethod, gotPath)
	}
	if perm == nil || perm.ID != 42 {
		t.Errorf("unexpected perm: %+v", perm)
	}
}

func TestGetPermissionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetPermission(1, nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetUserPermissions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"sites":{"read":{"access":"full"}}}}`)
	})
	perms, err := client.GetUserPermissions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions/me" || gotMethod != "GET" {
		t.Errorf("expected GET /permissions/me, got %s %s", gotMethod, gotPath)
	}
	if _, ok := perms["sites"]; !ok {
		t.Errorf("unexpected perms: %+v", perms)
	}
}

func TestGetUserPermissionsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetUserPermissions(); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetItemPermissions(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		fmt.Fprint(w, `{"data":{"update":{"access":true}}}`)
	})
	perms, err := client.GetItemPermissions("sites", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions/me/sites" {
		t.Errorf("expected /permissions/me/sites, got %s", gotPath)
	}
	if _, ok := perms["update"]; !ok {
		t.Errorf("unexpected perms: %+v", perms)
	}
}

func TestGetItemPermissionsWithKey(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		fmt.Fprint(w, `{"data":{"update":{"access":true}}}`)
	})
	if _, err := client.GetItemPermissions("sites", "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions/me/sites/5" {
		t.Errorf("expected /permissions/me/sites/5, got %s", gotPath)
	}
}

func TestGetItemPermissionsValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetItemPermissions("", ""); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "collection must be provided") {
		t.Errorf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty collection")
	}
}

func TestCreatePermission(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":1,"collection":"sites","action":"read"}}`)
	})
	perm, err := client.CreatePermission(&Permission{Collection: "sites", Action: "read"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions" || gotMethod != "POST" {
		t.Errorf("expected POST /permissions, got %s %s", gotMethod, gotPath)
	}
	if perm == nil || perm.ID != 1 {
		t.Errorf("unexpected perm: %+v", perm)
	}
}

func TestCreatePermissionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreatePermission(&Permission{}, nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestCreatePermissions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})
	perms, err := client.CreatePermissions([]Permission{{Action: "read"}, {Action: "create"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions" || gotMethod != "POST" {
		t.Errorf("expected POST /permissions, got %s %s", gotMethod, gotPath)
	}
	if len(perms) != 2 {
		t.Errorf("unexpected perms: %+v", perms)
	}
}

func TestPatchPermission(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":42,"action":"update"}}`)
	})
	perm, err := client.PatchPermission(42, &Permission{Action: "update"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions/42" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /permissions/42, got %s %s", gotMethod, gotPath)
	}
	if perm == nil || perm.Action != "update" {
		t.Errorf("unexpected perm: %+v", perm)
	}
}

func TestPatchPermissions(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})
	perms, err := client.PatchPermissions([]int{1, 2}, &Permission{Action: "read"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /permissions, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"keys":[1,2]`) || !strings.Contains(gotBody, `"data"`) {
		t.Errorf("expected keys+data body, got %s", gotBody)
	}
	if len(perms) != 2 {
		t.Errorf("unexpected perms: %+v", perms)
	}
}

func TestPatchPermissionsBatch(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":1}]}`)
	})
	perms, err := client.PatchPermissionsBatch([]Permission{{ID: 1, Action: "read"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /permissions, got %s %s", gotMethod, gotPath)
	}
	if len(perms) != 1 {
		t.Errorf("unexpected perms: %+v", perms)
	}
}

func TestDeletePermission(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeletePermission(42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions/42" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /permissions/42, got %s %s", gotMethod, gotPath)
	}
}

func TestDeletePermissions(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeletePermissions([]int{1, 2}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/permissions" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /permissions, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `[1,2]`) {
		t.Errorf("expected raw keys body, got %s", gotBody)
	}
}

func TestPermissionsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetPermissions(nil); err == nil {
		t.Error("GetPermissions: expected request build error, got nil")
	}
}
