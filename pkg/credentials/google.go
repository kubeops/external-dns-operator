package credentials

import (
	"context"
	"errors"
	"fmt"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const GoogleApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"

func validGoogleSecret(secret *core.Secret) bool {
	_, found := secret.Data["credentials.json"]
	return found
}

func setGoogleCredential(ctx context.Context, kc client.Client, edns *externaldnsv1alpha1.ExternalDNS) error {

	if err := resetEnvVariables(GoogleApplicationCredentials); err != nil {
		return err
	}

	// if ProviderSecretRef is nil then user is intended to use Workload Identity
	if edns.Spec.ProviderSecretRef == nil {
		return nil
	}

	secret, err := getSecret(ctx, kc, types.NamespacedName{Namespace: edns.Namespace, Name: edns.Spec.ProviderSecretRef.Name})
	if err != nil {
		return err
	}

	if !validGoogleSecret(secret) {
		return errors.New("invalid Google provider secret")
	}
	fileName := fmt.Sprintf("%s-%s-credential", edns.Namespace, edns.Name)
	filePath := fmt.Sprintf("/tmp/%s", fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b := secret.Data["credentials.json"]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	err = os.Setenv(GoogleApplicationCredentials, filePath)
	if err != nil {
		return err
	}

	return nil
}
