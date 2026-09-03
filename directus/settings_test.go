package directus

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSettingsRequestBuildError(t *testing.T) {
	client := badHostClient(t)

	if _, err := client.GetSettings(nil); err == nil {
		t.Error("GetSettings: expected request build error, got nil")
	}
	if _, err := client.PatchSettings(&Settings{}, nil); err == nil {
		t.Error("PatchSettings: expected request build error, got nil")
	}
}

func TestGetSettings(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":1,"project_name":"Acme","public_registration":true,"default_language":"en-US"}}`)
	})

	settings, err := client.GetSettings(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "GET" || gotPath != "/settings" {
		t.Errorf("expected GET /settings, got %s %s", gotMethod, gotPath)
	}
	if settings.ID != 1 || settings.ProjectName != "Acme" {
		t.Errorf("unexpected settings: %+v", settings)
	}
	if !settings.PublicRegistration || settings.DefaultLanguage != "en-US" {
		t.Errorf("unexpected settings fields: %+v", settings)
	}
}

func TestGetSettingsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetSettings(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetSettingsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetSettings(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestPatchSettings(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":1,"project_name":"Renamed"}}`)
	})

	settings, err := client.PatchSettings(&Settings{ProjectName: "Renamed"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PATCH" || gotPath != "/settings" {
		t.Errorf("expected PATCH /settings, got %s %s", gotMethod, gotPath)
	}
	if settings.ProjectName != "Renamed" {
		t.Errorf("unexpected settings: %+v", settings)
	}
}

func TestPatchSettingsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.PatchSettings(&Settings{}, nil); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestPatchSettingsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.PatchSettings(&Settings{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}
