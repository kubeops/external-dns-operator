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
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kmapi "kmodules.xyz/client-go/api/v1"
	"time"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// kubebuilder:validation:Enum:=sync;upsert-only;create-only
type Policy string

type ExternalDNSPhase string

// kubebuilder:validation:Enum:=aws;cloudflare
type Provider string

func (p Policy) String() string {
	return string(p)
}

func (p Provider) String() string {
	return string(p)
}

const (
	//Policy
	PolicySync       Policy = "sync"
	PolicyUpsertOnly Policy = "upsert-only"
	PolicyCreateOnly Policy = "create-only"

	//Provider
	ProviderAWS        Provider = "aws"
	providerCloudflare Provider = "cloudflare"

	//ExternalDNSPhase
	ExternalDNSPhaseCurrent    ExternalDNSPhase = "Current"
	ExternalDNSPhaseFailed     ExternalDNSPhase = "Failed"
	ExternalDNSPhaseInProgress ExternalDNSPhase = "InProgress"

	//ConditionType
	ConditionWatcher    = "CreateAndRegisterWatcher"
	ConditionCredential = "CreateAndSetCredential"
	ConditionPlan       = "CreateAndApplyPlan"
)

type TypeInfo struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

func (t TypeInfo) GroupVersionKind() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   t.Group,
		Version: t.Version,
		Kind:    t.Kind,
	}
}

type AWSProvider struct {
	// When using the AWS provider, filter for zones of this type. (support: public, private)
	// +optional
	ZoneType *string `json:"zoneType,omitempty"`

	// When using the AWS provider, filter for zones with these tags
	// +optional
	ZoneTagFilter []string `json:"zoneTagFilter,omitempty"`

	// When using the AWS provider, assume this IAM role. Useful for hosted zones in another AWS account. Specify the
	// full ARN, e.g. `arn:aws:iam::123455567:role/external-dns`
	// +optional
	AssumeRole *string `json:"assumeRole,omitempty"`

	// When using AWS provide, set the maximum number of changes that will be applied in each batch
	// +optional
	BatchChangeSize *int `json:"batchChangeSize,omitempty"`

	// When using the AWS provider, set the interval between batch changes.
	// +optional
	BatchChangeInterval *time.Duration `json:"batchChangeInterval,omitempty"`

	// When using the AWS provider, set whether to evaluate the health of the DNS target (default: enable, disable with --no-aws-evaluate-target-health)
	// +optional
	EvaluateTargetHealth *bool `json:"evaluateTargetHealth,omitempty"`

	// When using the AWS provider, set the maximum number of retries for API calls before giving up.
	// +optional
	APIRetries *int `json:"apiRetries,omitempty"`

	// When using the AWS provider, prefer using CNAME instead of ALIAS (default: disabled)
	// +optional
	PreferCNAME *bool `json:"preferCNAME,omitempty"`

	// When using the AWS provider, set the zones list cache TTL (0s to disable).
	// +optional
	ZoneCacheDuration *time.Duration `json:"zoneCacheDuration,omitempty"`

	// When using the AWS CloudMap provider, delete empty Services without endpoints (default: disabled)
	// +optional
	SDServiceCleanup *bool `json:"sdServiceCleanup,omitempty"`
}

type CloudflareProvider struct {
	// When using the Cloudflare provider, specify if the proxy mode must be enabled (default: disabled)
	// +optional
	Proxied *bool `json:"proxied,omitempty"`

	// When using the Cloudflare provider, specify how many zones per page listed, max. possible 50 (default: 50)
	// +optional
	ZonesPerPage *int `json:"zonesPerPage,omitempty"`
}

type ServiceConfig struct {
	// Limit sources of endpoints to a specific namespace (default: all namespaces)
	// +optional
	Namespace *string `json:"namespace,omitempty"`

	// Ignore hostname annotation when generating DNS names, valid only when using fqdn-template is set
	// +optional
	IgnoreHostnameAnnotation *bool `json:"ignoreHostnameAnnotation,omitempty"`

	// Combine FQDN template and Annotations instead of overwriting
	// +optional
	CombineFQDNAndAnnotation *bool `json:"combineFQDNAndAnnotation,omitempty"`

	// Filter sources managed by external-dns via label selector when listing all resources
	// +optional
	AnnotationFilter *string `json:"annotationFilter,omitempty"`

	// Filter sources managed by external-dns via annotation using label selector semantics
	// +optional
	LabelFilter *string `json:"labelFilter,omitempty"`

	// A templated string that's used to generate DNS names from source that don't define a hostname themselves, or to
	// add a hostname suffix when paired with the fake source
	// +optional
	FQDNTemplate *string `json:"fqdnTemplate,omitempty"`

	// Process  annotation semantics from legacy implementations
	// +optional
	Compatibility *string `json:"compatibility,omitempty"`

	// Allow  externals-dns to publish DNS records for ClusterIP services
	// +optional
	PublishInternal *bool `json:"publishInternal,omitempty"`

	// Allow external-dns to publish host-ip for headless services
	// +optional
	PublishHostIP *bool `json:"publishHostIP,omitempty"`

	// Always publish also not ready addresses for headless services
	// +optional
	AlwaysPublishNotReadyAddresses *bool `json:"alwaysPublishNotReadyAddresses"`

	// The service types to take care about (default all, expected: ClusterIP, NodePort, LoadBalancer or ExternalName)
	// +optional
	ServiceTypeFilter []string `json:"serviceTypeFilter,omitempty"`
}

type IngressConfig struct {
	// Limit sources of endpoints to a specific namespace (default: all namespaces)
	// +optional
	Namespace *string `json:"namespace,omitempty"`

	// Ignore hostname annotation when generating DNS names, valid only when using fqdn-template is set
	// +optional
	IgnoreHostnameAnnotation *bool `json:"ignoreHostnameAnnotation,omitempty"`

	// Combine FQDN template and Annotations instead of overwriting
	// +optional
	CombineFQDNAndAnnotation *bool `json:"combineFQDNAndAnnotation,omitempty"`

	// Filter sources managed by external-dns via label selector when listing all resources
	// +optional
	AnnotationFilter *string `json:"annotationFilter,omitempty"`

	// Filter sources managed by external-dns via annotation using label selector semantics
	// +optional
	LabelFilter *string `json:"labelFilter,omitempty"`

	// A templated string that's used to generate DNS names from source that don't define a hostname themselves, or to
	// add a hostname suffix when paired with the fake source
	// +optional
	FQDNTemplate *string `json:"fqdnTemplate,omitempty"`

	// Ignore TLS Spec section in ingresses resources, applicable only for ingress source
	// +optional
	IgnoreIngressTLSSpec *bool `json:"ignoreIngressTLSSpec,omitempty"`

	// Ignore rules spec section in ingresses resources, applicable only for ingress sources
	// +optional
	IgnoreIngressRulesSpec *bool `json:"ignoreIngressRulesSpec,omitempty"`
}

type NodeConfig struct {
	// A templated string that's used to generate DNS names from source that don't define a hostname themselves, or to
	// add a hostname suffix when paired with the fake source
	FQDNTemplate string `json:"fqdnTemplate,omitempty"`

	// Filter sources managed by external-dns via label selector when listing all resources
	// +optional
	AnnotationFilter *string `json:"annotationFilter,omitempty"`
}

type SourceConfig struct {
	Type TypeInfo `json:"type"`

	// +optional
	Node *NodeConfig `json:"node"`

	// +optional
	Service *ServiceConfig `json:"service"`

	// +optional
	Ingress *IngressConfig `json:"ingress"`
}

// ExternalDNSSpec defines the desired state of ExternalDNS
type ExternalDNSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ProviderSecretRef core.LocalObjectReference `json:"providerSecretRef,omitempty"`

	// Request timeout when calling Kubernetes API. 0s means no timeout
	// +optional
	RequestTimeout *time.Duration `json:"requestTimeout,omitempty"`

	// RELATED TO PROCESSING SOURCE
	// The resource types that are queried for endpoints; List of source. ex: source, ingress, node etc.
	Source SourceConfig `json:"source"`
	// sources:
	//    - group: ""
	//      version: v1
	//      kind: Service
	//    - group: ""
	//      version: v1
	//      kind: Node `

	// If source is openshift router then you can pass the ingress controller name. Based on this name the
	// external dns will select the respective router from the route status and map that routeCanonicalHostname
	// to the route host while creating a CNAME record.
	// +optional
	OCRouterName *string `json:"ocRouterName,omitempty"`

	// Limit Gateways of route endpoints to a specific namespace
	// +optional
	GatewayNamespace *string `json:"gatewayNamespace,omitempty"`

	// Filter Gateways of Route endpoints via label selector
	// +optional
	GatewayLabelFilter *string `json:"gatewayLabelFilter,omitempty"`

	// The server to connect for connector source, valid only when using connector source
	// +optional
	ConnectorSourceServer *string `json:"connectorSourceServer,omitempty"`

	// Comma separated list of record types to manage (default: A, CNAME; supported: A,CNAME,NS)
	// +optional
	ManageDNSRecordTypes []string `json:"manageDNSRecordTypes,omitempty"`

	// Set globally a list of default IP address that will apply as a target instead of source addresses.
	// +optional
	DefaultTargets []string `json:"defaultTargets,omitempty"`

	//

	// RELATED TO PROVIDERS
	// The DNS provider where the DNS records will be created. (AWS, Cloudflare)
	Provider Provider `json:"provider"`

	// Limit possible target zones by a domain suffix
	// +optional
	DomainFilter []string `json:"domainFilter,omitempty"`

	// Exclude subdomains
	// +optional
	ExcludeDomains []string `json:"excludeDomains,omitempty"`

	// Filter target zones by hosted zone id
	// +optional
	ZoneIDFilter []string `json:"zoneIDFilter,omitempty"`

	// AWS provider information
	// +optional
	AWS *AWSProvider `json:"aws,omitempty"`

	// Cloudflare provider information
	// +optional
	Cloudflare *CloudflareProvider `json:"cloudflare,omitempty"`

	//
	//POLICY INFORMATION
	//
	// Modify how DNS records are synchronized between sources and providers (default: sync, options: sync, upsert-only, create-only)
	// +optional
	Policy *Policy `json:"policy,omitempty"`

	//
	// REGISTRY information
	//
	// The registry implementation to use to keep track of DNS record ownership (default: txt, options: txt, noop, aws-sd)
	// +optional
	Registry *string `json:"registry,omitempty"`

	// When using the TXT registry, a name that identifies this instance of ExternalDNS (default: default)
	// +optional
	TXTOwnerID *string `json:"txtOwnerID,omitempty"`

	// When using the TXT registry, a custom string that's prefixed to each ownership DNS record (optional). Could
	// contain record type template like '%{record_type}-prefix-'. Mutual exclusive with txt-suffix!
	// +optional
	TXTPrefix *string `json:"txtPrefix,omitempty"`

	// When using the TXT registry, a custom string that's suffixed to the host portion of each ownership DNS
	// record. Could contain record type template like '-%{record_type}-suffix'. Mutual exclusive with txt-prefix!
	// +optional
	TXTSuffix *string `json:"txtSuffix,omitempty"`

	// When using the TXT registry, a custom string that's used instead of an asterisk for TXT records corresponding
	// to wildcard DNS records
	// +optional
	TXTWildcardReplacement *string `json:"txtWildcardReplacement,omitempty"`
}

// ExternalDNSStatus defines the observed state of ExternalDNS
type ExternalDNSStatus struct {
	// +optional
	Phase ExternalDNSPhase `json:"phase,omitempty"`

	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
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
