package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNotificationsRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetNotifications(nil); err == nil {
		t.Error("GetNotifications: expected request build error, got nil")
	}
	if _, err := client.GetNotification(1, nil); err == nil {
		t.Error("GetNotification: expected request build error, got nil")
	}
	if _, err := client.CreateNotification(&Notification{}, nil); err == nil {
		t.Error("CreateNotification: expected request build error, got nil")
	}
	if _, err := client.CreateNotifications([]Notification{{}}, nil); err == nil {
		t.Error("CreateNotifications: expected request build error, got nil")
	}
	if _, err := client.PatchNotification(1, &Notification{}, nil); err == nil {
		t.Error("PatchNotification: expected request build error, got nil")
	}
	if _, err := client.PatchNotifications([]int{1}, &Notification{}, nil); err == nil {
		t.Error("PatchNotifications: expected request build error, got nil")
	}
	if _, err := client.PatchNotificationsBatch([]Notification{{}}, nil); err == nil {
		t.Error("PatchNotificationsBatch: expected request build error, got nil")
	}
	if err := client.DeleteNotification(1); err == nil {
		t.Error("DeleteNotification: expected request build error, got nil")
	}
	if err := client.DeleteNotifications([]int{1}); err == nil {
		t.Error("DeleteNotifications: expected request build error, got nil")
	}
}

func TestGetNotifications(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1,"subject":"Hi"},{"id":2,"subject":"Yo"}]}`)
	})

	items, err := client.GetNotifications(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/notifications" {
		t.Errorf("expected GET /notifications, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 2 || items[0].ID != 1 || items[0].Subject != "Hi" {
		t.Errorf("unexpected items: %+v", items)
	}
}

func TestGetNotificationsHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetNotifications(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetNotificationsBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetNotifications(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetNotification(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":1,"status":"inbox","subject":"Hi"}}`)
	})

	item, err := client.GetNotification(1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/notifications/1" {
		t.Errorf("expected GET /notifications/1, got %s %s", gotMethod, gotPath)
	}
	if item == nil || item.ID != 1 || item.Status != "inbox" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestGetNotificationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.GetNotification(0, nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error message: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestGetNotificationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetNotification(999, nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetNotificationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetNotification(1, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateNotification(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody Notification
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":1,"subject":"Hi"}}`)
	})

	item, err := client.CreateNotification(&Notification{Subject: "Hi"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/notifications" {
		t.Errorf("expected POST /notifications, got %s %s", gotMethod, gotPath)
	}
	if gotBody.Subject != "Hi" {
		t.Errorf("unexpected request body: %+v", gotBody)
	}
	if item == nil || item.ID != 1 {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestCreateNotificationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	if _, err := client.CreateNotification(&Notification{}, nil); err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestCreateNotificationBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateNotification(&Notification{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateNotifications(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	items, err := client.CreateNotifications([]Notification{{Subject: "a"}, {Subject: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/notifications" {
		t.Errorf("expected POST /notifications, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestPatchNotification(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":1,"status":"archived"}}`)
	})

	item, err := client.PatchNotification(1, &Notification{Status: "archived"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/notifications/1" {
		t.Errorf("expected PATCH /notifications/1, got %s %s", gotMethod, gotPath)
	}
	if item == nil || item.Status != "archived" {
		t.Errorf("unexpected item: %+v", item)
	}
}

func TestPatchNotificationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	if _, err := client.PatchNotification(0, &Notification{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestPatchNotifications(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1},{"id":2}]}`)
	})

	items, err := client.PatchNotifications([]int{1, 2}, &Notification{Status: "archived"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/notifications" {
		t.Errorf("expected PATCH /notifications, got %s %s", gotMethod, gotPath)
	}
	if _, ok := gotBody["keys"]; !ok {
		t.Errorf("expected keys in body, got %+v", gotBody)
	}
	if _, ok := gotBody["data"]; !ok {
		t.Errorf("expected data in body, got %+v", gotBody)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestPatchNotificationsBatch(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":1}]}`)
	})

	items, err := client.PatchNotificationsBatch([]Notification{{ID: 1, Status: "archived"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/notifications" {
		t.Errorf("expected PATCH /notifications, got %s %s", gotMethod, gotPath)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestDeleteNotification(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteNotification(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/notifications/1" {
		t.Errorf("expected DELETE /notifications/1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteNotificationValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteNotification(0); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for invalid arguments")
	}
}

func TestDeleteNotificationHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if err := client.DeleteNotification(999); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestDeleteNotifications(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody []int
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteNotifications([]int{1, 2}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/notifications" {
		t.Errorf("expected DELETE /notifications, got %s %s", gotMethod, gotPath)
	}
	if len(gotBody) != 2 || gotBody[0] != 1 {
		t.Errorf("expected raw keys array body, got %+v", gotBody)
	}
}
