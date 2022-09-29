package credentials

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func SetCredential(secret *v1.Secret, key types.NamespacedName, provider string) {
	if provider == "aws" {
		setAWSCredential(secret, key)
	}
}
