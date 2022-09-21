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
	"time"
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

type CloudFoundry struct {
	// The fully-qualified domain name of the cloud foundry instance you are targeting
	// +optional
	CFApiEndpoint *string `json:"cfApiEndpoint,omitempty"`

	// The username to log into the cloud foundry API
	// +optional
	CFUsername *string `json:"cfUsername,omitempty"`

	// The password to log into cloud foundry API
	// +optional
	CFPassword *string `json:"cfPassword,omitempty"`
}

// ExternalDNSSpec defines the desired state of ExternalDNS
type ExternalDNSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Records *[]BasicInfo `json:"records"`

	// related to kubernetes

	// The kubernetes API server to connect to
	// +optional
	APIServerURL *string `json:"apiServerURL,omitempty"`

	// Path to kubernetes configuration file
	// +optional
	Kubeconfig *string `json:"kubeconfig,omitempty"`

	// Request timeout when calling Kubernetes API. 0s means no timeout
	// +optional
	RequestTimeout *time.Duration `json:"requestTimeout,omitempty"`

	// related to cloud foundry
	// +optional
	CloudFoundry *CloudFoundry `json:"cloudFoundry,omitempty"`

	// RELATED TO CONTOUR
	// The fully-qualified name of the Contour load balancer service. (Default: heptio-contour/contour)
	// +optional
	ContourLoadBalancer *string `json:"contourLoadBalancer,omitempty"`

	// RELATED TO PROCESSING SOURCE

	// The resource types that are queried for endpoints; List of source. ex: source, ingress, node etc.
	// +optional
	Sources *[]string `json:"sources,omitempty"`

	// If source is openshift router then you can pass the ingress controller name. Based on this name the
	// external dns will select the respective router from the route status and map that routeCanonicalHostname
	// to the route host while creating a CNAME record.
	// +optional
	OCRouterName *string `json:"ocRouterName,omitempty"`

	// Limit sources of endpoints to a specific namespace (default: all namespaces)
	// +optional
	Namespace *string `json:"namespace,omitempty"`

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

	// Combine FQDN template and Annotations instead of overwriting
	// +optional
	CombineFQDNAndAnnotation *bool `json:"combineFQDNAndAnnotation,omitempty"`

	// Ignore hostname annotation when generating DNS names, valid only when using fqdn-template is set
	// +optional
	IgnoreHostnameAnnotation *bool `json:"ignoreHostnameAnnotation,omitempty"`

	// Ignore TLS Spec section in ingresses resources, applicable only for ingress source
	// +optional
	IgnoreIngressTLSSpec *bool `json:"ignoreIngressTLSSpec,omitempty"`

	// Ignore rules spec section in ingresses resources, applicable only for ingress sources
	// +optional
	IgnoreIngressRulesSpec *bool `json:"ignoreIngressRulesSpec,omitempty"`

	// Limit Gateways of route endpoints to a specific namespace
	// +optional
	GatewayNamespace *string `json:"gatewayNamespace,omitempty"`

	// Filter Gateways of Route endpoints via label selector
	// +optional
	GatewayLabelFilter *string `json:"gatewayLabelFilter,omitempty"`

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

	// The server to connect for connector source, valid only when using connector source
	// +optional
	ConnectorSourceServer *string `json:"connectorSourceServer,omitempty"`
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

/*
type Config struct {
	APIServerURL                      string
	KubeConfig                        string
	RequestTimeout                    time.Duration
	DefaultTargets                    []string
	ContourLoadBalancerService        string
	GlooNamespace                     string
	SkipperRouteGroupVersion          string
	Sources                           []string
	Namespace                         string
	AnnotationFilter                  string
	LabelFilter                       string
	FQDNTemplate                      string
	CombineFQDNAndAnnotation          bool
	IgnoreHostnameAnnotation          bool
	IgnoreIngressTLSSpec              bool
	IgnoreIngressRulesSpec            bool
	GatewayNamespace                  string
	GatewayLabelFilter                string
	Compatibility                     string
	PublishInternal                   bool
	PublishHostIP                     bool
	AlwaysPublishNotReadyAddresses    bool
	ConnectorSourceServer             string
	Provider                          string
	GoogleProject                     string
	GoogleBatchChangeSize             int
	GoogleBatchChangeInterval         time.Duration
	GoogleZoneVisibility              string
	DomainFilter                      []string
	ExcludeDomains                    []string
	RegexDomainFilter                 *regexp.Regexp
	RegexDomainExclusion              *regexp.Regexp
	ZoneNameFilter                    []string
	ZoneIDFilter                      []string
	AlibabaCloudConfigFile            string
	AlibabaCloudZoneType              string
	AWSZoneType                       string
	AWSZoneTagFilter                  []string
	AWSAssumeRole                     string
	AWSBatchChangeSize                int
	AWSBatchChangeInterval            time.Duration
	AWSEvaluateTargetHealth           bool
	AWSAPIRetries                     int
	AWSPreferCNAME                    bool
	AWSZoneCacheDuration              time.Duration
	AWSSDServiceCleanup               bool
	AzureConfigFile                   string
	AzureResourceGroup                string
	AzureSubscriptionID               string
	AzureUserAssignedIdentityClientID string
	BluecatDNSConfiguration           string
	BluecatConfigFile                 string
	BluecatDNSView                    string
	BluecatGatewayHost                string
	BluecatRootZone                   string
	BluecatDNSServerName              string
	BluecatDNSDeployType              string
	BluecatSkipTLSVerify              bool
	CloudflareProxied                 bool
	CloudflareZonesPerPage            int
	CoreDNSPrefix                     string
	RcodezeroTXTEncrypt               bool
	AkamaiServiceConsumerDomain       string
	AkamaiClientToken                 string
	AkamaiClientSecret                string
	AkamaiAccessToken                 string
	AkamaiEdgercPath                  string
	AkamaiEdgercSection               string
	InfobloxGridHost                  string
	InfobloxWapiPort                  int
	InfobloxWapiUsername              string
	InfobloxWapiPassword              string `secure:"yes"`
	InfobloxWapiVersion               string
	InfobloxSSLVerify                 bool
	InfobloxView                      string
	InfobloxMaxResults                int
	InfobloxFQDNRegEx                 string
	InfobloxCreatePTR                 bool
	InfobloxCacheDuration             int
	DynCustomerName                   string
	DynUsername                       string
	DynPassword                       string `secure:"yes"`
	DynMinTTLSeconds                  int
	OCIConfigFile                     string
	InMemoryZones                     []string
	OVHEndpoint                       string
	OVHApiRateLimit                   int
	PDNSServer                        string
	PDNSAPIKey                        string `secure:"yes"`
	PDNSTLSEnabled                    bool
	TLSCA                             string
	TLSClientCert                     string
	TLSClientCertKey                  string
	Policy                            string
	Registry                          string
	TXTOwnerID                        string
	TXTPrefix                         string
	TXTSuffix                         string
	Interval                          time.Duration
	MinEventSyncInterval              time.Duration
	Once                              bool
	DryRun                            bool
	UpdateEvents                      bool
	LogFormat                         string
	MetricsAddress                    string
	LogLevel                          string
	TXTCacheInterval                  time.Duration
	TXTWildcardReplacement            string
	ExoscaleEndpoint                  string
	ExoscaleAPIKey                    string `secure:"yes"`
	ExoscaleAPISecret                 string `secure:"yes"`
	CRDSourceAPIVersion               string
	CRDSourceKind                     string
	ServiceTypeFilter                 []string
	CFAPIEndpoint                     string
	CFUsername                        string
	CFPassword                        string
	RFC2136Host                       string
	RFC2136Port                       int
	RFC2136Zone                       string
	RFC2136Insecure                   bool
	RFC2136GSSTSIG                    bool
	RFC2136KerberosRealm              string
	RFC2136KerberosUsername           string
	RFC2136KerberosPassword           string `secure:"yes"`
	RFC2136TSIGKeyName                string
	RFC2136TSIGSecret                 string `secure:"yes"`
	RFC2136TSIGSecretAlg              string
	RFC2136TAXFR                      bool
	RFC2136MinTTL                     time.Duration
	RFC2136BatchChangeSize            int
	NS1Endpoint                       string
	NS1IgnoreSSL                      bool
	NS1MinTTLSeconds                  int
	TransIPAccountName                string
	TransIPPrivateKeyFile             string
	DigitalOceanAPIPageSize           int
	ManagedDNSRecordTypes             []string
	GoDaddyAPIKey                     string `secure:"yes"`
	GoDaddySecretKey                  string `secure:"yes"`
	GoDaddyTTL                        int64
	GoDaddyOTE                        bool
	OCPRouterName                     string
	IBMCloudProxied                   bool
	IBMCloudConfigFile                string
}
*/
