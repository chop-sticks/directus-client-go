package directus

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// errorRoundTripper is a mock RoundTripper that returns an error
type errorRoundTripper struct{}

func (m *errorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("network error")
}

// readErrorBody is an io.ReadCloser that returns an error on Read
type readErrorBody struct{}

func (m *readErrorBody) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (m *readErrorBody) Close() error {
	return nil
}

// RoundTripFunc is a helper for mocking http.RoundTripper
type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func setupMockServer(t *testing.T) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default handler, can be overridden in tests
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{}")
	}))

	token := "test_token"
	host := "https://example.com"
	client, _ := NewClient(&host, &token)
	client.HostURL = server.URL

	return server, client
}

func TestNewClient(t *testing.T) {
	token := "test_token"
	host := "https://example.com"
	client, err := NewClient(&host, &token)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if client.Token != token {
		t.Errorf("expected token %s, got %s", token, client.Token)
	}
	if client.HostURL != host {
		t.Errorf("expected HostURL %s, got %s", host, client.HostURL)
	}
}

func TestDoRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"bad request"}`)
	}))
	defer server.Close()

	token := "test_token"
	host := "https://example.com"
	client, _ := NewClient(&host, &token)
	client.HostURL = server.URL

	req, _ := http.NewRequest("GET", client.HostURL, nil)
	_, err := client.doRequest(req)

	if err == nil {
		t.Error("expected error for 400 status, got nil")
	}
}

func TestDoRequestNetworkError(t *testing.T) {
	token := "test_token"
	host := "https://example.com"
	client, _ := NewClient(&host, &token)
	client.HTTPClient.Transport = &errorRoundTripper{}

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, err := client.doRequest(req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Errorf("expected error to contain 'network error', got %v", err)
	}
}

func TestDoRequestReadError(t *testing.T) {
	token := "test_token"
	host := "https://example.com"
	client, _ := NewClient(&host, &token)

	// Create a transport that returns a response with a body that fails on Read
	client.HTTPClient.Transport = RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &readErrorBody{},
			Header:     make(http.Header),
		}, nil
	})

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, err := client.doRequest(req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "read error") {
		t.Errorf("expected error to contain 'read error', got %v", err)
	}
}

// newMockClient starts an httptest server with the given handler and returns a
// Client pointed at it. The server is closed automatically at test end.
func newMockClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	token := "test_token"
	host := server.URL
	client, _ := NewClient(&host, &token)
	return client
}

func TestDoRequestSuccess(t *testing.T) {
	var gotAuth, gotContentType string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"ok":true}`)
	})

	req, _ := http.NewRequest(http.MethodGet, client.HostURL, nil)
	body, err := client.doRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Errorf("unexpected body: %s", body)
	}
	if gotAuth != "Bearer test_token" {
		t.Errorf("expected Authorization 'Bearer test_token', got %q", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected default Content-Type application/json, got %q", gotContentType)
	}
}

func TestDoRequestPreservesContentType(t *testing.T) {
	var gotContentType string
	client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodPost, client.HostURL, nil)
	req.Header.Set("Content-Type", "text/plain")
	if _, err := client.doRequest(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotContentType != "text/plain" {
		t.Errorf("expected preserved Content-Type text/plain, got %q", gotContentType)
	}
}

func TestDoRequestAcceptsSuccessStatuses(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			client := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			})
			req, _ := http.NewRequest(http.MethodGet, client.HostURL, nil)
			if _, err := client.doRequest(req); err != nil {
				t.Errorf("status %d: unexpected error: %v", status, err)
			}
		})
	}
}
