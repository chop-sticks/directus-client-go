//go:build integration

package integration

import (
	"strings"
	"testing"
)

func TestE2EServerPing(t *testing.T) {
	c := itestClient(t)
	pong, err := c.ServerPing()
	if err != nil {
		t.Fatalf("ServerPing: %v", err)
	}
	if !strings.Contains(pong, "pong") {
		t.Errorf("ServerPing = %q, want it to contain %q", pong, "pong")
	}
}

func TestE2EServerInfo(t *testing.T) {
	c := itestClient(t)
	info, err := c.ServerInfo()
	if err != nil {
		t.Fatalf("ServerInfo: %v", err)
	}
	if len(info) == 0 {
		t.Error("ServerInfo returned empty map")
	}
}

func TestE2EServerHealth(t *testing.T) {
	c := itestClient(t)
	health, err := c.ServerHealth()
	if err != nil {
		t.Fatalf("ServerHealth: %v", err)
	}
	if health == nil || health.Status == "" {
		t.Errorf("ServerHealth status is empty: %+v", health)
	}
}

func TestE2EServerReadOpenAPISpec(t *testing.T) {
	c := itestClient(t)
	spec, err := c.ReadOpenAPISpec()
	if err != nil {
		t.Fatalf("ReadOpenAPISpec: %v", err)
	}
	if len(spec) == 0 {
		t.Error("ReadOpenAPISpec returned empty map")
	}
}

func TestE2EServerReadGraphqlSDL(t *testing.T) {
	c := itestClient(t)

	itemSDL, err := c.ReadGraphqlSDL("item")
	if err != nil {
		t.Fatalf("ReadGraphqlSDL(item): %v", err)
	}
	if strings.TrimSpace(itemSDL) == "" {
		t.Error("ReadGraphqlSDL(item) returned empty string")
	}

	systemSDL, err := c.ReadGraphqlSDL("system")
	if err != nil {
		t.Fatalf("ReadGraphqlSDL(system): %v", err)
	}
	if strings.TrimSpace(systemSDL) == "" {
		t.Error("ReadGraphqlSDL(system) returned empty string")
	}
}

func TestE2EAuthReadProviders(t *testing.T) {
	c := itestClient(t)
	providers, err := c.ReadProviders(false)
	if err != nil {
		t.Fatalf("ReadProviders: %v", err)
	}
	if providers == nil {
		t.Error("ReadProviders returned nil slice")
	}
}

func TestE2EAuthLoginRefreshLogout(t *testing.T) {
	c := itestClient(t)

	auth, err := c.Login(defaultEmail, defaultPassword, "json", "")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if auth == nil || auth.AccessToken == "" || auth.RefreshToken == "" {
		t.Fatalf("Login returned empty tokens: %+v", auth)
	}

	// Refresh rotates the refresh token, so use the token from this login.
	refreshed, err := c.Refresh(auth.RefreshToken, "json")
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if refreshed == nil || refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		t.Fatalf("Refresh returned empty tokens: %+v", refreshed)
	}

	// Do a separate Login for Logout since the earlier refresh token was
	// rotated by Refresh and is no longer valid.
	logoutAuth, err := c.Login(defaultEmail, defaultPassword, "json", "")
	if err != nil {
		t.Fatalf("Login (for logout): %v", err)
	}
	if logoutAuth == nil || logoutAuth.RefreshToken == "" {
		t.Fatalf("Login (for logout) returned empty refresh token: %+v", logoutAuth)
	}
	if err := c.Logout(logoutAuth.RefreshToken, "json"); err != nil {
		t.Fatalf("Logout: %v", err)
	}
}

func TestE2EAuthPasswordRequest(t *testing.T) {
	c := itestClient(t)
	// Environment-dependent (email transport); call but do not fail on error.
	if err := c.PasswordRequest(defaultEmail, ""); err != nil {
		t.Logf("PasswordRequest: %v", err)
	}
}

func TestE2EAuthPasswordReset(t *testing.T) {
	c := itestClient(t)
	if err := c.PasswordReset("badtoken", "Newpassw0rd!"); err == nil {
		t.Error("PasswordReset with bad token: expected error, got nil")
	}
}
