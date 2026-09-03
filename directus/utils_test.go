package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func TestClearCache(t *testing.T) {
	var gotPath, gotRawQuery, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotRawQuery = r.URL.RawQuery
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.ClearCache(false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/utils/cache/clear" {
		t.Errorf("expected POST /utils/cache/clear, got %s %s", gotMethod, gotPath)
	}
	if gotRawQuery != "" {
		t.Errorf("expected no query, got %q", gotRawQuery)
	}
}

func TestClearCacheSystem(t *testing.T) {
	var gotRawQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.ClearCache(true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotRawQuery != "system" {
		t.Errorf("expected system query, got %q", gotRawQuery)
	}
}

func TestClearCacheHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.ClearCache(false); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestRandomString(t *testing.T) {
	var gotPath, gotRawQuery, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotRawQuery = r.URL.RawQuery
		gotMethod = r.Method
		fmt.Fprint(w, `{"data":"abc123"}`)
	})
	s, err := client.RandomString(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/utils/random/string" {
		t.Errorf("expected GET /utils/random/string, got %s %s", gotMethod, gotPath)
	}
	if gotRawQuery != "" {
		t.Errorf("expected no query, got %q", gotRawQuery)
	}
	if s != "abc123" {
		t.Errorf("unexpected decode: %q", s)
	}
}

func TestRandomStringWithLength(t *testing.T) {
	var gotRawQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		fmt.Fprint(w, `{"data":"xxxxx"}`)
	})
	if _, err := client.RandomString(5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotRawQuery != "length=5" {
		t.Errorf("expected length=5 query, got %q", gotRawQuery)
	}
}

func TestRandomStringHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.RandomString(0); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestRandomStringBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad}`)
	})
	if _, err := client.RandomString(0); err == nil {
		t.Error("expected JSON decode error")
	}
}

func TestUtilitySort(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.UtilitySort("articles", 5, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/utils/sort/articles" {
		t.Errorf("expected POST /utils/sort/articles, got %s %s", gotMethod, gotPath)
	}
	if gotBody["item"] != float64(5) || gotBody["to"] != float64(10) {
		t.Errorf("unexpected body: %v", gotBody)
	}
}

func TestUtilitySortValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.UtilitySort("", 1, 2); err == nil {
		t.Error("expected validation error for empty collection")
	}
	if called {
		t.Error("expected no HTTP call when validation fails")
	}
}

func TestUtilitySortHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.UtilitySort("articles", 1, 2); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestUtilsExport(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})
	query := map[string]any{"limit": float64(10)}
	file := map[string]any{"filename_download": "out.csv"}
	if err := client.UtilsExport("articles", "csv", query, file); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/utils/export/articles" {
		t.Errorf("expected POST /utils/export/articles, got %s %s", gotMethod, gotPath)
	}
	if gotBody["format"] != "csv" {
		t.Errorf("expected format csv, got %v", gotBody["format"])
	}
	if _, ok := gotBody["query"]; !ok {
		t.Error("expected query in body")
	}
	if _, ok := gotBody["file"]; !ok {
		t.Error("expected file in body")
	}
}

func TestUtilsExportValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.UtilsExport("", "csv", nil, nil); err == nil {
		t.Error("expected validation error for empty collection")
	}
	if err := client.UtilsExport("articles", "", nil, nil); err == nil {
		t.Error("expected validation error for empty format")
	}
	if called {
		t.Error("expected no HTTP call when validation fails")
	}
}

func TestUtilsExportHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.UtilsExport("articles", "csv", nil, nil); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestUtilsImport(t *testing.T) {
	var gotPath, gotMethod, gotContentType, gotFileName, gotFileContent string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")

		mediaType, params, _ := mime.ParseMediaType(gotContentType)
		if mediaType == "multipart/form-data" {
			mr := multipart.NewReader(r.Body, params["boundary"])
			for {
				part, err := mr.NextPart()
				if err != nil {
					break
				}
				if part.FormName() == "file" {
					gotFileName = part.FileName()
					content, _ := io.ReadAll(part)
					gotFileContent = string(content)
				}
			}
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.UtilsImport("articles", strings.NewReader("id,name\n1,foo\n"), "data.csv"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/utils/import/articles" {
		t.Errorf("expected POST /utils/import/articles, got %s %s", gotMethod, gotPath)
	}
	if !strings.HasPrefix(gotContentType, "multipart/form-data") {
		t.Errorf("expected multipart/form-data Content-Type, got %q", gotContentType)
	}
	if gotFileName != "data.csv" {
		t.Errorf("expected file part named data.csv, got %q", gotFileName)
	}
	if gotFileContent != "id,name\n1,foo\n" {
		t.Errorf("unexpected file content: %q", gotFileContent)
	}
}

func TestUtilsImportValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.UtilsImport("", strings.NewReader("x"), "f.csv"); err == nil {
		t.Error("expected validation error for empty collection")
	}
	if err := client.UtilsImport("articles", nil, "f.csv"); err == nil {
		t.Error("expected validation error for nil data")
	}
	if err := client.UtilsImport("articles", strings.NewReader("x"), ""); err == nil {
		t.Error("expected validation error for empty filename")
	}
	if called {
		t.Error("expected no HTTP call when validation fails")
	}
}

func TestUtilsImportHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.UtilsImport("articles", strings.NewReader("x"), "f.csv"); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestUtilsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.RandomString(0); err == nil {
		t.Error("RandomString: expected request build error")
	}
	if err := client.ClearCache(false); err == nil {
		t.Error("ClearCache: expected request build error")
	}
	if err := client.UtilsImport("articles", strings.NewReader("x"), "f.csv"); err == nil {
		t.Error("UtilsImport: expected request build error")
	}
}
