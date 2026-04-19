package providers

import "errors"

var GoogleEnvMissingError = errors.New("Google authentication is not configured: missing environment variables.")

var OAuthNotConfiguredError = errors.New("OAuth configuration error: no authentication methods enabled.")
