package credentials

import (
	core "k8s.io/api/core/v1"
	"os"
)

func setCloudflareCredentials(secret *core.Secret) error {

	//ignored the CF_API_TOKEN

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
