package credentials

import (
	"errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
)

func SetCredential(secret *v1.Secret, ednsKey types.NamespacedName, provider string) error {

	switch provider {
	case externaldnsv1alpha1.ProviderAWS.String():
		return setAWSCredential(secret, ednsKey)

	case externaldnsv1alpha1.ProviderCloudflare.String():
		// set environment variable
		return setCloudflareCredentials(secret)

	default:
		return errors.New("unknown provider name")
	}
}
