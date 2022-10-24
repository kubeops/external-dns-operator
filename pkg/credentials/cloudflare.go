package credentials

import (
	"errors"
	core "k8s.io/api/core/v1"
	"os"
)

func validateCFSecret(secret *core.Secret) error {
	if secret.Data["CF_API_TOKEN"] != nil {
		return nil
	} else if secret.Data["CF_API_KEY"] != nil && secret.Data["CF_API_EMAIL"] != nil {
		return nil
	} else {
		return errors.New("invalid secret format(s)")
	}
}

func setCloudflareCredentials(secret *core.Secret) error {

	if err := validateCFSecret(secret); err != nil {
		return err
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
