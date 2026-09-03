package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetOperations(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a","key":"one"},{"id":"b"}]}`)
	})

	ops, err := client.GetOperations(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/operations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(ops) != 2 || ops[0].Key != "one" {
		t.Errorf("unexpected operations: %+v", ops)
	}
}

func TestGetOperationsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetOperations(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetOperationsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetOperations(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetOperation(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","position_x":3}}`)
	})

	op, err := client.GetOperation("a", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/operations/a" {
		t.Errorf("expected path /operations/a, got %s", gotPath)
	}
	if op.PositionX != 3 {
		t.Errorf("unexpected operation: %+v", op)
	}
}

func TestGetOperationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetOperation("", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetOperationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetOperation("missing", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetOperationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.GetOperation("a", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateOperation(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Operation
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","key":"one"}}`)
	})

	op, err := client.CreateOperation(&Operation{Key: "one"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/operations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotBody.Key != "one" || op.ID != "a" {
		t.Errorf("unexpected body/op: %+v %+v", gotBody, op)
	}
}

func TestCreateOperationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.CreateOperation(&Operation{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateOperations(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	ops, err := client.CreateOperations([]Operation{{Key: "one"}, {Key: "two"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/operations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(ops) != 2 {
		t.Errorf("expected 2 operations, got %d", len(ops))
	}
}

func TestPatchOperation(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","key":"updated"}}`)
	})

	op, err := client.PatchOperation("a", &Operation{Key: "updated"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/operations/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if op.Key != "updated" {
		t.Errorf("unexpected op: %+v", op)
	}
}

func TestPatchOperationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchOperation("", &Operation{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchOperations(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []string  `json:"keys"`
		Data Operation `json:"data"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	ops, err := client.PatchOperations([]string{"a", "b"}, &Operation{Type: "log"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/operations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Data.Type != "log" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(ops) != 2 {
		t.Errorf("expected 2 operations, got %d", len(ops))
	}
}

func TestPatchOperationsBatch(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []Operation
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	ops, err := client.PatchOperationsBatch([]Operation{{ID: "a"}, {ID: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/operations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0].ID != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(ops) != 2 {
		t.Errorf("expected 2 operations, got %d", len(ops))
	}
}

func TestDeleteOperation(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteOperation("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/operations/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
}

func TestDeleteOperationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteOperation(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteOperations(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteOperations([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/operations" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != "a" {
		t.Errorf("expected raw keys body [a b], got %+v", gotBody)
	}
}

func TestOperationsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetOperations(nil); err == nil {
		t.Error("GetOperations: expected request build error, got nil")
	}
}
