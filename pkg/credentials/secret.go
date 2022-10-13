package credentials

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

func SetCredential(secret *v1.Secret, key types.NamespacedName, provider string) error {

	switch provider {
	case "aws":
		err := setAWSCredential(secret, key)
		if err != nil {
			klog.Info("failed to set credential, ", err.Error())
			return err
		}
	case "cloudflare":
		// set environment variable
		if err := setEnvVar("CF_API_KEY", string(secret.Data["CF_API_KEY"][:])); err != nil {
			klog.Info("failed to set environment variables")
			return err
		}

		if err := setEnvVar("CF_API_EMAIL", string(secret.Data["CF_API_EMAIL"][:])); err != nil {
			klog.Info("failed to set environment variables")
			return err
		}
	default:
		klog.Info("unknown provider name")
	}

	return nil
}
