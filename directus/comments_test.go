package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetComments(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"c1","comment":"hello","item":"5"}]}`)
	})

	comments, err := client.GetComments(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/comments" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(comments) != 1 || comments[0].Comment != "hello" || comments[0].Item != "5" {
		t.Errorf("unexpected comments: %+v", comments)
	}
}

func TestGetCommentsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetComments(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetCommentsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetComments(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetComment(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"abc","comment":"hi"}}`)
	})

	comment, err := client.GetComment("abc", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/comments/abc" {
		t.Errorf("expected path /comments/abc, got %s", gotPath)
	}
	if comment.ID != "abc" || comment.Comment != "hi" {
		t.Errorf("unexpected comment: %+v", comment)
	}
}

func TestGetCommentValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetComment("", nil); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetCommentBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetComment("abc", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateComment(t *testing.T) {
	var gotMethod, gotPath string
	var gotReq Comment
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotReq)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"new","comment":"hi"}}`)
	})

	comment, err := client.CreateComment(&Comment{Comment: "hi", Item: "5"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/comments" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotReq.Comment != "hi" || gotReq.Item != "5" {
		t.Errorf("unexpected request: %+v", gotReq)
	}
	if comment.ID != "new" {
		t.Errorf("expected id new, got %q", comment.ID)
	}
}

func TestCreateCommentBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.CreateComment(&Comment{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateComments(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	comments, err := client.CreateComments([]Comment{{Comment: "a"}, {Comment: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/comments" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}

func TestPatchComment(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"z","comment":"edited"}}`)
	})

	comment, err := client.PatchComment("z", &Comment{Comment: "edited"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/comments/z" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if comment.Comment != "edited" {
		t.Errorf("unexpected comment: %+v", comment)
	}
}

func TestPatchCommentValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchComment("", &Comment{}, nil); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteComment(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteComment("q"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/comments/q" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
}

func TestDeleteCommentValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteComment(""); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteComments(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []string `json:"keys"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteComments([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/comments" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Keys[0] != "a" {
		t.Errorf("expected keys wrapped in object, got %+v", gotBody)
	}
}

func TestCommentsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetComments(nil); err == nil {
		t.Error("GetComments: expected request build error, got nil")
	}
	if _, err := client.CreateComment(&Comment{}, nil); err == nil {
		t.Error("CreateComment: expected request build error, got nil")
	}
	if err := client.DeleteComments([]string{"a"}); err == nil {
		t.Error("DeleteComments: expected request build error, got nil")
	}
}
