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

func TestGetFiles(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"abc","title":"one"},{"id":"def","title":"two"}]}`)
	})

	files, err := client.GetFiles(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/files" {
		t.Errorf("expected GET /files, got %s %s", gotMethod, gotPath)
	}
	if len(files) != 2 || files[0].ID != "abc" || files[1].Title != "two" {
		t.Errorf("unexpected files: %+v", files)
	}
}

func TestGetFilesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetFiles(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetFilesBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[`)
	})
	if _, err := client.GetFiles(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetFile(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"abc","title":"one","width":640,"height":480,"tags":["a","b"]}}`)
	})

	file, err := client.GetFile("abc", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/files/abc" {
		t.Errorf("expected GET /files/abc, got %s %s", gotMethod, gotPath)
	}
	if file == nil || file.ID != "abc" || file.Width != 640 || file.Height != 480 {
		t.Errorf("unexpected file: %+v", file)
	}
	if len(file.Tags) != 2 || file.Tags[0] != "a" {
		t.Errorf("unexpected tags: %+v", file.Tags)
	}
}

func TestGetFileValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.GetFile("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetFileHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetFile("abc", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetFileBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetFile("abc", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestImportFile(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"abc","title":"imported"}}`)
	})

	file, err := client.ImportFile("https://example.com/x.png", &File{Title: "imported"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/files/import" {
		t.Errorf("expected POST /files/import, got %s %s", gotMethod, gotPath)
	}
	if gotBody["url"] != "https://example.com/x.png" {
		t.Errorf("unexpected url in body: %+v", gotBody)
	}
	if _, ok := gotBody["data"]; !ok {
		t.Errorf("expected data key in body: %+v", gotBody)
	}
	if file == nil || file.Title != "imported" {
		t.Errorf("unexpected file: %+v", file)
	}
}

func TestImportFileValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.ImportFile("", nil, nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "url must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty url")
	}
}

func TestImportFileHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.ImportFile("https://example.com/x.png", nil, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestImportFileBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.ImportFile("https://example.com/x.png", nil, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchFile(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"abc","title":"renamed"}}`)
	})

	file, err := client.PatchFile("abc", &File{Title: "renamed"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" || gotPath != "/files/abc" {
		t.Errorf("expected PATCH /files/abc, got %s %s", gotMethod, gotPath)
	}
	if file == nil || file.Title != "renamed" {
		t.Errorf("unexpected file: %+v", file)
	}
}

func TestPatchFileValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchFile("", nil, nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchFileHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchFile("abc", nil, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestPatchFiles(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	files, err := client.PatchFiles([]string{"a", "b"}, &File{Title: "x"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" || gotPath != "/files" {
		t.Errorf("expected PATCH /files, got %s %s", gotMethod, gotPath)
	}
	keys, ok := gotBody["keys"].([]any)
	if !ok || len(keys) != 2 {
		t.Errorf("expected keys array in body: %+v", gotBody)
	}
	if _, ok := gotBody["data"]; !ok {
		t.Errorf("expected data key in body: %+v", gotBody)
	}
	if len(files) != 2 {
		t.Errorf("unexpected files: %+v", files)
	}
}

func TestPatchFilesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchFiles([]string{"a"}, nil, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestPatchFilesBatch(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	files, err := client.PatchFilesBatch([]File{{ID: "a", Title: "x"}, {ID: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" || gotPath != "/files" {
		t.Errorf("expected PATCH /files, got %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0]["id"] != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(files) != 2 {
		t.Errorf("unexpected files: %+v", files)
	}
}

func TestPatchFilesBatchHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.PatchFilesBatch([]File{{ID: "a"}}, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDeleteFile(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteFile("abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/files/abc" {
		t.Errorf("expected DELETE /files/abc, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteFileValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteFile(""); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteFileHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeleteFile("abc"); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestDeleteFiles(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteFiles([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/files" {
		t.Errorf("expected DELETE /files, got %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
}

func TestDeleteFilesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeleteFiles([]string{"a"}); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestUploadFile(t *testing.T) {
	var gotPath, gotMethod, gotContentType, gotFilename, gotField, gotFileContent string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")

		mediaType, params, err := mime.ParseMediaType(gotContentType)
		if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
			t.Errorf("expected multipart content type, got %q", gotContentType)
		} else {
			mr := multipart.NewReader(r.Body, params["boundary"])
			for {
				part, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("reading part: %v", err)
				}
				data, _ := io.ReadAll(part)
				switch part.FormName() {
				case "file":
					gotFilename = part.FileName()
					gotFileContent = string(data)
				case "title":
					gotField = string(data)
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"abc","filename_download":"pic.png"}}`)
	})

	file, err := client.UploadFile(strings.NewReader("PNGDATA"), "pic.png", map[string]string{"title": "hello"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" || gotPath != "/files" {
		t.Errorf("expected POST /files, got %s %s", gotMethod, gotPath)
	}
	if gotFilename != "pic.png" {
		t.Errorf("expected file part filename pic.png, got %q", gotFilename)
	}
	if gotFileContent != "PNGDATA" {
		t.Errorf("expected file content PNGDATA, got %q", gotFileContent)
	}
	if gotField != "hello" {
		t.Errorf("expected title field hello, got %q", gotField)
	}
	if file == nil || file.ID != "abc" || file.FilenameDownload != "pic.png" {
		t.Errorf("unexpected file: %+v", file)
	}
}

func TestUploadFileValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.UploadFile(nil, "pic.png", nil, nil); err == nil {
		t.Error("expected validation error for nil data, got nil")
	}
	if _, err := client.UploadFile(strings.NewReader("x"), "", nil, nil); err == nil {
		t.Error("expected validation error for empty filename, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestUploadFileHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.UploadFile(strings.NewReader("x"), "pic.png", nil, nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestUploadFileBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.UploadFile(strings.NewReader("x"), "pic.png", nil, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestFilesRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetFiles(nil); err == nil {
		t.Error("GetFiles: expected request build error, got nil")
	}
	if _, err := client.UploadFile(strings.NewReader("x"), "pic.png", nil, nil); err == nil {
		t.Error("UploadFile: expected request build error, got nil")
	}
}
