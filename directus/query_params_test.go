package directus

import (
	"strings"
	"testing"
)

func intPtr(i int) *int { return &i }

func TestBuildQueryString(t *testing.T) {
	cases := []struct {
		name     string
		query    *Query
		contains []string
		want     string // exact match when set
	}{
		{"nil", nil, nil, ""},
		{"empty", &Query{}, nil, ""},
		{"limit", &Query{Limit: intPtr(10)}, nil, "?limit=10"},
		{"offset+page", &Query{Offset: intPtr(5), Page: intPtr(2)}, []string{"offset=5", "page=2"}, ""},
		{"fields joined", &Query{Fields: []string{"id", "title"}}, []string{"fields=id%2Ctitle"}, ""},
		{"sort joined", &Query{Sort: []string{"-date", "id"}}, []string{"sort=-date%2Cid"}, ""},
		{"search", &Query{Search: "hello world"}, []string{"search=hello+world"}, ""},
		{"filter json", &Query{Filter: map[string]any{"status": map[string]any{"_eq": "published"}}}, []string{"filter="}, ""},
		{"aggregate json", &Query{Aggregate: map[string]any{"count": "*"}}, []string{"aggregate="}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildQueryString(tc.query)
			if tc.want != "" || (tc.query == nil || len(tc.contains) == 0) {
				if got != tc.want {
					t.Errorf("buildQueryString() = %q, want %q", got, tc.want)
				}
			}
			for _, sub := range tc.contains {
				if !strings.Contains(got, sub) {
					t.Errorf("buildQueryString() = %q, want it to contain %q", got, sub)
				}
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	base := "https://example.com"
	cases := []struct {
		name  string
		path  string
		query *Query
		want  string
	}{
		{"no query", "/collections", nil, "https://example.com/collections"},
		{"empty query", "/items/articles", &Query{}, "https://example.com/items/articles"},
		{"with limit", "/items/articles", &Query{Limit: intPtr(1)}, "https://example.com/items/articles?limit=1"},
		{"empty path", "", nil, "https://example.com"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildURL(base, tc.path, tc.query); got != tc.want {
				t.Errorf("buildURL(%q, %q, %v) = %q, want %q", base, tc.path, tc.query, got, tc.want)
			}
		})
	}
}
