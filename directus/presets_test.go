package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetPresets(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1,"collection":"posts"},{"id":2}]}`)
	})

	presets, err := client.GetPresets(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/presets" {
		t.Errorf("expected path /presets, got %s", gotPath)
	}
	if len(presets) != 2 || presets[0].ID != 1 || presets[0].Collection != "posts" {
		t.Errorf("unexpected presets: %+v", presets)
	}
}

func TestGetPresetsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetPresets(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetPresetsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetPresets(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetPreset(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":7,"bookmark":"My view","refresh_interval":30}}`)
	})

	preset, err := client.GetPreset(7, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/presets/7" {
		t.Errorf("expected path /presets/7, got %s", gotPath)
	}
	if preset.ID != 7 || preset.Bookmark != "My view" || preset.RefreshInterval != 30 {
		t.Errorf("unexpected preset: %+v", preset)
	}
}

func TestGetPresetBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetPreset(1, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreatePreset(t *testing.T) {
	var gotMethod, gotPath string
	var gotReq Preset
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotReq)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":3,"collection":"posts"}}`)
	})

	preset, err := client.CreatePreset(&Preset{Collection: "posts"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/presets" {
		t.Errorf("expected path /presets, got %s", gotPath)
	}
	if gotReq.Collection != "posts" {
		t.Errorf("expected request collection 'posts', got %q", gotReq.Collection)
	}
	if preset.ID != 3 {
		t.Errorf("expected id 3, got %d", preset.ID)
	}
}

func TestCreatePresetBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.CreatePreset(&Preset{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreatePresets(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	presets, err := client.CreatePresets([]Preset{{Collection: "a"}, {Collection: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/presets" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(presets) != 2 {
		t.Fatalf("expected 2 presets, got %d", len(presets))
	}
}

func TestPatchPreset(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":5,"color":"#fff"}}`)
	})

	preset, err := client.PatchPreset(5, &Preset{Color: "#fff"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", gotMethod)
	}
	if gotPath != "/presets/5" {
		t.Errorf("expected path /presets/5, got %s", gotPath)
	}
	if preset.ID != 5 || preset.Color != "#fff" {
		t.Errorf("unexpected preset: %+v", preset)
	}
}

func TestPatchPresets(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []int  `json:"keys"`
		Data Preset `json:"data"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	presets, err := client.PatchPresets([]int{1, 2}, &Preset{Color: "#000"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/presets" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Keys[0] != 1 || gotBody.Data.Color != "#000" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(presets) != 2 {
		t.Fatalf("expected 2 presets, got %d", len(presets))
	}
}

func TestPatchPresetsBatch(t *testing.T) {
	var gotMethod, gotPath string
	var gotItems []Preset
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotItems)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	presets, err := client.PatchPresetsBatch([]Preset{{ID: 1}, {ID: 2}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/presets" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotItems) != 2 || gotItems[0].ID != 1 {
		t.Errorf("unexpected items: %+v", gotItems)
	}
	if len(presets) != 2 {
		t.Fatalf("expected 2 presets, got %d", len(presets))
	}
}

func TestDeletePreset(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeletePreset(9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/presets/9" {
		t.Errorf("expected path /presets/9, got %s", gotPath)
	}
}

func TestDeletePresets(t *testing.T) {
	var gotMethod, gotPath string
	var gotKeys []int
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotKeys)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeletePresets([]int{4, 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/presets" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotKeys) != 2 || gotKeys[0] != 4 {
		t.Errorf("expected raw keys array, got %+v", gotKeys)
	}
}

func TestDeletePresetsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.DeletePresets([]int{1}); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestPresetsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetPresets(nil); err == nil {
		t.Error("GetPresets: expected request build error, got nil")
	}
	if _, err := client.GetPreset(1, nil); err == nil {
		t.Error("GetPreset: expected request build error, got nil")
	}
	if _, err := client.CreatePreset(&Preset{}, nil); err == nil {
		t.Error("CreatePreset: expected request build error, got nil")
	}
	if err := client.DeletePreset(1); err == nil {
		t.Error("DeletePreset: expected request build error, got nil")
	}
}
