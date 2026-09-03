package directus

// Settings is the Directus project settings singleton.
type Settings struct {
	ID                            int    `json:"id,omitempty"`
	ProjectName                   string `json:"project_name,omitempty"`
	ProjectDescriptor             string `json:"project_descriptor,omitempty"`
	ProjectURL                    string `json:"project_url,omitempty"`
	ReportErrorURL                string `json:"report_error_url,omitempty"`
	ReportBugURL                  string `json:"report_bug_url,omitempty"`
	DefaultLanguage               string `json:"default_language,omitempty"`
	DefaultAppearance             string `json:"default_appearance,omitempty"`
	ProjectColor                  string `json:"project_color,omitempty"`
	ProjectLogo                   any    `json:"project_logo,omitempty"`
	PublicRegistration            bool   `json:"public_registration,omitempty"`
	PublicRegistrationVerifyEmail bool   `json:"public_registration_verify_email,omitempty"`
	StorageAssetTransform         string `json:"storage_asset_transform,omitempty"`
	CustomCSS                     string `json:"custom_css,omitempty"`
	ModuleBar                     any    `json:"module_bar,omitempty"`
	MapboxKey                     string `json:"mapbox_key,omitempty"`
	DefaultThemeLight             string `json:"default_theme_light,omitempty"`
	DefaultThemeDark              string `json:"default_theme_dark,omitempty"`
}

// GetSettings retrieves the project settings singleton.
func (c *Client) GetSettings(q *Query) (*Settings, error) {
	return request[*Settings](c, "GET", "/settings", q, nil)
}

// PatchSettings updates the project settings singleton.
func (c *Client) PatchSettings(item *Settings, q *Query) (*Settings, error) {
	return request[*Settings](c, "PATCH", "/settings", q, item)
}
