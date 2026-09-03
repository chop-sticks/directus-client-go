package directus

import "fmt"

// AuthenticationData holds the credentials returned by the authentication
// endpoints (login, refresh) and share authentication.
type AuthenticationData struct {
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

// AuthProvider describes a configured authentication provider as returned by
// GET /auth.
type AuthProvider struct {
	Name   string  `json:"name"`
	Driver string  `json:"driver"`
	Label  *string `json:"label"`
	Icon   *string `json:"icon"`
}

// Login authenticates with email and password, returning access and refresh
// tokens. mode defaults to "json" when empty; otp is only sent when non-empty.
func (c *Client) Login(email, password, mode, otp string) (*AuthenticationData, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password must be provided")
	}
	if mode == "" {
		mode = "json"
	}
	body := map[string]any{
		"email":    email,
		"password": password,
		"mode":     mode,
	}
	if otp != "" {
		body["otp"] = otp
	}
	return request[*AuthenticationData](c, "POST", "/auth/login", nil, body)
}

// Refresh exchanges a refresh token for a new set of tokens. mode defaults to
// "json" when empty; refresh_token is only sent when non-empty (session mode
// reads it from a cookie).
func (c *Client) Refresh(refreshToken, mode string) (*AuthenticationData, error) {
	if mode == "" {
		mode = "json"
	}
	body := map[string]any{
		"mode": mode,
	}
	if refreshToken != "" {
		body["refresh_token"] = refreshToken
	}
	return request[*AuthenticationData](c, "POST", "/auth/refresh", nil, body)
}

// Logout invalidates the current session. mode defaults to "json" when empty;
// refresh_token is only sent when non-empty.
func (c *Client) Logout(refreshToken, mode string) error {
	if mode == "" {
		mode = "json"
	}
	body := map[string]any{
		"mode": mode,
	}
	if refreshToken != "" {
		body["refresh_token"] = refreshToken
	}
	return c.execute("POST", "/auth/logout", nil, body)
}

// PasswordRequest triggers a password reset email for the given address. When
// resetURL is non-empty it overrides the configured reset page.
func (c *Client) PasswordRequest(email, resetURL string) error {
	if email == "" {
		return fmt.Errorf("email must be provided")
	}
	body := map[string]any{
		"email": email,
	}
	if resetURL != "" {
		body["reset_url"] = resetURL
	}
	return c.execute("POST", "/auth/password/request", nil, body)
}

// PasswordReset completes a password reset using the token from the reset email.
func (c *Client) PasswordReset(token, password string) error {
	if token == "" || password == "" {
		return fmt.Errorf("token and password must be provided")
	}
	body := map[string]any{
		"token":    token,
		"password": password,
	}
	return c.execute("POST", "/auth/password/reset", nil, body)
}

// ReadProviders lists the configured authentication providers. When sessionOnly
// is true, only providers usable for session-based auth are returned.
func (c *Client) ReadProviders(sessionOnly bool) ([]AuthProvider, error) {
	path := "/auth"
	if sessionOnly {
		path = "/auth?sessionOnly"
	}
	return request[[]AuthProvider](c, "GET", path, nil, nil)
}
