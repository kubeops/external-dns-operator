/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type BasicInfo struct {
	Source     *string `json:"source"`
	Domain     *string `json:"domain"`
	Provider   *string `json:"provider"`
	Policy     *string `json:"policy"`
	AWSZone    *string `json:"aws_zone"`
	Registry   *string `json:"registry"`
	TxtOwnerID *string `json:"txt_owner_id"`
	TxtPrefix  *string `json:"txt_prefix"`
	/*
	 */
}

// ExternalDNSSpec defines the desired state of ExternalDNS
type ExternalDNSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Records *[]BasicInfo `json:"records"`

	/*
		// related to kubernetes
		APIServerURL         *string `json:"api_server_url"`
		Kubeconfig     *string        `json:"kubeconfig"`
		RequestTimeout *time.Duration `json:"request_timeout"`

		// related to cloud foundry
		CFApiEndpoint *string `json:"cf_api_endpoint"`
		CFUsername    *string `json:"cf_username"`
		CFPassword    *string `json:"cf_password"`

		//related to contour

	*/

}

// ExternalDNSStatus defines the observed state of ExternalDNS
type ExternalDNSStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ExternalDNS is the Schema for the externaldns API
type ExternalDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalDNSSpec   `json:"spec,omitempty"`
	Status ExternalDNSStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExternalDNSList contains a list of ExternalDNS
type ExternalDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalDNS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExternalDNS{}, &ExternalDNSList{})
}
