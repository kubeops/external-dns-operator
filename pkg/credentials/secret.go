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

func getSecret(ctx context.Context, kc client.Client, key types.NamespacedName) (*core.Secret, error) {
	secret := &core.Secret{}
	if err := kc.Get(ctx, key, secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func resetEnvVariables(list ...string) error {
	for _, item := range list {
		if err := os.Unsetenv(item); err != nil {
			return err
		}
	}
	return nil
}

// providerEnvVars enumerates every environment variable any of the
// supported credential setters may write. Listing them in one place lets
// SetCredential clear stale entries left over from a previous provider
// before configuring the new one, so switching providers (or deleting an
// old ExternalDNS and creating a new one) cannot leak credentials across
// reconciles.
var providerEnvVars = []string{
	AWSSharedCredentialsFile,
	GoogleApplicationCredentials,
	CFBaseURL,
	CFApiToken,
	CFApiKey,
	CFApiEmail,
}

// credentialFilePath returns the on-disk path used by the file-based
// provider credential setters (AWS / Azure / Google). It must stay in
// sync with the path each setter writes to.
func credentialFilePath(edns *api.ExternalDNS) string {
	return fmt.Sprintf("/tmp/%s-%s-credential", edns.Namespace, edns.Name)
}

// CleanupCredential removes any on-disk credential files written by a
// previous SetCredential call for this ExternalDNS. Safe to call when
// the file does not exist. Intended to run during finalizer-driven
// deletion so we don't leak secret material on /tmp across the lifetime
// of the operator pod.
func CleanupCredential(edns *api.ExternalDNS) error {
	switch edns.Spec.Provider {
	case api.ProviderAWS, api.ProviderAzure, api.ProviderGoogle:
		if err := os.Remove(credentialFilePath(edns)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func SetCredential(ctx context.Context, kc client.Client, edns *api.ExternalDNS) error {
	if err := resetEnvVariables(providerEnvVars...); err != nil {
		return err
	}

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
