package providers

import (
	"os"
	"testing"
)

func TestProvider(t *testing.T) {
	t.Run("Should receive error if none OAUTH is enable", func(t *testing.T) {
		os.Setenv("GOOGLE_ENABLE", "0")
		os.Setenv("GITHUB_ENABLE", "0")

		_, err := NewOAuthOptions()

		if err != OAuthNotConfiguredError {
			t.Errorf("Expect -> %v, Receive -> %v", OAuthNotConfiguredError, err)
		}

	})

	t.Run("Should receive error if some Google env is missing", func(t *testing.T) {
		os.Setenv("GOOGLE_ENABLE", "1")

		_, err := NewOAuthOptions()

		if err != GoogleEnvMissingError {
			t.Errorf("Expect -> %v, Receive -> %v", GoogleEnvMissingError, err)
		}

	})
}
