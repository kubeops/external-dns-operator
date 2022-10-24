package credentials

import (
	"errors"
	core "k8s.io/api/core/v1"
	"os"
)

func validCFSecret(secret *core.Secret) bool {

	if _, foundToken := secret.Data["CF_API_TOKEN"]; foundToken {
		return true
	} else {
		_, foundKey := secret.Data["CF_API_KEY"]
		_, foundEmail := secret.Data["CF_API_EMAIL"]

		return foundKey && foundEmail
	}
}

func setCloudflareCredentials(secret *core.Secret) error {

	if !validCFSecret(secret) {
		return errors.New("invalid cloudflare provider secret")
	}

	if string(secret.Data["CF_API_TOKEN"][:]) != "" {
		return os.Setenv("CF_API_TOKEN", string(secret.Data["CF_API_TOKEN"][:]))
	} else {
		if err := os.Setenv("CF_API_KEY", string(secret.Data["CF_API_KEY"][:])); err != nil {
			return err
		}
		if err := os.Setenv("CF_API_EMAIL", string(secret.Data["CF_API_EMAIL"][:])); err != nil {
			return err
		}
	}
	return nil
}
