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

func validCFSecret(secret *core.Secret) bool {

	if _, foundToken := secret.Data["CF_API_TOKEN"]; foundToken {
		return true
	} else {
		_, foundKey := secret.Data["CF_API_KEY"]
		_, foundEmail := secret.Data["CF_API_EMAIL"]

		return foundKey && foundEmail
	}
}

func setCloudflareCredentials(ctx context.Context, kc client.Client, edns *externaldnsv1alpha1.ExternalDNS) error {

	if err := resetEnvVariables("CF_API_TOKEN", "CF_API_KEY", "CF_API_EMAIL"); err != nil {
		return err
	}

	// if ProviderSecretRef is nil then user is intended to use IRSA (IAM Role for Service Account)
	if edns.Spec.ProviderSecretRef == nil {
		// handle for not providing the providerSecretRef
		// probably clear the environment variables
		return nil
	}

	secret, err := getSecret(ctx, kc, types.NamespacedName{Namespace: edns.Namespace, Name: edns.Spec.ProviderSecretRef.Name})
	if err != nil {
		return err
	}

	if !validCFSecret(secret) {
		return errors.New("invalid cloudflare provider secret")
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
