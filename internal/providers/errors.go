package providers

import "errors"

var OAuthNotConfiguredError = errors.New("OAuth configuration error: no authentication methods enabled.")

var GoogleEnvMissingError = errors.New("Google authentication is not configured: missing environment variables.")

var GithuEnvMissingError = errors.New("Github authentication is not configured: missing environment variables.")
