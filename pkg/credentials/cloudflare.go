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

	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

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

func validCFSecret(secret *core.Secret, tokenKey, apiKey, apiEmail string) bool {
	if _, foundToken := secret.Data[tokenKey]; foundToken {
		return true
	} else {
		_, foundKey := secret.Data[apiKey]
		_, foundEmail := secret.Data[apiEmail]

		return foundKey && foundEmail
	}
}

func setCloudflareCredentials(ctx context.Context, kc client.Client, edns *externaldnsv1alpha1.ExternalDNS) error {
	if err := resetEnvVariables(CFApiToken, CFApiKey, CFApiEmail, CFBaseURL); err != nil {
		return err
	}

	// ProviderSecretRef is required for cloudflare
	if edns.Spec.Cloudflare == nil || edns.Spec.Cloudflare.SecretRef == nil {
		return errors.New("providerSecretRef is not given for cloudflare provider")
	}

	secret, err := getSecret(ctx, kc, types.NamespacedName{Namespace: edns.Namespace, Name: edns.Spec.Cloudflare.SecretRef.Name})
	if err != nil {
		return err
	}

	if !validCFSecret(secret, edns.Spec.Cloudflare.SecretRef.APITokenKey, edns.Spec.Cloudflare.SecretRef.APIKey, edns.Spec.Cloudflare.SecretRef.APIEmailKey) {
		return errors.New("invalid cloudflare provider secret")
	}

	if edns.Spec.Cloudflare.BaseURL != "" {
		if err := os.Setenv(CFBaseURL, edns.Spec.Cloudflare.BaseURL); err != nil {
			return err
		}
	}
	if len(secret.Data[edns.Spec.Cloudflare.SecretRef.APITokenKey]) > 0 {
		if err := os.Setenv(CFApiToken, string(secret.Data[edns.Spec.Cloudflare.SecretRef.APITokenKey])); err != nil {
			return err
		}
	} else {
		if err := os.Setenv(CFApiKey, string(secret.Data[edns.Spec.Cloudflare.SecretRef.APIKey])); err != nil {
			return err
		}
		if err := os.Setenv(CFApiEmail, string(secret.Data[edns.Spec.Cloudflare.SecretRef.APIEmailKey])); err != nil {
			return err
		}
	}
	return nil
}
