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

const (
	CFApiToken = "CF_API_TOKEN"
	CFApiKey   = "CF_API_KEY"
	CFApiEmail = "CF_API_EMAIL"
)

func validCFSecret(secret *core.Secret) bool {

	if _, foundToken := secret.Data[CFApiToken]; foundToken {
		return true
	} else {
		_, foundKey := secret.Data[CFApiKey]
		_, foundEmail := secret.Data[CFApiEmail]

		return foundKey && foundEmail
	}
}

func setCloudflareCredentials(ctx context.Context, kc client.Client, edns *externaldnsv1alpha1.ExternalDNS) error {

	if err := resetEnvVariables(CFApiToken, CFApiKey, CFApiEmail); err != nil {
		return err
	}

	// ProviderSecretRef is required for cloudflare
	if edns.Spec.ProviderSecretRef == nil {
		return errors.New("providerSecretRef is not given for cloudflare provider")
	}

	secret, err := getSecret(ctx, kc, types.NamespacedName{Namespace: edns.Namespace, Name: edns.Spec.ProviderSecretRef.Name})
	if err != nil {
		return err
	}

	if !validCFSecret(secret) {
		return errors.New("invalid cloudflare provider secret")
	}

	if string(secret.Data[CFApiToken][:]) != "" {
		return os.Setenv("CF_API_TOKEN", string(secret.Data["CF_API_TOKEN"][:]))
	} else {
		if err := os.Setenv(CFApiKey, string(secret.Data[CFApiKey][:])); err != nil {
			return err
		}
		if err := os.Setenv(CFApiEmail, string(secret.Data[CFApiEmail][:])); err != nil {
			return err
		}
	}
	return nil
}
