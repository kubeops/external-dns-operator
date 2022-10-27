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

func validAzureSecret(secret *core.Secret) bool {
	_, found := secret.Data["azure.json"]
	return found
}

func setAzureCredential(ctx context.Context, kc client.Client, edns *externaldnsv1alpha1.ExternalDNS) error {

	namespace := edns.Namespace
	name := edns.Name

	// for azure, user must have to provide ProviderSecretRef
	if edns.Spec.ProviderSecretRef == nil {
		return errors.New("invalid ProviderSecretRef for azure provider")
	}

	secret, err := getSecret(ctx, kc, types.NamespacedName{Namespace: namespace, Name: edns.Spec.ProviderSecretRef.Name})
	if err != nil {
		return err
	}

	if !validAzureSecret(secret) {
		return errors.New("invalid Azure provider secret")
	}
	fileName := fmt.Sprintf("%s-%s-credential", namespace, name)
	filepath := fmt.Sprintf("/tmp/%s", fileName)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	b := secret.Data["azure.json"]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
