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

const (
	CFBaseURL  = "CF_BASE_URL"
	CFApiToken = "CF_API_TOKEN"
	CFApiKey   = "CF_API_KEY"
	CFApiEmail = "CF_API_EMAIL"
)

type cfAuthMode int

const (
	cfAuthInvalid cfAuthMode = iota
	cfAuthToken
	cfAuthKeyAndEmail
)

// cfSecretMode picks which Cloudflare auth flavor to use given the
// referenced secret. Token-based auth wins when the token key is
// present and non-empty; otherwise key+email is used when both are
// present and non-empty; otherwise the secret is invalid.
func cfSecretMode(secret *core.Secret, tokenKey, apiKey, apiEmail string) cfAuthMode {
	if len(secret.Data[tokenKey]) > 0 {
		return cfAuthToken
	}
	if len(secret.Data[apiKey]) > 0 && len(secret.Data[apiEmail]) > 0 {
		return cfAuthKeyAndEmail
	}
	return cfAuthInvalid
}

func setCloudflareCredentials(ctx context.Context, kc client.Client, edns *api.ExternalDNS) error {
	if err := resetEnvVariables(CFApiToken, CFApiKey, CFApiEmail, CFBaseURL); err != nil {
		return err
	}

	// ProviderSecretRef is required for cloudflare
	if edns.Spec.Cloudflare == nil || edns.Spec.Cloudflare.SecretRef == nil {
		return errors.New("providerSecretRef is not given for cloudflare provider")
	}

	ref := edns.Spec.Cloudflare.SecretRef
	secret, err := getSecret(ctx, kc, types.NamespacedName{Namespace: edns.Namespace, Name: ref.Name})
	if err != nil {
		return err
	}

	if edns.Spec.Cloudflare.BaseURL != "" {
		if err := os.Setenv(CFBaseURL, edns.Spec.Cloudflare.BaseURL); err != nil {
			return err
		}
	}

	switch cfSecretMode(secret, ref.APITokenKey, ref.APIKey, ref.APIEmailKey) {
	case cfAuthToken:
		return os.Setenv(CFApiToken, string(secret.Data[ref.APITokenKey]))
	case cfAuthKeyAndEmail:
		if err := os.Setenv(CFApiKey, string(secret.Data[ref.APIKey])); err != nil {
			return err
		}
		return os.Setenv(CFApiEmail, string(secret.Data[ref.APIEmailKey]))
	default:
		return errors.New("invalid cloudflare provider secret")
	}
}
