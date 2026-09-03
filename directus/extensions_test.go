package directus

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestExtensionsRequestBuildError(t *testing.T) {
	client := badHostClient(t)

	if _, err := client.GetExtensions(); err == nil {
		t.Error("GetExtensions: expected request build error, got nil")
	}
	if _, err := client.PatchExtension("my-ext", map[string]any{}); err == nil {
		t.Error("PatchExtension: expected request build error, got nil")
	}
	if _, err := client.PatchBundleExtension("bundle", "my-ext", map[string]any{}); err == nil {
		t.Error("PatchBundleExtension: expected request build error, got nil")
	}
	if err := client.DeleteExtension("abc"); err == nil {
		t.Error("DeleteExtension: expected request build error, got nil")
	}
	if err := client.InstallRegistryExtension("abc", "1.0.0"); err == nil {
		t.Error("InstallRegistryExtension: expected request build error, got nil")
	}
	if err := client.UninstallRegistryExtension("abc"); err == nil {
		t.Error("UninstallRegistryExtension: expected request build error, got nil")
	}
	if _, err := client.GetRegistryExtensions(nil); err == nil {
		t.Error("GetRegistryExtensions: expected request build error, got nil")
	}
}

func TestGetExtensions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"my-ext","meta":{"enabled":true},"schema":{"type":"interface","name":"my-ext"}}]}`)
	})

	extensions, err := client.GetExtensions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/extensions/" {
		t.Errorf("expected GET /extensions/, got %s %s", gotMethod, gotPath)
	}
	if len(extensions) != 1 {
		t.Fatalf("expected 1 extension, got %d", len(extensions))
	}
	if extensions[0].ID != "my-ext" || !extensions[0].Meta.Enabled {
		t.Errorf("unexpected extension: %+v", extensions[0])
	}
	if extensions[0].Schema["name"] != "my-ext" {
		t.Errorf("unexpected schema: %+v", extensions[0].Schema)
	}
}

func TestGetExtensionsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetExtensions(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetExtensionsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.GetExtensions(); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchExtension(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"my-ext","meta":{"enabled":false}}}`)
	})

	ext, err := client.PatchExtension("my-ext", map[string]any{"meta": map[string]any{"enabled": false}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" || gotPath != "/extensions/my-ext" {
		t.Errorf("expected PATCH /extensions/my-ext, got %s %s", gotMethod, gotPath)
	}
	if ext.ID != "my-ext" || ext.Meta.Enabled {
		t.Errorf("unexpected extension: %+v", ext)
	}
}

func TestPatchExtensionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchExtension("", map[string]any{}); err == nil || !strings.Contains(err.Error(), "name must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchExtensionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.PatchExtension("my-ext", map[string]any{}); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestPatchExtensionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.PatchExtension("my-ext", map[string]any{}); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchBundleExtension(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"child","meta":{"enabled":true}}}`)
	})

	ext, err := client.PatchBundleExtension("bundle", "child", map[string]any{"meta": map[string]any{"enabled": true}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" || gotPath != "/extensions/bundle/child" {
		t.Errorf("expected PATCH /extensions/bundle/child, got %s %s", gotMethod, gotPath)
	}
	if ext.ID != "child" || !ext.Meta.Enabled {
		t.Errorf("unexpected extension: %+v", ext)
	}
}

func TestPatchBundleExtensionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchBundleExtension("", "child", map[string]any{}); err == nil || !strings.Contains(err.Error(), "bundle and name must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if _, err := client.PatchBundleExtension("bundle", "", map[string]any{}); err == nil {
		t.Error("expected validation error for empty name, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteExtension(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteExtension("abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/extensions/abc" {
		t.Errorf("expected DELETE /extensions/abc, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteExtensionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if err := client.DeleteExtension(""); err == nil || !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteExtensionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if err := client.DeleteExtension("abc"); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestInstallRegistryExtension(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.InstallRegistryExtension("abc", "1.0.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/extensions/registry/install" {
		t.Errorf("expected POST /extensions/registry/install, got %s %s", gotMethod, gotPath)
	}
}

func TestInstallRegistryExtensionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if err := client.InstallRegistryExtension("", "1.0.0"); err == nil || !strings.Contains(err.Error(), "extensionID and version must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if err := client.InstallRegistryExtension("abc", ""); err == nil {
		t.Error("expected validation error for empty version, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestInstallRegistryExtensionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.InstallRegistryExtension("abc", "1.0.0"); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestUninstallRegistryExtension(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.UninstallRegistryExtension("abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/extensions/registry/uninstall/abc" {
		t.Errorf("expected DELETE /extensions/registry/uninstall/abc, got %s %s", gotMethod, gotPath)
	}
}

func TestUninstallRegistryExtensionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if err := client.UninstallRegistryExtension(""); err == nil || !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("expected validation error, got %v", err)
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetRegistryExtensions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"abc","name":"Cool Extension"}]}`)
	})

	extensions, err := client.GetRegistryExtensions(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/extensions/registry" {
		t.Errorf("expected GET /extensions/registry, got %s %s", gotMethod, gotPath)
	}
	if len(extensions) != 1 || extensions[0]["name"] != "Cool Extension" {
		t.Errorf("unexpected extensions: %+v", extensions)
	}
}

func TestGetRegistryExtensionsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetRegistryExtensions(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetRegistryExtensionsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.GetRegistryExtensions(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}
