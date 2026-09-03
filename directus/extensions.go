package directus

import "fmt"

// ExtensionMeta holds the enable state of an extension.
type ExtensionMeta struct {
	Enabled bool `json:"enabled"`
}

// Extension describes an installed Directus extension.
type Extension struct {
	ID     string         `json:"id"`
	Bundle any            `json:"bundle,omitempty"`
	Schema map[string]any `json:"schema,omitempty"`
	Meta   ExtensionMeta  `json:"meta"`
}

// GetExtensions lists all installed extensions.
func (c *Client) GetExtensions() ([]Extension, error) {
	return request[[]Extension](c, "GET", "/extensions/", nil, nil)
}

// PatchExtension updates a single extension's metadata by name.
func (c *Client) PatchExtension(name string, item map[string]any) (*Extension, error) {
	if name == "" {
		return nil, fmt.Errorf("name must be provided")
	}
	return request[*Extension](c, "PATCH", fmt.Sprintf("/extensions/%s", name), nil, item)
}

// PatchBundleExtension updates a single extension within a bundle.
func (c *Client) PatchBundleExtension(bundle string, name string, item map[string]any) (*Extension, error) {
	if bundle == "" || name == "" {
		return nil, fmt.Errorf("bundle and name must be provided")
	}
	return request[*Extension](c, "PATCH", fmt.Sprintf("/extensions/%s/%s", bundle, name), nil, item)
}

// DeleteExtension deletes an installed extension by id.
func (c *Client) DeleteExtension(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/extensions/%s", id), nil, nil)
}

// InstallRegistryExtension installs an extension from the registry.
func (c *Client) InstallRegistryExtension(extensionID string, version string) error {
	if extensionID == "" || version == "" {
		return fmt.Errorf("extensionID and version must be provided")
	}
	body := map[string]any{"extension": extensionID, "version": version}
	return c.execute("POST", "/extensions/registry/install", nil, body)
}

// UninstallRegistryExtension uninstalls a registry extension by id.
func (c *Client) UninstallRegistryExtension(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/extensions/registry/uninstall/%s", id), nil, nil)
}

// GetRegistryExtensions lists extensions available in the registry.
func (c *Client) GetRegistryExtensions(q *Query) ([]map[string]any, error) {
	return request[[]map[string]any](c, "GET", "/extensions/registry", q, nil)
}
