package directus

import "fmt"

// User is a Directus user (system collection `directus_users`).
type User struct {
	ID                  string         `json:"id,omitempty"`
	Status              string         `json:"status,omitempty"`
	FirstName           string         `json:"first_name,omitempty"`
	LastName            string         `json:"last_name,omitempty"`
	Email               string         `json:"email,omitempty"`
	Password            string         `json:"password,omitempty"`
	Token               string         `json:"token,omitempty"`
	LastAccess          string         `json:"last_access,omitempty"`
	LastPage            string         `json:"last_page,omitempty"`
	ExternalIdentifier  string         `json:"external_identifier,omitempty"`
	TFASecret           string         `json:"tfa_secret,omitempty"`
	AuthData            map[string]any `json:"auth_data,omitempty"`
	Provider            string         `json:"provider,omitempty"`
	Appearance          string         `json:"appearance,omitempty"`
	ThemeLight          string         `json:"theme_light,omitempty"`
	ThemeDark           string         `json:"theme_dark,omitempty"`
	ThemeLightOverrides map[string]any `json:"theme_light_overrides,omitempty"`
	ThemeDarkOverrides  map[string]any `json:"theme_dark_overrides,omitempty"`
	Role                any            `json:"role,omitempty"`
	Policies            []any          `json:"policies,omitempty"`
	Language            string         `json:"language,omitempty"`
	TextDirection       string         `json:"text_direction,omitempty"`
	Avatar              any            `json:"avatar,omitempty"`
	Title               string         `json:"title,omitempty"`
	Description         string         `json:"description,omitempty"`
	Location            string         `json:"location,omitempty"`
	Tags                []string       `json:"tags,omitempty"`
	EmailNotifications  bool           `json:"email_notifications,omitempty"`
}

// GetUsers lists users.
func (c *Client) GetUsers(q *Query) ([]User, error) {
	return request[[]User](c, "GET", "/users", q, nil)
}

// GetUser retrieves a single user by ID.
func (c *Client) GetUser(id string, q *Query) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*User](c, "GET", fmt.Sprintf("/users/%s", id), q, nil)
}

// GetUsersMe retrieves the currently authenticated user.
func (c *Client) GetUsersMe(q *Query) (*User, error) {
	return request[*User](c, "GET", "/users/me", q, nil)
}

// CreateUser creates a single user.
func (c *Client) CreateUser(item *User, q *Query) (*User, error) {
	return request[*User](c, "POST", "/users", q, item)
}

// CreateUsers creates multiple users.
func (c *Client) CreateUsers(items []User, q *Query) ([]User, error) {
	return request[[]User](c, "POST", "/users", q, items)
}

// PatchUser updates a single user by ID.
func (c *Client) PatchUser(id string, item *User, q *Query) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}
	return request[*User](c, "PATCH", fmt.Sprintf("/users/%s", id), q, item)
}

// PatchUsers updates multiple users identified by keys with the same data.
func (c *Client) PatchUsers(keys []string, item *User, q *Query) ([]User, error) {
	body := map[string]any{"keys": keys, "data": item}
	return request[[]User](c, "PATCH", "/users", q, body)
}

// PatchUsersBatch updates multiple users from a batch of items.
func (c *Client) PatchUsersBatch(items []User, q *Query) ([]User, error) {
	return request[[]User](c, "PATCH", "/users", q, items)
}

// PatchUsersMe updates the currently authenticated user.
func (c *Client) PatchUsersMe(item *User, q *Query) (*User, error) {
	return request[*User](c, "PATCH", "/users/me", q, item)
}

// DeleteUser deletes a single user by ID.
func (c *Client) DeleteUser(id string) error {
	if id == "" {
		return fmt.Errorf("id must be provided")
	}
	return c.execute("DELETE", fmt.Sprintf("/users/%s", id), nil, nil)
}

// DeleteUsers deletes multiple users identified by keys.
func (c *Client) DeleteUsers(keys []string) error {
	return c.execute("DELETE", "/users", nil, keys)
}

// InviteUser invites a new user by email to the given role. inviteURL is
// optional; pass an empty string to omit it.
func (c *Client) InviteUser(email, role, inviteURL string) error {
	if email == "" || role == "" {
		return fmt.Errorf("email and role must be provided")
	}
	body := map[string]any{"email": email, "role": role}
	if inviteURL != "" {
		body["invite_url"] = inviteURL
	}
	return c.execute("POST", "/users/invite", nil, body)
}

// AcceptUserInvite accepts a user invite using the invite token and sets the
// account password.
func (c *Client) AcceptUserInvite(token, password string) error {
	if token == "" || password == "" {
		return fmt.Errorf("token and password must be provided")
	}
	body := map[string]any{"token": token, "password": password}
	return c.execute("POST", "/users/invite/accept", nil, body)
}

// RegisterUser registers a new user. Additional fields (e.g. first_name,
// verification_url) may be supplied through opts and are merged into the body.
func (c *Client) RegisterUser(email, password string, opts map[string]any) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password must be provided")
	}
	body := map[string]any{}
	for k, v := range opts {
		body[k] = v
	}
	body["email"] = email
	body["password"] = password
	return c.execute("POST", "/users/register", nil, body)
}

// RegisterUserVerify verifies a registered user's email using the token sent to
// them.
func (c *Client) RegisterUserVerify(token string) error {
	if token == "" {
		return fmt.Errorf("token must be provided")
	}
	return c.execute("GET", fmt.Sprintf("/users/register/verify-email?token=%s", token), nil, nil)
}

// GenerateTwoFactorSecret generates a two-factor authentication secret for the
// current user, returning the secret and otpauth URL.
func (c *Client) GenerateTwoFactorSecret(password string) (map[string]any, error) {
	if password == "" {
		return nil, fmt.Errorf("password must be provided")
	}
	body := map[string]any{"password": password}
	return request[map[string]any](c, "POST", "/users/me/tfa/generate", nil, body)
}

// EnableTwoFactor enables two-factor authentication for the current user.
func (c *Client) EnableTwoFactor(secret, otp string) error {
	if secret == "" || otp == "" {
		return fmt.Errorf("secret and otp must be provided")
	}
	body := map[string]any{"secret": secret, "otp": otp}
	return c.execute("POST", "/users/me/tfa/enable", nil, body)
}

// DisableTwoFactor disables two-factor authentication for the current user.
func (c *Client) DisableTwoFactor(otp string) error {
	if otp == "" {
		return fmt.Errorf("otp must be provided")
	}
	body := map[string]any{"otp": otp}
	return c.execute("POST", "/users/me/tfa/disable", nil, body)
}
