package directus

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServerRequestBuildError(t *testing.T) {
	client := badHostClient(t)
	if _, err := client.ServerHealth(); err == nil {
		t.Error("ServerHealth: expected request build error, got nil")
	}
	if _, err := client.ServerPing(); err == nil {
		t.Error("ServerPing: expected request build error, got nil")
	}
	if _, err := client.ServerInfo(); err == nil {
		t.Error("ServerInfo: expected request build error, got nil")
	}
	if _, err := client.ReadOpenAPISpec(); err == nil {
		t.Error("ReadOpenAPISpec: expected request build error, got nil")
	}
	if _, err := client.ReadGraphqlSDL(""); err == nil {
		t.Error("ReadGraphqlSDL: expected request build error, got nil")
	}
}

func TestServerHealth(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		// Not enveloped.
		fmt.Fprint(w, `{"status":"ok","releaseId":"11.0.0","serviceId":"abc","checks":{"pg:responseTime":[{"status":"ok"}]}}`)
	})

	health, err := client.ServerHealth()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/server/health" {
		t.Errorf("expected path /server/health, got %s", gotPath)
	}
	if health == nil {
		t.Fatal("expected health, got nil")
	}
	if health.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", health.Status)
	}
	if health.ReleaseID != "11.0.0" {
		t.Errorf("expected releaseId '11.0.0', got %q", health.ReleaseID)
	}
	if health.ServiceID != "abc" {
		t.Errorf("expected serviceId 'abc', got %q", health.ServiceID)
	}
	if health.Checks == nil {
		t.Error("expected checks, got nil")
	}
}

func TestServerHealthHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.ServerHealth(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestServerHealthBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{not json`)
	})
	if _, err := client.ServerHealth(); err == nil {
		t.Error("expected JSON decode error, got nil")
	}
}

func TestServerPing(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	})

	pong, err := client.ServerPing()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/server/ping" {
		t.Errorf("expected path /server/ping, got %s", gotPath)
	}
	if pong != "pong" {
		t.Errorf("expected 'pong', got %q", pong)
	}
}

func TestServerPingHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.ServerPing(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestServerInfo(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		// Enveloped.
		fmt.Fprint(w, `{"data":{"project":{"project_name":"Directus"},"version":"11.0.0"}}`)
	})

	info, err := client.ServerInfo()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/server/info" {
		t.Errorf("expected path /server/info, got %s", gotPath)
	}
	if info["version"] != "11.0.0" {
		t.Errorf("expected version '11.0.0', got %v", info["version"])
	}
}

func TestServerInfoHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.ServerInfo(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestServerInfoBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{not json`)
	})
	if _, err := client.ServerInfo(); err == nil {
		t.Error("expected JSON decode error, got nil")
	}
}

func TestReadOpenAPISpec(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		// Not enveloped.
		fmt.Fprint(w, `{"openapi":"3.0.1","info":{"title":"Directus API"}}`)
	})

	spec, err := client.ReadOpenAPISpec()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/server/specs/oas" {
		t.Errorf("expected path /server/specs/oas, got %s", gotPath)
	}
	if spec["openapi"] != "3.0.1" {
		t.Errorf("expected openapi '3.0.1', got %v", spec["openapi"])
	}
}

func TestReadOpenAPISpecHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.ReadOpenAPISpec(); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}

func TestReadOpenAPISpecBadJSON(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{not json`)
	})
	if _, err := client.ReadOpenAPISpec(); err == nil {
		t.Error("expected JSON decode error, got nil")
	}
}

func TestReadGraphqlSDL(t *testing.T) {
	var gotPath, gotMethod string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "type Query { foo: String }")
	})

	sdl, err := client.ReadGraphqlSDL("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/server/specs/graphql" {
		t.Errorf("expected path /server/specs/graphql, got %s", gotPath)
	}
	if sdl != "type Query { foo: String }" {
		t.Errorf("unexpected SDL: %q", sdl)
	}
}

func TestReadGraphqlSDLSystem(t *testing.T) {
	var gotPath string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "type Query { sys: String }")
	})

	if _, err := client.ReadGraphqlSDL("system"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/server/specs/graphql/system" {
		t.Errorf("expected path /server/specs/graphql/system, got %s", gotPath)
	}
}

func TestReadGraphqlSDLHTTPError(t *testing.T) {
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"errors":[]}`)
	})
	if _, err := client.ReadGraphqlSDL(""); err == nil {
		t.Error("expected error for 500 status, got nil")
	}
}
