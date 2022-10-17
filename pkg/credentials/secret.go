package credentials

import (
	"errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"kubeops.dev/external-dns-operator/pkg/constant"
)

func SetCredential(secret *v1.Secret, ednsKey types.NamespacedName, provider string) error {

	switch provider {
	case constant.ProviderAWS.String():
		return setAWSCredential(secret, ednsKey)

	case constant.ProviderCloudflare.String():
		// set environment variable
		return setCloudflareCredentials(secret)

	default:
		return errors.New("unknown provider name")
	}
}
