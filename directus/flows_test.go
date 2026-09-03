package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetFlows(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a","name":"one"},{"id":"b","name":"two"}]}`)
	})

	flows, err := client.GetFlows(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/flows" {
		t.Errorf("expected path /flows, got %s", gotPath)
	}
	if len(flows) != 2 || flows[0].ID != "a" || flows[1].Name != "two" {
		t.Errorf("unexpected flows: %+v", flows)
	}
}

func TestGetFlowsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetFlows(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetFlowsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetFlows(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetFlow(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"one","options":{"foo":"bar"}}}`)
	})

	flow, err := client.GetFlow("a", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/flows/a" {
		t.Errorf("expected path /flows/a, got %s", gotPath)
	}
	if flow.ID != "a" || flow.Options["foo"] != "bar" {
		t.Errorf("unexpected flow: %+v", flow)
	}
}

func TestGetFlowValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetFlow("", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetFlowHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetFlow("missing", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetFlowBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.GetFlow("a", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateFlow(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Flow
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"one"}}`)
	})

	flow, err := client.CreateFlow(&Flow{Name: "one"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/flows" {
		t.Errorf("expected path /flows, got %s", gotPath)
	}
	if gotBody.Name != "one" {
		t.Errorf("expected body name 'one', got %q", gotBody.Name)
	}
	if flow.ID != "a" {
		t.Errorf("unexpected flow: %+v", flow)
	}
}

func TestCreateFlowBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.CreateFlow(&Flow{Name: "one"}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateFlows(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	flows, err := client.CreateFlows([]Flow{{Name: "one"}, {Name: "two"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/flows" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(flows) != 2 {
		t.Errorf("expected 2 flows, got %d", len(flows))
	}
}

func TestPatchFlow(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"updated"}}`)
	})

	flow, err := client.PatchFlow("a", &Flow{Name: "updated"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/flows/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if flow.Name != "updated" {
		t.Errorf("unexpected flow: %+v", flow)
	}
}

func TestPatchFlowValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchFlow("", &Flow{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchFlows(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []string `json:"keys"`
		Data Flow     `json:"data"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	flows, err := client.PatchFlows([]string{"a", "b"}, &Flow{Status: "active"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/flows" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Keys[0] != "a" {
		t.Errorf("expected keys [a b], got %+v", gotBody.Keys)
	}
	if gotBody.Data.Status != "active" {
		t.Errorf("expected data status 'active', got %q", gotBody.Data.Status)
	}
	if len(flows) != 2 {
		t.Errorf("expected 2 flows, got %d", len(flows))
	}
}

func TestPatchFlowsBatch(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []Flow
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	flows, err := client.PatchFlowsBatch([]Flow{{ID: "a", Name: "one"}, {ID: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/flows" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0].ID != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(flows) != 2 {
		t.Errorf("expected 2 flows, got %d", len(flows))
	}
}

func TestDeleteFlow(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteFlow("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/flows/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
}

func TestDeleteFlowValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteFlow(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteFlows(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteFlows([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/flows" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != "a" {
		t.Errorf("expected raw keys body [a b], got %+v", gotBody)
	}
}

func TestTriggerFlowGet(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.Query().Get("foo")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `raw-response`)
	})

	body, err := client.TriggerFlow("GET", "a", map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/flows/trigger/a" {
		t.Errorf("expected path /flows/trigger/a, got %s", gotPath)
	}
	if gotQuery != "bar" {
		t.Errorf("expected query foo=bar, got %q", gotQuery)
	}
	if string(body) != "raw-response" {
		t.Errorf("expected raw body, got %q", string(body))
	}
}

func TestTriggerFlowPost(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"ok":true}`)
	})

	body, err := client.TriggerFlow("POST", "a", map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/flows/trigger/a" {
		t.Errorf("expected path /flows/trigger/a, got %s", gotPath)
	}
	if gotBody["foo"] != "bar" {
		t.Errorf("expected body foo=bar, got %+v", gotBody)
	}
	if string(body) != `{"ok":true}` {
		t.Errorf("unexpected body: %q", string(body))
	}
}

func TestTriggerFlowValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.TriggerFlow("GET", "", nil); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if _, err := client.TriggerFlow("PUT", "a", nil); err == nil {
		t.Error("expected error for unsupported method, got nil")
	}
	if called {
		t.Error("expected no HTTP call")
	}
}

func TestFlowsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetFlows(nil); err == nil {
		t.Error("GetFlows: expected request build error, got nil")
	}
	if _, err := client.TriggerFlow("GET", "a", nil); err == nil {
		t.Error("TriggerFlow: expected request build error, got nil")
	}
}
