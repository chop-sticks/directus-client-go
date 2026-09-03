package directus

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestActivityRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetActivities(nil); err == nil {
		t.Error("GetActivities: expected request build error, got nil")
	}
	if _, err := client.GetActivity(1, nil); err == nil {
		t.Error("GetActivity: expected request build error, got nil")
	}
}

func TestGetActivities(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1,"action":"create"},{"id":2,"action":"update"}]}`)
	})

	items, err := client.GetActivities(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/activity" {
		t.Errorf("expected GET /activity, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 2 || items[0].ID != 1 || items[0].Action != "create" {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestGetActivitiesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetActivities(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetActivitiesBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetActivities(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetActivity(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":7,"action":"login","collection":"directus_users"}}`)
	})

	item, err := client.GetActivity(7, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/activity/7" {
		t.Errorf("expected GET /activity/7, got %s %s", gotMethod, gotPath)
	}
	if item == nil || item.ID != 7 || item.Action != "login" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestGetActivityValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	for _, id := range []int{0, -1} {
		if _, err := client.GetActivity(id, nil); err == nil {
			t.Errorf("id=%d: expected validation error, got nil", id)
		} else if !strings.Contains(err.Error(), "id must be provided") {
			t.Errorf("id=%d: unexpected error message: %v", id, err)
		}
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetActivityHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetActivity(999, nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetActivityBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetActivity(1, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}
