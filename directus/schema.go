package directus

import (
	"fmt"
	"net/url"
	"strings"
)

// SchemaSnapshot is a point-in-time snapshot of the instance's data model.
type SchemaSnapshot struct {
	Version      int              `json:"version"`
	Directus     string           `json:"directus"`
	Vendor       string           `json:"vendor"`
	Collections  []map[string]any `json:"collections"`
	Fields       []map[string]any `json:"fields"`
	SystemFields []map[string]any `json:"systemFields"`
	Relations    []map[string]any `json:"relations"`
}

// SchemaDiff describes the difference between a snapshot and the current schema.
type SchemaDiff struct {
	Hash string         `json:"hash"`
	Diff map[string]any `json:"diff"`
}

// GetSchemaSnapshot retrieves the current schema snapshot. Only one of
// includeCollections or excludeCollections is applied (includeCollections takes
// precedence); the collection names are comma-joined into a query parameter.
func (c *Client) GetSchemaSnapshot(includeCollections, excludeCollections []string) (*SchemaSnapshot, error) {
	path := "/schema/snapshot"
	params := url.Values{}
	if len(includeCollections) > 0 {
		params.Set("includeCollections", strings.Join(includeCollections, ","))
	} else if len(excludeCollections) > 0 {
		params.Set("excludeCollections", strings.Join(excludeCollections, ","))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	return request[*SchemaSnapshot](c, "GET", path, nil, nil)
}

// SchemaDiffSnapshot compares the given snapshot against the current schema and
// returns the diff to apply. When there is no difference the server responds
// with an empty body (204), in which case nil is returned. The force and mode
// parameters are sent as query parameters when set.
func (c *Client) SchemaDiffSnapshot(snapshot *SchemaSnapshot, force bool, mode string) (*SchemaDiff, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot must be provided")
	}
	path := "/schema/diff"
	params := url.Values{}
	if force {
		params.Set("force", "true")
	}
	if mode != "" {
		params.Set("mode", mode)
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	return request[*SchemaDiff](c, "POST", path, nil, snapshot)
}

// SchemaApply applies the given diff to the instance's schema. Set force to
// bypass version/vendor checks.
func (c *Client) SchemaApply(diff *SchemaDiff, force bool) error {
	if diff == nil {
		return fmt.Errorf("diff must be provided")
	}
	path := "/schema/apply"
	if force {
		path += "?force=true"
	}
	return c.execute("POST", path, nil, diff)
}
