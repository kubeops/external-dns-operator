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
	"os"

	api "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
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

func SetCredential(ctx context.Context, kc client.Client, edns *api.ExternalDNS) error {
	switch edns.Spec.Provider {
	case api.ProviderAWS:
		return setAWSCredential(ctx, kc, edns)

	case api.ProviderCloudflare:
		return setCloudflareCredentials(ctx, kc, edns)

	case api.ProviderAzure:
		return setAzureCredential(ctx, kc, edns)

	case api.ProviderGoogle:
		return setGoogleCredential(ctx, kc, edns)

	default:
		return errors.New("unknown provider name")
	}
}
