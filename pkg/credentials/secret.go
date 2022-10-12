package credentials

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

func SetCredential(secret *v1.Secret, key types.NamespacedName, provider string) error {
	if provider == "aws" {
		err := setAWSCredential(secret, key)
		if err != nil {
			klog.Info("failed to set credential, ", err.Error())
			return err
		}
	}

	return nil
}
