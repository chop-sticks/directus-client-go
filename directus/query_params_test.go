package directus

import "testing"

func TestBuildQueryString(t *testing.T) {
	// buildQueryString currently ignores its argument and always returns "".
	cases := []struct {
		name   string
		params interface{}
	}{
		{"untyped nil", nil},
		{"arbitrary struct", struct{ Limit int }{Limit: 10}},
		{"string", "anything"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildQueryString(tc.params); got != "" {
				t.Errorf("expected empty query string, got %q", got)
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	cases := []struct {
		name    string
		baseURL string
		path    string
		params  interface{}
		want    string
	}{
		{"no params", "https://example.com", "/collections", nil, "https://example.com/collections"},
		{"params ignored", "https://example.com", "/items/articles", struct{ X int }{X: 1}, "https://example.com/items/articles"},
		{"empty path", "https://example.com", "", nil, "https://example.com"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildURL(tc.baseURL, tc.path, tc.params); got != tc.want {
				t.Errorf("buildURL(%q, %q, %v) = %q, want %q", tc.baseURL, tc.path, tc.params, got, tc.want)
			}
		})
	}
}
