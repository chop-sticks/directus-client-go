package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetPanels(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a","name":"one"},{"id":"b"}]}`)
	})

	panels, err := client.GetPanels(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/panels" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(panels) != 2 || panels[0].Name != "one" {
		t.Errorf("unexpected panels: %+v", panels)
	}
}

func TestGetPanelsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetPanels(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetPanelsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetPanels(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetPanel(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","show_header":true,"width":4}}`)
	})

	panel, err := client.GetPanel("a", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/panels/a" {
		t.Errorf("expected path /panels/a, got %s", gotPath)
	}
	if !panel.ShowHeader || panel.Width != 4 {
		t.Errorf("unexpected panel: %+v", panel)
	}
}

func TestGetPanelValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetPanel("", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetPanelHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetPanel("missing", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetPanelBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.GetPanel("a", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreatePanel(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Panel
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"one"}}`)
	})

	panel, err := client.CreatePanel(&Panel{Name: "one"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/panels" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotBody.Name != "one" || panel.ID != "a" {
		t.Errorf("unexpected body/panel: %+v %+v", gotBody, panel)
	}
}

func TestCreatePanelBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.CreatePanel(&Panel{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreatePanels(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	panels, err := client.CreatePanels([]Panel{{Name: "one"}, {Name: "two"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/panels" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(panels) != 2 {
		t.Errorf("expected 2 panels, got %d", len(panels))
	}
}

func TestPatchPanel(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"updated"}}`)
	})

	panel, err := client.PatchPanel("a", &Panel{Name: "updated"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/panels/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if panel.Name != "updated" {
		t.Errorf("unexpected panel: %+v", panel)
	}
}

func TestPatchPanelValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchPanel("", &Panel{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchPanels(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []string `json:"keys"`
		Data Panel    `json:"data"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	panels, err := client.PatchPanels([]string{"a", "b"}, &Panel{Color: "red"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/panels" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Data.Color != "red" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(panels) != 2 {
		t.Errorf("expected 2 panels, got %d", len(panels))
	}
}

func TestPatchPanelsBatch(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []Panel
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	panels, err := client.PatchPanelsBatch([]Panel{{ID: "a"}, {ID: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/panels" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0].ID != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(panels) != 2 {
		t.Errorf("expected 2 panels, got %d", len(panels))
	}
}

func TestDeletePanel(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeletePanel("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/panels/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
}

func TestDeletePanelValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeletePanel(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeletePanels(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeletePanels([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/panels" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != "a" {
		t.Errorf("expected raw keys body [a b], got %+v", gotBody)
	}
}

func TestPanelsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetPanels(nil); err == nil {
		t.Error("GetPanels: expected request build error, got nil")
	}
}
