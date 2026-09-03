package directus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetAsset(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "RAWBYTES")
	})

	data, err := client.GetAsset("abc", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/assets/abc" {
		t.Errorf("expected GET /assets/abc, got %s %s", gotMethod, gotPath)
	}
	if string(data) != "RAWBYTES" {
		t.Errorf("unexpected data: %s", data)
	}
}

func TestGetAssetValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.GetAsset("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetAssetHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetAsset("abc", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestDownloadFilesZip(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ZIPBYTES")
	})

	data, err := client.DownloadFilesZip([]string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/assets/files/" {
		t.Errorf("expected POST /assets/files/, got %s %s", gotMethod, gotPath)
	}
	ids, ok := gotBody["ids"].([]any)
	if !ok || len(ids) != 2 || ids[0] != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if string(data) != "ZIPBYTES" {
		t.Errorf("unexpected data: %s", data)
	}
}

func TestDownloadFilesZipHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.DownloadFilesZip([]string{"a"}); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDownloadFolderZip(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ZIPBYTES")
	})

	data, err := client.DownloadFolderZip("abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/assets/folder/abc" {
		t.Errorf("expected POST /assets/folder/abc, got %s %s", gotMethod, gotPath)
	}
	if string(data) != "ZIPBYTES" {
		t.Errorf("unexpected data: %s", data)
	}
}

func TestDownloadFolderZipValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.DownloadFolderZip(""); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDownloadFolderZipHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.DownloadFolderZip("abc"); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestAssetsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetAsset("abc", nil); err == nil {
		t.Error("GetAsset: expected request build error, got nil")
	}
	if _, err := client.DownloadFilesZip([]string{"a"}); err == nil {
		t.Error("DownloadFilesZip: expected request build error, got nil")
	}
	if _, err := client.DownloadFolderZip("abc"); err == nil {
		t.Error("DownloadFolderZip: expected request build error, got nil")
	}
}
