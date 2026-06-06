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
	"testing"

	core "k8s.io/api/core/v1"
)

func TestValidCFSecret(t *testing.T) {
	const (
		tokenKey = "api-token"
		apiKey   = "api-key"
		apiEmail = "api-email"
	)

	tests := []struct {
		name string
		data map[string][]byte
		want bool
	}{
		{
			name: "token present is sufficient",
			data: map[string][]byte{tokenKey: []byte("tok")},
			want: true,
		},
		{
			name: "key+email present is sufficient",
			data: map[string][]byte{apiKey: []byte("k"), apiEmail: []byte("e@example.com")},
			want: true,
		},
		{
			name: "key only is not sufficient",
			data: map[string][]byte{apiKey: []byte("k")},
			want: false,
		},
		{
			name: "email only is not sufficient",
			data: map[string][]byte{apiEmail: []byte("e@example.com")},
			want: false,
		},
		{
			name: "empty secret is invalid",
			data: map[string][]byte{},
			want: false,
		},
		{
			// The current implementation treats a present-but-empty token
			// key as sufficient; lock that in so the behavior is explicit.
			name: "empty-but-present token still passes",
			data: map[string][]byte{tokenKey: nil},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := validCFSecret(&core.Secret{Data: tc.data}, tokenKey, apiKey, apiEmail)
			if got != tc.want {
				t.Fatalf("validCFSecret = %v, want %v (data=%v)", got, tc.want, tc.data)
			}
		})
	}
}

func TestValidAWSSecret(t *testing.T) {
	const key = "credentials"

	tests := []struct {
		name string
		data map[string][]byte
		want bool
	}{
		{name: "key present", data: map[string][]byte{key: []byte("...")}, want: true},
		{name: "key absent", data: map[string][]byte{"other": []byte("...")}, want: false},
		{name: "empty data", data: map[string][]byte{}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := validAWSSecret(&core.Secret{Data: tc.data}, key); got != tc.want {
				t.Fatalf("validAWSSecret = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestValidGoogleSecret(t *testing.T) {
	const key = "service-account.json"

	if !validGoogleSecret(&core.Secret{Data: map[string][]byte{key: []byte(`{"type":"x"}`)}}, key) {
		t.Fatalf("expected key present to validate")
	}
	if validGoogleSecret(&core.Secret{Data: map[string][]byte{}}, key) {
		t.Fatalf("expected empty data to fail")
	}
}

func TestValidAzureSecret(t *testing.T) {
	const key = "azure.json"

	if !validAzureSecret(&core.Secret{Data: map[string][]byte{key: []byte("{}")}}, key) {
		t.Fatalf("expected key present to validate")
	}
	if validAzureSecret(&core.Secret{Data: map[string][]byte{}}, key) {
		t.Fatalf("expected empty data to fail")
	}
}
