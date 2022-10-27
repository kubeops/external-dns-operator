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

func validAWSSecret(secret *core.Secret) bool {
	_, found := secret.Data["credentials"]
	return found
}

func setAWSCredential(ctx context.Context, kc client.Client, edns *externaldnsv1alpha1.ExternalDNS) error {

	if err := resetEnvVariables("AWS_SHARED_CREDENTIALS_FILE"); err != nil {
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

	if !validAWSSecret(secret) {
		return errors.New("invalid aws provider secret")
	}

	fileName := fmt.Sprintf("%s-%s-credential", edns.Namespace, edns.Name)

	filePath := fmt.Sprintf("/tmp/%s", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b := secret.Data["credentials"]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	err = os.Setenv("AWS_SHARED_CREDENTIALS_FILE", filePath)
	if err != nil {
		return err
	}
	return nil
}
