package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestVersionsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetContentVersions(nil); err == nil {
		t.Error("GetContentVersions: expected request build error, got nil")
	}
	if _, err := client.GetContentVersion("v1", nil); err == nil {
		t.Error("GetContentVersion: expected request build error, got nil")
	}
	if _, err := client.CreateContentVersion(&Version{}, nil); err == nil {
		t.Error("CreateContentVersion: expected request build error, got nil")
	}
	if _, err := client.CreateContentVersions([]Version{{}}, nil); err == nil {
		t.Error("CreateContentVersions: expected request build error, got nil")
	}
	if _, err := client.PatchContentVersion("v1", &Version{}, nil); err == nil {
		t.Error("PatchContentVersion: expected request build error, got nil")
	}
	if _, err := client.PatchContentVersions([]string{"v1"}, &Version{}, nil); err == nil {
		t.Error("PatchContentVersions: expected request build error, got nil")
	}
	if _, err := client.PatchContentVersionsBatch([]Version{{}}, nil); err == nil {
		t.Error("PatchContentVersionsBatch: expected request build error, got nil")
	}
	if err := client.DeleteContentVersion("v1"); err == nil {
		t.Error("DeleteContentVersion: expected request build error, got nil")
	}
	if err := client.DeleteContentVersions([]string{"v1"}); err == nil {
		t.Error("DeleteContentVersions: expected request build error, got nil")
	}
	if _, err := client.SaveToContentVersion("v1", map[string]any{}); err == nil {
		t.Error("SaveToContentVersion: expected request build error, got nil")
	}
	if _, err := client.CompareContentVersion("v1"); err == nil {
		t.Error("CompareContentVersion: expected request build error, got nil")
	}
	if _, err := client.PromoteContentVersion("v1", "hash", nil); err == nil {
		t.Error("PromoteContentVersion: expected request build error, got nil")
	}
}

func TestGetContentVersions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"v1","key":"draft"},{"id":"v2","key":"wip"}]}`)
	})

	items, err := client.GetContentVersions(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/versions" {
		t.Errorf("expected GET /versions, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 2 || items[0].ID != "v1" || items[0].Key != "draft" {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestGetContentVersionsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetContentVersions(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetContentVersionsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetContentVersions(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetContentVersion(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"v1","key":"draft","hash":"abc"}}`)
	})

	item, err := client.GetContentVersion("v1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/versions/v1" {
		t.Errorf("expected GET /versions/v1, got %s %s", gotMethod, gotPath)
	}
	if item == nil || item.ID != "v1" || item.Hash != "abc" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestGetContentVersionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.GetContentVersion("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetContentVersionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetContentVersion("missing", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetContentVersionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetContentVersion("v1", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateContentVersion(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody Version
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"v1","key":"draft"}}`)
	})

	item, err := client.CreateContentVersion(&Version{Key: "draft"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/versions" {
		t.Errorf("expected POST /versions, got %s %s", gotMethod, gotPath)
	}
	if gotBody.Key != "draft" {
		t.Errorf("unexpected request body: %+v", gotBody)
	}
	if item == nil || item.ID != "v1" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestCreateContentVersionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.CreateContentVersion(&Version{}, nil); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestCreateContentVersionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateContentVersion(&Version{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateContentVersions(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"v1"},{"id":"v2"}]}`)
	})

	items, err := client.CreateContentVersions([]Version{{Key: "a"}, {Key: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/versions" {
		t.Errorf("expected POST /versions, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestPatchContentVersion(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"v1","name":"Draft"}}`)
	})

	item, err := client.PatchContentVersion("v1", &Version{Name: "Draft"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/versions/v1" {
		t.Errorf("expected PATCH /versions/v1, got %s %s", gotMethod, gotPath)
	}
	if item == nil || item.Name != "Draft" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestPatchContentVersionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchContentVersion("", &Version{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchContentVersions(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"v1"},{"id":"v2"}]}`)
	})

	items, err := client.PatchContentVersions([]string{"v1", "v2"}, &Version{Name: "Draft"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/versions" {
		t.Errorf("expected PATCH /versions, got %s %s", gotMethod, gotPath)
	}
	if _, ok := gotBody["keys"]; !ok {
		t.Errorf("expected keys in body, got %+v", gotBody)
	}
	if _, ok := gotBody["data"]; !ok {
		t.Errorf("expected data in body, got %+v", gotBody)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestPatchContentVersionsBatch(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"v1"}]}`)
	})

	items, err := client.PatchContentVersionsBatch([]Version{{ID: "v1", Name: "Draft"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/versions" {
		t.Errorf("expected PATCH /versions, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestDeleteContentVersion(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteContentVersion("v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/versions/v1" {
		t.Errorf("expected DELETE /versions/v1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteContentVersionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteContentVersion(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteContentVersionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if err := client.DeleteContentVersion("missing"); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestDeleteContentVersions(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteContentVersions([]string{"v1", "v2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/versions" {
		t.Errorf("expected DELETE /versions, got %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != "v1" {
		t.Errorf("expected raw keys array body, got %+v", gotBody)
	}
}

func TestSaveToContentVersion(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"title":"Hello","status":"draft"}}`)
	})

	res, err := client.SaveToContentVersion("v1", map[string]any{"title": "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/versions/v1/save" {
		t.Errorf("expected POST /versions/v1/save, got %s %s", gotMethod, gotPath)
	}
	if gotBody["title"] != "Hello" {
		t.Errorf("unexpected request body: %+v", gotBody)
	}
	if res["title"] != "Hello" || res["status"] != "draft" {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestSaveToContentVersionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.SaveToContentVersion("", map[string]any{}); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestSaveToContentVersionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.SaveToContentVersion("v1", map[string]any{}); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestSaveToContentVersionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.SaveToContentVersion("v1", map[string]any{}); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCompareContentVersion(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"outdated":true,"mainHash":"abc","current":{"title":"c"},"main":{"title":"m"}}}`)
	})

	res, err := client.CompareContentVersion("v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/versions/v1/compare" {
		t.Errorf("expected GET /versions/v1/compare, got %s %s", gotMethod, gotPath)
	}
	if res == nil || !res.Outdated || res.MainHash != "abc" {
		t.Errorf("unexpected result: %+v", res)
	}
	if res.Current["title"] != "c" || res.Main["title"] != "m" {
		t.Errorf("unexpected current/main: %+v", res)
	}
}

func TestCompareContentVersionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.CompareContentVersion(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestCompareContentVersionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.CompareContentVersion("missing"); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestCompareContentVersionBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CompareContentVersion("v1"); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPromoteContentVersion(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":"42"}`)
	})

	res, err := client.PromoteContentVersion("v1", "abc", []string{"title", "body"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/versions/v1/promote" {
		t.Errorf("expected POST /versions/v1/promote, got %s %s", gotMethod, gotPath)
	}
	if gotBody["mainHash"] != "abc" {
		t.Errorf("expected mainHash in body, got %+v", gotBody)
	}
	if _, ok := gotBody["fields"]; !ok {
		t.Errorf("expected fields in body, got %+v", gotBody)
	}
	if res != "42" {
		t.Errorf("unexpected result: %v", res)
	}
}

func TestPromoteContentVersionOmitsFields(t *testing.T) {
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":"42"}`)
	})

	if _, err := client.PromoteContentVersion("v1", "abc", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotBody["fields"]; ok {
		t.Errorf("expected fields omitted from body, got %+v", gotBody)
	}
	if gotBody["mainHash"] != "abc" {
		t.Errorf("expected mainHash in body, got %+v", gotBody)
	}
}

func TestPromoteContentVersionValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PromoteContentVersion("", "abc", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPromoteContentVersionHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.PromoteContentVersion("v1", "abc", nil); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}
