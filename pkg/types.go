package pkg

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"regexp"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/pkg/apis/externaldns"
	"time"
)

/*
// overriding the config structure
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
	TargetNetFilter                   []string
	ExcludeTargetNets                 []string
	AlibabaCloudConfigFile            string
	AlibabaCloudZoneType              string
	AWSZoneType                       string
	AWSZoneTagFilter                  []string
	AWSAssumeRole                     string
	AWSAssumeRoleExternalID           string
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
var defaultConfig = &externaldns.Config{
	APIServerURL:                "",
	KubeConfig:                  "",
	RequestTimeout:              time.Second * 30,
	DefaultTargets:              []string{},
	ContourLoadBalancerService:  "heptio-contour/contour",
	GlooNamespace:               "gloo-system",
	SkipperRouteGroupVersion:    "zalando.org/v1",
	Sources:                     nil,
	Namespace:                   "",
	AnnotationFilter:            "",
	LabelFilter:                 labels.Everything().String(),
	FQDNTemplate:                "",
	CombineFQDNAndAnnotation:    false,
	IgnoreHostnameAnnotation:    false,
	IgnoreIngressTLSSpec:        false,
	IgnoreIngressRulesSpec:      false,
	GatewayNamespace:            "",
	GatewayLabelFilter:          "",
	Compatibility:               "",
	PublishInternal:             false,
	PublishHostIP:               false,
	ConnectorSourceServer:       "localhost:8080",
	Provider:                    "",
	GoogleProject:               "",
	GoogleBatchChangeSize:       1000,
	GoogleBatchChangeInterval:   time.Second,
	GoogleZoneVisibility:        "",
	DomainFilter:                []string{},
	ExcludeDomains:              []string{},
	RegexDomainFilter:           regexp.MustCompile(""),
	RegexDomainExclusion:        regexp.MustCompile(""),
	TargetNetFilter:             []string{}, //config type er sob fields ei tinta ken pailo na?
	ExcludeTargetNets:           []string{},
	AlibabaCloudConfigFile:      "/etc/kubernetes/alibaba-cloud.json",
	AWSZoneType:                 "",
	AWSZoneTagFilter:            []string{},
	AWSAssumeRole:               "",
	AWSAssumeRoleExternalID:     "",
	AWSBatchChangeSize:          1000,
	AWSBatchChangeInterval:      time.Second,
	AWSEvaluateTargetHealth:     true,
	AWSAPIRetries:               3,
	AWSPreferCNAME:              false,
	AWSZoneCacheDuration:        0 * time.Second,
	AWSSDServiceCleanup:         false,
	AzureConfigFile:             "/etc/kubernetes/azure.json",
	AzureResourceGroup:          "",
	AzureSubscriptionID:         "",
	BluecatConfigFile:           "/etc/kubernetes/bluecat.json",
	BluecatDNSDeployType:        "no-deploy",
	CloudflareProxied:           false,
	CloudflareZonesPerPage:      50,
	CoreDNSPrefix:               "/skydns/",
	RcodezeroTXTEncrypt:         false,
	AkamaiServiceConsumerDomain: "",
	AkamaiClientToken:           "",
	AkamaiClientSecret:          "",
	AkamaiAccessToken:           "",
	AkamaiEdgercSection:         "",
	AkamaiEdgercPath:            "",
	InfobloxGridHost:            "",
	InfobloxWapiPort:            443,
	InfobloxWapiUsername:        "admin",
	InfobloxWapiPassword:        "",
	InfobloxWapiVersion:         "2.3.1",
	InfobloxSSLVerify:           true,
	InfobloxView:                "",
	InfobloxMaxResults:          0,
	InfobloxFQDNRegEx:           "",
	InfobloxCreatePTR:           false,
	InfobloxCacheDuration:       0,
	OCIConfigFile:               "/etc/kubernetes/oci.yaml",
	InMemoryZones:               []string{},
	OVHEndpoint:                 "ovh-eu",
	OVHApiRateLimit:             20,
	PDNSServer:                  "http://localhost:8081",
	PDNSAPIKey:                  "",
	PDNSTLSEnabled:              false,
	TLSCA:                       "",
	TLSClientCert:               "",
	TLSClientCertKey:            "",
	Policy:                      "sync",
	Registry:                    "txt",
	TXTOwnerID:                  "default",
	TXTPrefix:                   "",
	TXTSuffix:                   "",
	TXTCacheInterval:            0,
	TXTWildcardReplacement:      "",
	MinEventSyncInterval:        5 * time.Second,
	Interval:                    time.Minute,
	Once:                        false,
	DryRun:                      false,
	UpdateEvents:                false,
	LogFormat:                   "text",
	MetricsAddress:              ":7979",
	LogLevel:                    logrus.InfoLevel.String(),
	ExoscaleEndpoint:            "https://api.exoscale.ch/dns",
	ExoscaleAPIKey:              "",
	ExoscaleAPISecret:           "",
	CRDSourceAPIVersion:         "externaldns.k8s.io/v1alpha1",
	CRDSourceKind:               "DNSEndpoint",
	ServiceTypeFilter:           []string{},
	CFAPIEndpoint:               "",
	CFUsername:                  "",
	CFPassword:                  "",
	RFC2136Host:                 "",
	RFC2136Port:                 0,
	RFC2136Zone:                 "",
	RFC2136Insecure:             false,
	RFC2136GSSTSIG:              false,
	RFC2136KerberosRealm:        "",
	RFC2136KerberosUsername:     "",
	RFC2136KerberosPassword:     "",
	RFC2136TSIGKeyName:          "",
	RFC2136TSIGSecret:           "",
	RFC2136TSIGSecretAlg:        "",
	RFC2136TAXFR:                true,
	RFC2136MinTTL:               0,
	RFC2136BatchChangeSize:      50,
	NS1Endpoint:                 "",
	NS1IgnoreSSL:                false,
	TransIPAccountName:          "",
	TransIPPrivateKeyFile:       "",
	DigitalOceanAPIPageSize:     50,
	ManagedDNSRecordTypes:       []string{endpoint.RecordTypeA, endpoint.RecordTypeCNAME},
	GoDaddyAPIKey:               "",
	GoDaddySecretKey:            "",
	GoDaddyTTL:                  600,
	GoDaddyOTE:                  false,
	IBMCloudProxied:             false,
	IBMCloudConfigFile:          "/etc/kubernetes/ibmcloud.json",
}

func ConvertCRDtoCfg(crd externaldnsv1alpha1.ExternalDNS) externaldns.Config {
	cfg := externaldns.Config{}

	// basic fields are given for testing purpose

	//

	return cfg
}
