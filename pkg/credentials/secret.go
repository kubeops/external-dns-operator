package credentials

import (
	"context"
	"errors"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getSecret(ctx context.Context, kc client.Client, key types.NamespacedName) (*core.Secret, error) {
	secret := &core.Secret{}
	if err := kc.Get(ctx, key, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func resetEnvVariables(list ...string) error {
	for _, item := range list {
		err := os.Setenv(item, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func SetCredential(ctx context.Context, kc client.Client, edns *externaldnsv1alpha1.ExternalDNS) error {

	switch edns.Spec.Provider.String() {
	case externaldnsv1alpha1.ProviderAWS.String():
		return setAWSCredential(ctx, kc, edns)

	case externaldnsv1alpha1.ProviderCloudflare.String():
		return setCloudflareCredentials(ctx, kc, edns)

	case externaldnsv1alpha1.ProviderAzure.String():
		return setAzureCredential(ctx, kc, edns)

	case externaldnsv1alpha1.ProviderGoogle.String():
		return setGoogleCredential(ctx, kc, edns)

	default:
		return errors.New("unknown provider name")
	}
}
