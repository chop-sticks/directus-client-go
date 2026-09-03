package directus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetFolders(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a","name":"one"},{"id":"b","name":"two"}]}`)
	})

	folders, err := client.GetFolders(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/folders" {
		t.Errorf("expected GET /folders, got %s %s", gotMethod, gotPath)
	}
	if len(folders) != 2 || folders[0].ID != "a" || folders[1].Name != "two" {
		t.Errorf("unexpected folders: %+v", folders)
	}
}

func TestGetFoldersHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetFolders(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetFoldersBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.GetFolders(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetFolder(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"one","parent":"root"}}`)
	})

	folder, err := client.GetFolder("a", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/folders/a" {
		t.Errorf("expected GET /folders/a, got %s %s", gotMethod, gotPath)
	}
	if folder == nil || folder.ID != "a" || folder.Parent != "root" {
		t.Errorf("unexpected folder: %+v", folder)
	}
}

func TestGetFolderValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.GetFolder("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetFolderHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetFolder("a", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetFolderBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetFolder("a", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateFolder(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"one"}}`)
	})

	folder, err := client.CreateFolder(&Folder{Name: "one"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/folders" {
		t.Errorf("expected POST /folders, got %s %s", gotMethod, gotPath)
	}
	if gotBody["name"] != "one" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if folder == nil || folder.Name != "one" {
		t.Errorf("unexpected folder: %+v", folder)
	}
}

func TestCreateFolderHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.CreateFolder(&Folder{Name: "one"}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestCreateFolderBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateFolder(&Folder{Name: "one"}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateFolders(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	folders, err := client.CreateFolders([]Folder{{Name: "one"}, {Name: "two"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/folders" {
		t.Errorf("expected POST /folders, got %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0]["name"] != "one" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(folders) != 2 {
		t.Errorf("unexpected folders: %+v", folders)
	}
}

func TestCreateFoldersHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.CreateFolders([]Folder{{Name: "one"}}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDeleteFolder(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteFolder("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/folders/a" {
		t.Errorf("expected DELETE /folders/a, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteFolderValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteFolder(""); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteFolderHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeleteFolder("a"); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDeleteFolders(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteFolders([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/folders" {
		t.Errorf("expected DELETE /folders, got %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
}

func TestDeleteFoldersHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeleteFolders([]string{"a"}); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestFoldersRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetFolders(nil); err == nil {
		t.Error("GetFolders: expected request build error, got nil")
	}
}
