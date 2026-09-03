package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetDashboards(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a","name":"one"},{"id":"b"}]}`)
	})

	dashboards, err := client.GetDashboards(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/dashboards" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(dashboards) != 2 || dashboards[0].Name != "one" {
		t.Errorf("unexpected dashboards: %+v", dashboards)
	}
}

func TestGetDashboardsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetDashboards(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetDashboardsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetDashboards(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetDashboard(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"one","color":"#fff"}}`)
	})

	dashboard, err := client.GetDashboard("a", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/dashboards/a" {
		t.Errorf("expected path /dashboards/a, got %s", gotPath)
	}
	if dashboard.Color != "#fff" {
		t.Errorf("unexpected dashboard: %+v", dashboard)
	}
}

func TestGetDashboardValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetDashboard("", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetDashboardHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetDashboard("missing", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetDashboardBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.GetDashboard("a", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateDashboard(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Dashboard
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"one"}}`)
	})

	dashboard, err := client.CreateDashboard(&Dashboard{Name: "one"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/dashboards" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotBody.Name != "one" || dashboard.ID != "a" {
		t.Errorf("unexpected body/dashboard: %+v %+v", gotBody, dashboard)
	}
}

func TestCreateDashboardBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `nope`)
	})
	if _, err := client.CreateDashboard(&Dashboard{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateDashboards(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	dashboards, err := client.CreateDashboards([]Dashboard{{Name: "one"}, {Name: "two"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/dashboards" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(dashboards) != 2 {
		t.Errorf("expected 2 dashboards, got %d", len(dashboards))
	}
}

func TestPatchDashboard(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"a","name":"updated"}}`)
	})

	dashboard, err := client.PatchDashboard("a", &Dashboard{Name: "updated"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/dashboards/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if dashboard.Name != "updated" {
		t.Errorf("unexpected dashboard: %+v", dashboard)
	}
}

func TestPatchDashboardValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchDashboard("", &Dashboard{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchDashboards(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []string  `json:"keys"`
		Data Dashboard `json:"data"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	dashboards, err := client.PatchDashboards([]string{"a", "b"}, &Dashboard{Color: "blue"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/dashboards" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Data.Color != "blue" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(dashboards) != 2 {
		t.Errorf("expected 2 dashboards, got %d", len(dashboards))
	}
}

func TestPatchDashboardsBatch(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []Dashboard
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	dashboards, err := client.PatchDashboardsBatch([]Dashboard{{ID: "a"}, {ID: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/dashboards" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0].ID != "a" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(dashboards) != 2 {
		t.Errorf("expected 2 dashboards, got %d", len(dashboards))
	}
}

func TestDeleteDashboard(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteDashboard("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/dashboards/a" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
}

func TestDeleteDashboardValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteDashboard(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteDashboards(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteDashboards([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/dashboards" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != "a" {
		t.Errorf("expected raw keys body [a b], got %+v", gotBody)
	}
}

func TestDashboardsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetDashboards(nil); err == nil {
		t.Error("GetDashboards: expected request build error, got nil")
	}
}
