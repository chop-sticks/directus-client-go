package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetTranslations(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a","language":"en-US","key":"hello","value":"Hello"}]}`)
	})

	tr, err := client.GetTranslations(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/translations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(tr) != 1 || tr[0].Key != "hello" || tr[0].Value != "Hello" {
		t.Errorf("unexpected translations: %+v", tr)
	}
}

func TestGetTranslationsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetTranslations(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetTranslationsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetTranslations(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetTranslation(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"abc","key":"hi"}}`)
	})

	tr, err := client.GetTranslation("abc", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/translations/abc" {
		t.Errorf("expected path /translations/abc, got %s", gotPath)
	}
	if tr.ID != "abc" || tr.Key != "hi" {
		t.Errorf("unexpected translation: %+v", tr)
	}
}

func TestGetTranslationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetTranslation("", nil); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetTranslationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetTranslation("abc", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateTranslation(t *testing.T) {
	var gotMethod, gotPath string
	var gotReq Translation
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotReq)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"x","key":"hi"}}`)
	})

	tr, err := client.CreateTranslation(&Translation{Key: "hi", Value: "Hi"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/translations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotReq.Key != "hi" || gotReq.Value != "Hi" {
		t.Errorf("unexpected request: %+v", gotReq)
	}
	if tr.ID != "x" {
		t.Errorf("expected id x, got %q", tr.ID)
	}
}

func TestCreateTranslationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.CreateTranslation(&Translation{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateTranslations(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	tr, err := client.CreateTranslations([]Translation{{Key: "a"}, {Key: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/translations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(tr) != 2 {
		t.Fatalf("expected 2 translations, got %d", len(tr))
	}
}

func TestPatchTranslation(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"z","value":"new"}}`)
	})

	tr, err := client.PatchTranslation("z", &Translation{Value: "new"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/translations/z" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if tr.Value != "new" {
		t.Errorf("unexpected translation: %+v", tr)
	}
}

func TestPatchTranslationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchTranslation("", &Translation{}, nil); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchTranslations(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []string    `json:"keys"`
		Data Translation `json:"data"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	tr, err := client.PatchTranslations([]string{"a", "b"}, &Translation{Value: "v"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/translations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Data.Value != "v" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(tr) != 2 {
		t.Fatalf("expected 2 translations, got %d", len(tr))
	}
}

func TestPatchTranslationsBatch(t *testing.T) {
	var gotMethod, gotPath string
	var gotItems []Translation
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotItems)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	tr, err := client.PatchTranslationsBatch([]Translation{{ID: "a"}, {ID: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/translations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotItems) != 2 || gotItems[0].ID != "a" {
		t.Errorf("unexpected items: %+v", gotItems)
	}
	if len(tr) != 2 {
		t.Fatalf("expected 2 translations, got %d", len(tr))
	}
}

func TestDeleteTranslation(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteTranslation("q"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/translations/q" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
}

func TestDeleteTranslationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteTranslation(""); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteTranslations(t *testing.T) {
	var gotMethod, gotPath string
	var gotKeys []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotKeys)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteTranslations([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/translations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotKeys) != 2 || gotKeys[0] != "a" {
		t.Errorf("expected raw keys array, got %+v", gotKeys)
	}
}

func TestTranslationsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetTranslations(nil); err == nil {
		t.Error("GetTranslations: expected request build error, got nil")
	}
	if _, err := client.CreateTranslation(&Translation{}, nil); err == nil {
		t.Error("CreateTranslation: expected request build error, got nil")
	}
}
