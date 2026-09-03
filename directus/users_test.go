package directus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetUsers(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"u1","email":"a@x.io"},{"id":"u2","email":"b@x.io"}]}`)
	})
	users, err := client.GetUsers(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users" || gotMethod != "GET" {
		t.Errorf("expected GET /users, got %s %s", gotMethod, gotPath)
	}
	if len(users) != 2 || users[0].ID != "u1" {
		t.Errorf("unexpected users: %+v", users)
	}
}

func TestGetUsersHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	if _, err := client.GetUsers(nil); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestGetUsersBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetUsers(nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetUser(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"u1","first_name":"Ada"}}`)
	})
	user, err := client.GetUser("u1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/u1" || gotMethod != "GET" {
		t.Errorf("expected GET /users/u1, got %s %s", gotMethod, gotPath)
	}
	if user == nil || user.FirstName != "Ada" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestGetUserValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GetUser("", nil); err == nil {
		t.Error("expected validation error, got nil")
	} else if !strings.Contains(err.Error(), "id must be provided") {
		t.Errorf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestGetUserHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if _, err := client.GetUser("u1", nil); err == nil {
		t.Error("expected error for 404 status, got nil")
	}
}

func TestGetUserBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GetUser("u1", nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestGetUsersMe(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		fmt.Fprint(w, `{"data":{"id":"me","email":"me@x.io"}}`)
	})
	user, err := client.GetUsersMe(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/me" {
		t.Errorf("expected /users/me, got %s", gotPath)
	}
	if user == nil || user.ID != "me" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestCreateUser(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"data":{"id":"u1","email":"a@x.io"}}`)
	})
	user, err := client.CreateUser(&User{Email: "a@x.io"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users" || gotMethod != "POST" {
		t.Errorf("expected POST /users, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"email":"a@x.io"`) {
		t.Errorf("unexpected body: %s", gotBody)
	}
	if user == nil || user.ID != "u1" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestCreateUserBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.CreateUser(&User{}, nil); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestCreateUsers(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"u1"},{"id":"u2"}]}`)
	})
	users, err := client.CreateUsers([]User{{Email: "a@x.io"}, {Email: "b@x.io"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users" || gotMethod != "POST" {
		t.Errorf("expected POST /users, got %s %s", gotMethod, gotPath)
	}
	if len(users) != 2 {
		t.Errorf("unexpected users: %+v", users)
	}
}

func TestPatchUser(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"u1","first_name":"Ada"}}`)
	})
	user, err := client.PatchUser("u1", &User{FirstName: "Ada"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/u1" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /users/u1, got %s %s", gotMethod, gotPath)
	}
	if user == nil || user.FirstName != "Ada" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestPatchUserValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.PatchUser("", &User{}, nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestPatchUsers(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		fmt.Fprint(w, `{"data":[{"id":"u1"},{"id":"u2"}]}`)
	})
	users, err := client.PatchUsers([]string{"u1", "u2"}, &User{Status: "active"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /users, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"keys"`) || !strings.Contains(gotBody, `"data"`) {
		t.Errorf("expected keys+data body, got %s", gotBody)
	}
	if len(users) != 2 {
		t.Errorf("unexpected users: %+v", users)
	}
}

func TestPatchUsersBatch(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":[{"id":"u1"}]}`)
	})
	users, err := client.PatchUsersBatch([]User{{ID: "u1", Status: "active"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /users, got %s %s", gotMethod, gotPath)
	}
	if len(users) != 1 {
		t.Errorf("unexpected users: %+v", users)
	}
}

func TestPatchUsersMe(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"id":"me","first_name":"Ada"}}`)
	})
	user, err := client.PatchUsersMe(&User{FirstName: "Ada"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/me" || gotMethod != "PATCH" {
		t.Errorf("expected PATCH /users/me, got %s %s", gotMethod, gotPath)
	}
	if user == nil || user.FirstName != "Ada" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestDeleteUser(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteUser("u1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/u1" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /users/u1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteUserValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DeleteUser(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty id")
	}
}

func TestDeleteUsers(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DeleteUsers([]string{"u1", "u2"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users" || gotMethod != "DELETE" {
		t.Errorf("expected DELETE /users, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `["u1","u2"]`) {
		t.Errorf("expected raw keys body, got %s", gotBody)
	}
}

func TestInviteUser(t *testing.T) {
	var gotPath, gotMethod, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.InviteUser("a@x.io", "r1", "https://app/invite"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/invite" || gotMethod != "POST" {
		t.Errorf("expected POST /users/invite, got %s %s", gotMethod, gotPath)
	}
	if !strings.Contains(gotBody, `"email":"a@x.io"`) || !strings.Contains(gotBody, `"invite_url":"https://app/invite"`) {
		t.Errorf("unexpected body: %s", gotBody)
	}
}

func TestInviteUserOmitsURL(t *testing.T) {
	var gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.InviteUser("a@x.io", "r1", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(gotBody, "invite_url") {
		t.Errorf("expected no invite_url, got %s", gotBody)
	}
}

func TestInviteUserValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.InviteUser("", "r1", ""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for missing email")
	}
}

func TestAcceptUserInvite(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.AcceptUserInvite("tok", "pw"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/invite/accept" {
		t.Errorf("expected /users/invite/accept, got %s", gotPath)
	}
}

func TestAcceptUserInviteValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.AcceptUserInvite("", "pw"); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for missing token")
	}
}

func TestRegisterUser(t *testing.T) {
	var gotPath, gotBody string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.RegisterUser("a@x.io", "pw", map[string]any{"first_name": "Ada"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/register" {
		t.Errorf("expected /users/register, got %s", gotPath)
	}
	if !strings.Contains(gotBody, `"email":"a@x.io"`) || !strings.Contains(gotBody, `"first_name":"Ada"`) {
		t.Errorf("expected merged body, got %s", gotBody)
	}
}

func TestRegisterUserValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.RegisterUser("a@x.io", "", nil); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for missing password")
	}
}

func TestRegisterUserVerify(t *testing.T) {
	var gotPath, gotQuery string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.RawQuery
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.RegisterUserVerify("tok123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/register/verify-email" {
		t.Errorf("expected /users/register/verify-email, got %s", gotPath)
	}
	if gotQuery != "token=tok123" {
		t.Errorf("expected token query, got %s", gotQuery)
	}
}

func TestRegisterUserVerifyValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.RegisterUserVerify(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty token")
	}
}

func TestGenerateTwoFactorSecret(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		fmt.Fprint(w, `{"data":{"secret":"S","otpauth_url":"otpauth://x"}}`)
	})
	res, err := client.GenerateTwoFactorSecret("pw")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/me/tfa/generate" || gotMethod != "POST" {
		t.Errorf("expected POST /users/me/tfa/generate, got %s %s", gotMethod, gotPath)
	}
	if res["secret"] != "S" {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestGenerateTwoFactorSecretValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if _, err := client.GenerateTwoFactorSecret(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty password")
	}
}

func TestGenerateTwoFactorSecretBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{bad`)
	})
	if _, err := client.GenerateTwoFactorSecret("pw"); err == nil {
		t.Error("expected JSON error, got nil")
	}
}

func TestEnableTwoFactor(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.EnableTwoFactor("S", "123456"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/me/tfa/enable" {
		t.Errorf("expected /users/me/tfa/enable, got %s", gotPath)
	}
}

func TestEnableTwoFactorValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.EnableTwoFactor("S", ""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for missing otp")
	}
}

func TestDisableTwoFactor(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	if err := client.DisableTwoFactor("123456"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/users/me/tfa/disable" {
		t.Errorf("expected /users/me/tfa/disable, got %s", gotPath)
	}
}

func TestDisableTwoFactorValidation(t *testing.T) {
	called := false
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	if err := client.DisableTwoFactor(""); err == nil {
		t.Error("expected validation error, got nil")
	}
	if called {
		t.Error("expected no HTTP call for empty otp")
	}
}

func TestUsersRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.GetUsers(nil); err == nil {
		t.Error("GetUsers: expected request build error, got nil")
	}
}
