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

	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
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
