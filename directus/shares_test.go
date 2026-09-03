package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestGetShares(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"s1","name":"Link","times_used":3}]}`)
	})

	shares, err := client.GetShares(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/shares" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(shares) != 1 || shares[0].Name != "Link" || shares[0].TimesUsed != 3 {
		t.Errorf("unexpected shares: %+v", shares)
	}
}

func TestGetSharesHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetShares(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetSharesBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetShares(nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestGetShare(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"abc","collection":"posts"}}`)
	})

	share, err := client.GetShare("abc", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/shares/abc" {
		t.Errorf("expected path /shares/abc, got %s", gotPath)
	}
	if share.ID != "abc" || share.Collection != "posts" {
		t.Errorf("unexpected share: %+v", share)
	}
}

func TestGetShareValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetShare("", nil); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetShareBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.GetShare("abc", nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateShare(t *testing.T) {
	var gotMethod, gotPath string
	var gotReq Share
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotReq)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"new","name":"L"}}`)
	})

	share, err := client.CreateShare(&Share{Name: "L", Collection: "posts"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/shares" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotReq.Name != "L" {
		t.Errorf("unexpected request: %+v", gotReq)
	}
	if share.ID != "new" {
		t.Errorf("expected id new, got %q", share.ID)
	}
}

func TestCreateShareBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.CreateShare(&Share{}, nil); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestCreateShares(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	shares, err := client.CreateShares([]Share{{Name: "a"}, {Name: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/shares" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(shares) != 2 {
		t.Fatalf("expected 2 shares, got %d", len(shares))
	}
}

func TestPatchShare(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"z","name":"renamed"}}`)
	})

	share, err := client.PatchShare("z", &Share{Name: "renamed"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/shares/z" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if share.Name != "renamed" {
		t.Errorf("unexpected share: %+v", share)
	}
}

func TestPatchShareValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchShare("", &Share{}, nil); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchShares(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Keys []string `json:"keys"`
		Data Share    `json:"data"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	shares, err := client.PatchShares([]string{"a", "b"}, &Share{Name: "n"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/shares" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotBody.Keys) != 2 || gotBody.Data.Name != "n" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if len(shares) != 2 {
		t.Fatalf("expected 2 shares, got %d", len(shares))
	}
}

func TestPatchSharesBatch(t *testing.T) {
	var gotMethod, gotPath string
	var gotItems []Share
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotItems)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}]}`)
	})

	shares, err := client.PatchSharesBatch([]Share{{ID: "a"}, {ID: "b"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/shares" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotItems) != 2 || gotItems[0].ID != "a" {
		t.Errorf("unexpected items: %+v", gotItems)
	}
	if len(shares) != 2 {
		t.Fatalf("expected 2 shares, got %d", len(shares))
	}
}

func TestDeleteShare(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteShare("q"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/shares/q" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
}

func TestDeleteShareValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteShare(""); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteShares(t *testing.T) {
	var gotMethod, gotPath string
	var gotKeys []string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotKeys)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteShares([]string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/shares" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if len(gotKeys) != 2 || gotKeys[0] != "a" {
		t.Errorf("expected raw keys array, got %+v", gotKeys)
	}
}

func TestAuthenticateShare(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Share    string `json:"share"`
		Password string `json:"password"`
		Mode     string `json:"mode"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"access_token":"tok","expires":900,"refresh_token":"ref"}}`)
	})

	auth, err := client.AuthenticateShare("s1", "pw", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/shares/auth" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotBody.Share != "s1" || gotBody.Password != "pw" || gotBody.Mode != "json" {
		t.Errorf("unexpected body: %+v", gotBody)
	}
	if auth.AccessToken != "tok" || auth.RefreshToken != "ref" {
		t.Errorf("unexpected auth: %+v", auth)
	}
}

func TestAuthenticateShareValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.AuthenticateShare("", "pw", "json"); err == nil {
		t.Error("expected validation error for empty share, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty share")
	}
}

func TestInviteShare(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody struct {
		Share  string   `json:"share"`
		Emails []string `json:"emails"`
	}
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.InviteShare("s1", []string{"a@x.com", "b@x.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/shares/invite" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if gotBody.Share != "s1" || len(gotBody.Emails) != 2 {
		t.Errorf("unexpected body: %+v", gotBody)
	}
}

func TestInviteShareValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.InviteShare("", []string{"a@x.com"}); err == nil {
		t.Error("expected validation error for empty share, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty share")
	}
}

func TestReadShareInfo(t *testing.T) {
	var gotMethod, gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"id":"s1","collection":"posts","item":"5","password":null,"times_used":2,"max_uses":null}}`)
	})

	info, err := client.ReadShareInfo("s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/shares/info/s1" {
		t.Errorf("unexpected method/path: %s %s", gotMethod, gotPath)
	}
	if info.ID != "s1" || info.Collection != "posts" || info.Item != "5" {
		t.Errorf("unexpected info: %+v", info)
	}
	if info.Password != nil {
		t.Errorf("expected nil password, got %v", info.Password)
	}
	if info.TimesUsed == nil || *info.TimesUsed != 2 {
		t.Errorf("expected times_used 2, got %v", info.TimesUsed)
	}
	if info.MaxUses != nil {
		t.Errorf("expected nil max_uses, got %v", info.MaxUses)
	}
}

func TestReadShareInfoValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.ReadShareInfo(""); err == nil {
		t.Error("expected validation error for empty id, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestReadShareInfoBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `not json`)
	})
	if _, err := client.ReadShareInfo("s1"); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}

func TestSharesRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetShares(nil); err == nil {
		t.Error("GetShares: expected request build error, got nil")
	}
	if _, err := client.CreateShare(&Share{}, nil); err == nil {
		t.Error("CreateShare: expected request build error, got nil")
	}
	if _, err := client.AuthenticateShare("s1", "pw", "json"); err == nil {
		t.Error("AuthenticateShare: expected request build error, got nil")
	}
	if err := client.InviteShare("s1", nil); err == nil {
		t.Error("InviteShare: expected request build error, got nil")
	}
}
