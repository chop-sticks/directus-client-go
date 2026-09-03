package directus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		fmt.Fprint(w, `{"data":{"access_token":"at","expires":900000,"refresh_token":"rt"}}`)
	})

	data, err := client.Login("a@b.com", "secret", "", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/auth/login" {
		t.Errorf("expected POST /auth/login, got %s %s", gotMethod, gotPath)
	}
	if gotBody["mode"] != "json" {
		t.Errorf("expected default mode json, got %v", gotBody["mode"])
	}
	if gotBody["otp"] != "123456" {
		t.Errorf("expected otp forwarded, got %v", gotBody["otp"])
	}
	if data == nil || data.AccessToken != "at" || data.Expires != 900000 || data.RefreshToken != "rt" {
		t.Errorf("unexpected decode: %+v", data)
	}
}

func TestLoginOmitsEmptyOTP(t *testing.T) {
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		fmt.Fprint(w, `{"data":{}}`)
	})
	if _, err := client.Login("a@b.com", "secret", "cookie", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotBody["otp"]; ok {
		t.Error("expected otp omitted when empty")
	}
	if gotBody["mode"] != "cookie" {
		t.Errorf("expected mode cookie, got %v", gotBody["mode"])
	}
}

func TestLoginValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.Login("", "secret", "", ""); err == nil {
		t.Error("expected validation error for empty email")
	}
	if _, err := client.Login("a@b.com", "", "", ""); err == nil {
		t.Error("expected validation error for empty password")
	}
	if called {
		t.Error("expected no HTTP call when validation fails")
	}
}

func TestLoginHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	if _, err := client.Login("a@b.com", "secret", "", ""); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestLoginBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{not json}`)
	})
	if _, err := client.Login("a@b.com", "secret", "", ""); err == nil {
		t.Error("expected JSON decode error")
	}
}

func TestRefresh(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		fmt.Fprint(w, `{"data":{"access_token":"at2","expires":1,"refresh_token":"rt2"}}`)
	})

	data, err := client.Refresh("rt", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/auth/refresh" {
		t.Errorf("expected POST /auth/refresh, got %s %s", gotMethod, gotPath)
	}
	if gotBody["mode"] != "json" {
		t.Errorf("expected default mode json, got %v", gotBody["mode"])
	}
	if gotBody["refresh_token"] != "rt" {
		t.Errorf("expected refresh_token forwarded, got %v", gotBody["refresh_token"])
	}
	if data == nil || data.AccessToken != "at2" {
		t.Errorf("unexpected decode: %+v", data)
	}
}

func TestRefreshOmitsEmptyToken(t *testing.T) {
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		fmt.Fprint(w, `{"data":{}}`)
	})
	if _, err := client.Refresh("", "cookie"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotBody["refresh_token"]; ok {
		t.Error("expected refresh_token omitted when empty")
	}
}

func TestRefreshHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	if _, err := client.Refresh("rt", ""); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestRefreshBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad}`)
	})
	if _, err := client.Refresh("rt", ""); err == nil {
		t.Error("expected JSON decode error")
	}
}

func TestLogout(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.Logout("rt", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/auth/logout" {
		t.Errorf("expected POST /auth/logout, got %s %s", gotMethod, gotPath)
	}
	if gotBody["mode"] != "json" {
		t.Errorf("expected default mode json, got %v", gotBody["mode"])
	}
	if gotBody["refresh_token"] != "rt" {
		t.Errorf("expected refresh_token forwarded, got %v", gotBody["refresh_token"])
	}
}

func TestLogoutHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if err := client.Logout("rt", ""); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestPasswordRequest(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.PasswordRequest("a@b.com", "https://reset"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/auth/password/request" {
		t.Errorf("expected POST /auth/password/request, got %s %s", gotMethod, gotPath)
	}
	if gotBody["email"] != "a@b.com" || gotBody["reset_url"] != "https://reset" {
		t.Errorf("unexpected body: %v", gotBody)
	}
}

func TestPasswordRequestOmitsEmptyResetURL(t *testing.T) {
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.PasswordRequest("a@b.com", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotBody["reset_url"]; ok {
		t.Error("expected reset_url omitted when empty")
	}
}

func TestPasswordRequestValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.PasswordRequest("", ""); err == nil {
		t.Error("expected validation error for empty email")
	}
	if called {
		t.Error("expected no HTTP call when validation fails")
	}
}

func TestPasswordReset(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]any
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.PasswordReset("tok", "newpass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/auth/password/reset" {
		t.Errorf("expected POST /auth/password/reset, got %s %s", gotMethod, gotPath)
	}
	if gotBody["token"] != "tok" || gotBody["password"] != "newpass" {
		t.Errorf("unexpected body: %v", gotBody)
	}
}

func TestPasswordResetValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.PasswordReset("", "newpass"); err == nil {
		t.Error("expected validation error for empty token")
	}
	if err := client.PasswordReset("tok", ""); err == nil {
		t.Error("expected validation error for empty password")
	}
	if called {
		t.Error("expected no HTTP call when validation fails")
	}
}

func TestReadProviders(t *testing.T) {
	var gotPath, gotRawQuery, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotRawQuery = r.URL.RawQuery
		gotMethod = r.Method
		fmt.Fprint(w, `{"data":[{"name":"default","driver":"local","label":"Local","icon":null}]}`)
	})

	providers, err := client.ReadProviders(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/auth" {
		t.Errorf("expected GET /auth, got %s %s", gotMethod, gotPath)
	}
	if gotRawQuery != "" {
		t.Errorf("expected no query, got %q", gotRawQuery)
	}
	if len(providers) != 1 || providers[0].Name != "default" || providers[0].Driver != "local" {
		t.Errorf("unexpected decode: %+v", providers)
	}
	if providers[0].Label == nil || *providers[0].Label != "Local" {
		t.Errorf("expected label Local, got %v", providers[0].Label)
	}
	if providers[0].Icon != nil {
		t.Errorf("expected nil icon, got %v", providers[0].Icon)
	}
}

func TestReadProvidersSessionOnly(t *testing.T) {
	var gotRawQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		fmt.Fprint(w, `{"data":[]}`)
	})
	if _, err := client.ReadProviders(true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotRawQuery != "sessionOnly" {
		t.Errorf("expected sessionOnly query, got %q", gotRawQuery)
	}
}

func TestReadProvidersHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.ReadProviders(false); err == nil {
		t.Error("expected HTTP error propagation")
	}
}

func TestReadProvidersBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad}`)
	})
	if _, err := client.ReadProviders(false); err == nil {
		t.Error("expected JSON decode error")
	}
}

func TestAuthRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.Login("a@b.com", "secret", "", ""); err == nil {
		t.Error("Login: expected request build error")
	}
	if _, err := client.ReadProviders(false); err == nil {
		t.Error("ReadProviders: expected request build error")
	}
}
