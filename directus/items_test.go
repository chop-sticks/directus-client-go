package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetItems(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1,"title":"a"},{"id":2,"title":"b"}]}`)
	})

	items, err := client.GetItems("articles", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	if len(items) != 2 || items[0]["title"] != "a" {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestGetItemsValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.GetItems("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "collection must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetItemsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetItems("articles", nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetItemsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetItems("articles", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetItem(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":7,"title":"hello"}}`)
	})

	item, err := client.GetItem("articles", "7", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/items/articles/7" {
		t.Errorf("expected path /items/articles/7, got %s", gotPath)
	}
	if item["title"] != "hello" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestGetItemValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	cases := []struct{ collection, key string }{
		{"", "7"},
		{"articles", ""},
		{"", ""},
	}
	for _, tc := range cases {
		if _, err := client.GetItem(tc.collection, tc.key, nil); err == nil {
			t.Errorf("collection=%q key=%q: expected validation error, got nil", tc.collection, tc.key)
		} else if !strings.Contains(err.Error(), "collection and key must be provided") {
			t.Errorf("unexpected error message: %v", err)
		}
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetItemHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetItem("articles", "missing", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetItemBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetItem("articles", "7", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateItem(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":9,"title":"new"}}`)
	})

	item, err := client.CreateItem("articles", map[string]any{"title": "new"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	if gotBody["title"] != "new" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if item["title"] != "new" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestCreateItemValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.CreateItem("", map[string]any{"title": "x"}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestCreateItemHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.CreateItem("articles", map[string]any{"title": "x"}, nil); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestCreateItemBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateItem("articles", map[string]any{"title": "x"}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateItems(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	items, err := client.CreateItems("articles", []map[string]any{{"title": "a"}, {"title": "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	if len(gotBody) != 2 {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestCreateItemsValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.CreateItems("", nil, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestCreateItemsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.CreateItems("articles", []map[string]any{{"title": "a"}}, nil); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestCreateItemsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.CreateItems("articles", []map[string]any{{"title": "a"}}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchItem(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":7,"title":"upd"}}`)
	})

	item, err := client.PatchItem("articles", "7", map[string]any{"title": "upd"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" {
		t.Errorf("expected PATCH, got %s", gotMethod)
	}
	if gotPath != "/items/articles/7" {
		t.Errorf("expected path /items/articles/7, got %s", gotPath)
	}
	if gotBody["title"] != "upd" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if item["title"] != "upd" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestPatchItemValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	cases := []struct{ collection, key string }{
		{"", "7"},
		{"articles", ""},
	}
	for _, tc := range cases {
		if _, err := client.PatchItem(tc.collection, tc.key, map[string]any{"a": 1}, nil); err == nil {
			t.Errorf("collection=%q key=%q: expected validation error, got nil", tc.collection, tc.key)
		}
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchItemHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchItem("articles", "7", map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestPatchItemBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.PatchItem("articles", "7", map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchItems(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	items, err := client.PatchItems("articles", []any{1, 2}, map[string]any{"status": "published"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" {
		t.Errorf("expected PATCH, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	keys, ok := gotBody["keys"].([]any)
	if !ok || len(keys) != 2 {
		t.Errorf("unexpected keys in body: %+v", gotBody)
	}
	data, ok := gotBody["data"].(map[string]any)
	if !ok || data["status"] != "published" {
		t.Errorf("unexpected data in body: %+v", gotBody)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestPatchItemsValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchItems("", []any{1}, map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchItemsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchItems("articles", []any{1}, map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestPatchItemsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.PatchItems("articles", []any{1}, map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchItemsBatch(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	items, err := client.PatchItemsBatch("articles", []map[string]any{{"id": 1, "title": "a"}, {"id": 2, "title": "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" {
		t.Errorf("expected PATCH, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	if len(gotBody) != 2 {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestPatchItemsBatchValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchItemsBatch("", nil, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchItemsBatchHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchItemsBatch("articles", []map[string]any{{"id": 1}}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestPatchItemsBatchBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.PatchItemsBatch("articles", []map[string]any{{"id": 1}}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestDeleteItem(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteItem("articles", "7"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/items/articles/7" {
		t.Errorf("expected path /items/articles/7, got %s", gotPath)
	}
}

func TestDeleteItemValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	cases := []struct{ collection, key string }{
		{"", "7"},
		{"articles", ""},
	}
	for _, tc := range cases {
		if err := client.DeleteItem(tc.collection, tc.key); err == nil {
			t.Errorf("collection=%q key=%q: expected validation error, got nil", tc.collection, tc.key)
		}
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteItemHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeleteItem("articles", "7"); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDeleteItems(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteItems("articles", []any{1, 2, 3}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	keys, ok := gotBody["keys"].([]any)
	if !ok || len(keys) != 3 {
		t.Errorf("unexpected keys in body: %+v", gotBody)
	}
}

func TestDeleteItemsValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteItems("", []any{1}); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteItemsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeleteItems("articles", []any{1}); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDeleteItemsByQuery(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteItemsByQuery("articles", map[string]any{"filter": map[string]any{"status": "draft"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	if _, ok := gotBody["query"].(map[string]any); !ok {
		t.Errorf("unexpected query in body: %+v", gotBody)
	}
}

func TestDeleteItemsByQueryValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteItemsByQuery("", map[string]any{"filter": nil}); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteItemsByQueryHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeleteItemsByQuery("articles", map[string]any{"filter": nil}); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetSingleton(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":1,"site_name":"Acme"}}`)
	})

	item, err := client.GetSingleton("globals", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/items/globals" {
		t.Errorf("expected path /items/globals, got %s", gotPath)
	}
	if item["site_name"] != "Acme" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestGetSingletonValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.GetSingleton("", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetSingletonHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetSingleton("globals", nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetSingletonBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetSingleton("globals", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchSingleton(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":1,"site_name":"New"}}`)
	})

	item, err := client.PatchSingleton("globals", map[string]any{"site_name": "New"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" {
		t.Errorf("expected PATCH, got %s", gotMethod)
	}
	if gotPath != "/items/globals" {
		t.Errorf("expected path /items/globals, got %s", gotPath)
	}
	if gotBody["site_name"] != "New" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if item["site_name"] != "New" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestPatchSingletonValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchSingleton("", map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchSingletonHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchSingleton("globals", map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestPatchSingletonBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.PatchSingleton("globals", map[string]any{"a": 1}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestAggregate(t *testing.T) {
	var gotPath, gotMethod, gotQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"count":{"id":42}}]}`)
	})

	q := &Query{
		Aggregate: map[string]any{"count": "id"},
		GroupBy:   []string{"status"},
	}
	result, err := client.Aggregate("articles", q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/items/articles" {
		t.Errorf("expected path /items/articles, got %s", gotPath)
	}
	if !strings.Contains(gotQuery, "aggregate") || !strings.Contains(gotQuery, "groupBy") {
		t.Errorf("expected aggregate and groupBy in query, got %s", gotQuery)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

func TestAggregateValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.Aggregate("", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestAggregateHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.Aggregate("articles", nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestAggregateBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.Aggregate("articles", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}
