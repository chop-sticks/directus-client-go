package directus

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Query models the Directus global query parameters shared by most read and
// write endpoints. Only non-zero fields are serialized, so a nil *Query or an
// empty Query yields an empty query string.
//
// Reference: https://docs.directus.io/reference/query.html
type Query struct {
	// Fields to return. Supports dot-notation for relations and "*" wildcards.
	Fields []string
	// Filter is a Directus filter object, e.g. {"status": {"_eq": "published"}}.
	Filter map[string]any
	// Sort fields; prefix a field with "-" for descending order.
	Sort []string
	// Limit caps the number of returned items; -1 returns all.
	Limit *int
	// Offset skips the given number of items.
	Offset *int
	// Page is an alternative to Offset (1-indexed).
	Page *int
	// Search performs a full-text search across the collection.
	Search string
	// Deep applies query parameters to nested relational data.
	Deep map[string]any
	// Aggregate performs aggregation, e.g. {"count": "*"}.
	Aggregate map[string]any
	// GroupBy groups aggregated results by the given fields.
	GroupBy []string
	// Alias renames fields in the response, e.g. {"headline": "title"}.
	Alias map[string]string
	// Version returns the given content version of an item.
	Version string
	// Export sets the export file format (csv, json, xml, yaml).
	Export string
	// Meta requests metadata (e.g. "total_count", "filter_count", "*").
	Meta string
}

// buildQueryString serializes q into a URL query string prefixed with "?".
// A nil or empty query returns "".
func buildQueryString(q *Query) string {
	if q == nil {
		return ""
	}

	values := url.Values{}

	if len(q.Fields) > 0 {
		values.Set("fields", strings.Join(q.Fields, ","))
	}
	if len(q.Sort) > 0 {
		values.Set("sort", strings.Join(q.Sort, ","))
	}
	if len(q.GroupBy) > 0 {
		values.Set("groupBy", strings.Join(q.GroupBy, ","))
	}
	if q.Limit != nil {
		values.Set("limit", strconv.Itoa(*q.Limit))
	}
	if q.Offset != nil {
		values.Set("offset", strconv.Itoa(*q.Offset))
	}
	if q.Page != nil {
		values.Set("page", strconv.Itoa(*q.Page))
	}
	if q.Search != "" {
		values.Set("search", q.Search)
	}
	if q.Version != "" {
		values.Set("version", q.Version)
	}
	if q.Export != "" {
		values.Set("export", q.Export)
	}
	if q.Meta != "" {
		values.Set("meta", q.Meta)
	}
	if len(q.Filter) > 0 {
		if b, err := json.Marshal(q.Filter); err == nil {
			values.Set("filter", string(b))
		}
	}
	if len(q.Deep) > 0 {
		if b, err := json.Marshal(q.Deep); err == nil {
			values.Set("deep", string(b))
		}
	}
	if len(q.Aggregate) > 0 {
		if b, err := json.Marshal(q.Aggregate); err == nil {
			values.Set("aggregate", string(b))
		}
	}
	if len(q.Alias) > 0 {
		if b, err := json.Marshal(q.Alias); err == nil {
			values.Set("alias", string(b))
		}
	}

	if len(values) == 0 {
		return ""
	}
	return "?" + values.Encode()
}

// buildURL joins the base URL, path, and serialized query into a full request URL.
func buildURL(baseURL, path string, q *Query) string {
	return fmt.Sprintf("%s%s%s", baseURL, path, buildQueryString(q))
}
