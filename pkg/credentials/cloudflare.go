package credentials

import (
	core "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"os"
)

func setCloudflareCredentials(secret *core.Secret) error {

	//ignored the CF_API_TOKEN

	if err := os.Setenv("CF_API_KEY", string(secret.Data["CF_API_KEY"][:])); err != nil {
		klog.Error("failed to set environment variables")
		return err
	}

	if err := os.Setenv("CF_API_EMAIL", string(secret.Data["CF_API_EMAIL"][:])); err != nil {
		klog.Error("failed to set environment variables")
		return err
	}
	return nil
}
