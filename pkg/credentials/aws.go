/*
Copyright AppsCode Inc. and Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"

	api "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const AWSSharedCredentialsFile = "AWS_SHARED_CREDENTIALS_FILE"

func validAWSSecret(secret *core.Secret, key string) bool {
	_, found := secret.Data[key]
	return found
}

func setAWSCredential(ctx context.Context, kc client.Client, edns *api.ExternalDNS) error {
	if err := resetEnvVariables(AWSSharedCredentialsFile); err != nil {
		return err
	}

	// if ProviderSecretRef is nil then user is intended to use IRSA (IAM Role for Service Account)
	if edns.Spec.AWS == nil || edns.Spec.AWS.SecretRef == nil {
		return nil
	}

	secret, err := getSecret(ctx, kc, types.NamespacedName{Namespace: edns.Namespace, Name: edns.Spec.AWS.SecretRef.Name})
	if err != nil {
		return err
	}

	if !validAWSSecret(secret, edns.Spec.AWS.SecretRef.CredentialKey) {
		return errors.New("invalid aws provider secret")
	}

	fileName := fmt.Sprintf("%s-%s-credential", edns.Namespace, edns.Name)

	filePath := fmt.Sprintf("/tmp/%s", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close() // nolint:errcheck

	b := secret.Data[edns.Spec.AWS.SecretRef.CredentialKey]
	_, err = file.Write(b)
	if err != nil {
		return err
	}

	err = os.Setenv(AWSSharedCredentialsFile, filePath)
	if err != nil {
		return err
	}
	return nil
}
