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

package plan

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	api "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	sd "github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	log "github.com/sirupsen/logrus"
	"gomodules.xyz/sets"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/pkg/apis/externaldns"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
	"sigs.k8s.io/external-dns/provider/akamai"
	"sigs.k8s.io/external-dns/provider/alibabacloud"
	"sigs.k8s.io/external-dns/provider/aws"
	"sigs.k8s.io/external-dns/provider/awssd"
	"sigs.k8s.io/external-dns/provider/azure"
	"sigs.k8s.io/external-dns/provider/civo"
	"sigs.k8s.io/external-dns/provider/cloudflare"
	"sigs.k8s.io/external-dns/provider/coredns"
	"sigs.k8s.io/external-dns/provider/digitalocean"
	"sigs.k8s.io/external-dns/provider/dnsimple"
	"sigs.k8s.io/external-dns/provider/exoscale"
	"sigs.k8s.io/external-dns/provider/gandi"
	"sigs.k8s.io/external-dns/provider/godaddy"
	"sigs.k8s.io/external-dns/provider/google"
	"sigs.k8s.io/external-dns/provider/inmemory"
	"sigs.k8s.io/external-dns/provider/linode"
	"sigs.k8s.io/external-dns/provider/ns1"
	"sigs.k8s.io/external-dns/provider/oci"
	"sigs.k8s.io/external-dns/provider/ovh"
	"sigs.k8s.io/external-dns/provider/pdns"
	"sigs.k8s.io/external-dns/provider/pihole"
	"sigs.k8s.io/external-dns/provider/plural"
	"sigs.k8s.io/external-dns/provider/rfc2136"
	"sigs.k8s.io/external-dns/provider/scaleway"
	"sigs.k8s.io/external-dns/provider/transip"
	"sigs.k8s.io/external-dns/provider/webhook"
	"sigs.k8s.io/external-dns/registry"
	"sigs.k8s.io/external-dns/source"
	"sigs.k8s.io/external-dns/source/wrappers"
)

var defaultConfig = externaldns.Config{
	AkamaiAccessToken:           "",
	AkamaiClientSecret:          "",
	AkamaiClientToken:           "",
	AkamaiEdgercPath:            "",
	AkamaiEdgercSection:         "",
	AkamaiServiceConsumerDomain: "",
	AlibabaCloudConfigFile:      "/etc/kubernetes/alibaba-cloud.json",
	AnnotationFilter:            "",
	APIServerURL:                "",
	AWSAPIRetries:               3,
	AWSAssumeRole:               "",
	AWSAssumeRoleExternalID:     "",
	AWSBatchChangeInterval:      time.Second,
	AWSBatchChangeSize:          1000,
	AWSBatchChangeSizeBytes:     32000,
	AWSBatchChangeSizeValues:    1000,
	AWSDynamoDBRegion:           "",
	AWSDynamoDBTable:            "external-dns",
	AWSEvaluateTargetHealth:     true,
	AWSPreferCNAME:              false,
	AWSSDCreateTag:              map[string]string{}, // new
	AWSSDServiceCleanup:         false,
	AWSZoneCacheDuration:        0 * time.Second,
	AWSZoneMatchParent:          false,
	AWSZoneTagFilter:            []string{},
	AWSZoneType:                 "",
	AzureConfigFile:             "/etc/kubernetes/azure.json",
	AzureResourceGroup:          "",
	AzureSubscriptionID:         "",
	AzureZonesCacheDuration:     0 * time.Second,
	AzureMaxRetriesCount:        3,
	CFAPIEndpoint:               "",
	CFPassword:                  "",
	CFUsername:                  "",
	CloudflareCustomHostnamesCertificateAuthority: "none",
	CloudflareCustomHostnames:                     false,
	CloudflareCustomHostnamesMinTLSVersion:        "1.0",
	CloudflareDNSRecordsPerPage:                   100,
	CloudflareProxied:                             false,
	CloudflareRegionalServices:                    false,
	CloudflareRegionKey:                           "earth",

	CombineFQDNAndAnnotation:     false,
	Compatibility:                "",
	ConnectorSourceServer:        "localhost:8080",
	CoreDNSPrefix:                "/skydns/",
	CRDSourceAPIVersion:          "externaldns.k8s.io/v1alpha1",
	CRDSourceKind:                "DNSEndpoint",
	DefaultTargets:               []string{},
	DigitalOceanAPIPageSize:      50,
	DomainFilter:                 []string{},
	DryRun:                       false,
	ExcludeDNSRecordTypes:        []string{},
	ExcludeDomains:               []string{},
	ExcludeTargetNets:            []string{},
	EmitEvents:                   []string{},
	ExcludeUnschedulable:         true,
	ExoscaleAPIEnvironment:       "api",
	ExoscaleAPIKey:               "",
	ExoscaleAPISecret:            "",
	ExoscaleAPIZone:              "ch-gva-2",
	ExposeInternalIPV6:           false,
	FQDNTemplate:                 "",
	GatewayLabelFilter:           "",
	GatewayName:                  "",
	GatewayNamespace:             "",
	GlooNamespaces:               []string{"gloo-system"},
	GoDaddyAPIKey:                "",
	GoDaddyOTE:                   false,
	GoDaddySecretKey:             "",
	GoDaddyTTL:                   600,
	GoogleBatchChangeInterval:    time.Second,
	GoogleBatchChangeSize:        1000,
	GoogleProject:                "",
	GoogleZoneVisibility:         "",
	IgnoreHostnameAnnotation:     false,
	IgnoreIngressRulesSpec:       false,
	IgnoreIngressTLSSpec:         false,
	IngressClassNames:            nil,
	InMemoryZones:                []string{},
	Interval:                     time.Minute,
	KubeConfig:                   "",
	LabelFilter:                  labels.Everything().String(),
	LogFormat:                    "text",
	LogLevel:                     log.InfoLevel.String(),
	ManagedDNSRecordTypes:        []string{endpoint.RecordTypeA, endpoint.RecordTypeAAAA, endpoint.RecordTypeCNAME},
	MetricsAddress:               ":7979",
	MinEventSyncInterval:         5 * time.Second,
	Namespace:                    "",
	NAT64Networks:                []string{},
	NS1Endpoint:                  "",
	NS1IgnoreSSL:                 false,
	OCIConfigFile:                "/etc/kubernetes/oci.yaml",
	OCIZoneCacheDuration:         0 * time.Second,
	OCIZoneScope:                 "GLOBAL",
	Once:                         false,
	OVHApiRateLimit:              20,
	OVHEnableCNAMERelative:       false,
	OVHEndpoint:                  "ovh-eu",
	PDNSAPIKey:                   "",
	PDNSServer:                   "http://localhost:8081",
	PDNSServerID:                 "localhost",
	PDNSSkipTLSVerify:            false,
	PiholeApiVersion:             "5",
	PiholePassword:               "",
	PiholeServer:                 "",
	PiholeTLSInsecureSkipVerify:  false,
	PluralCluster:                "",
	PluralProvider:               "",
	PodSourceDomain:              "",
	Policy:                       "sync",
	Provider:                     "",
	ProviderCacheTime:            0,
	PublishHostIP:                false,
	PublishInternal:              false,
	RegexDomainExclusion:         regexp.MustCompile(""),
	RegexDomainFilter:            regexp.MustCompile(""),
	Registry:                     "txt",
	RequestTimeout:               time.Second * 30,
	RFC2136BatchChangeSize:       50,
	RFC2136GSSTSIG:               false,
	RFC2136Host:                  []string{""},
	RFC2136Insecure:              false,
	RFC2136KerberosPassword:      "",
	RFC2136KerberosRealm:         "",
	RFC2136KerberosUsername:      "",
	RFC2136LoadBalancingStrategy: "disabled",
	RFC2136MinTTL:                0,
	RFC2136Port:                  0,
	RFC2136SkipTLSVerify:         false,
	RFC2136TAXFR:                 true,
	RFC2136TSIGKeyName:           "",
	RFC2136TSIGSecret:            "",
	RFC2136TSIGSecretAlg:         "",
	RFC2136UseTLS:                false,
	RFC2136Zone:                  []string{},
	ServiceTypeFilter:            []string{},
	SkipperRouteGroupVersion:     "zalando.org/v1",
	Sources:                      nil,
	TargetNetFilter:              []string{},
	TLSCA:                        "",
	TLSClientCert:                "",
	TLSClientCertKey:             "",
	TraefikEnableLegacy:          false,
	TraefikDisableNew:            false,
	TransIPAccountName:           "",
	TransIPPrivateKeyFile:        "",
	TXTCacheInterval:             0,
	TXTEncryptAESKey:             "",
	TXTEncryptEnabled:            false,
	TXTOwnerID:                   "default",
	TXTPrefix:                    "",
	TXTSuffix:                    "",
	TXTWildcardReplacement:       "",
	UpdateEvents:                 false,
	WebhookProviderReadTimeout:   5 * time.Second,
	WebhookProviderURL:           "http://localhost:8888",
	WebhookProviderWriteTimeout:  10 * time.Second,
	WebhookServer:                false,
	ZoneIDFilter:                 []string{},
	ForceDefaultTargets:          false,
}

func SetDNSRecords(ctx context.Context, edns *api.ExternalDNS) ([]api.DNSRecord, error) {
	cfg := convertEDNSObjectToCfg(edns)

	endpointsSource, err := createEndpointsSource(ctx, cfg)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	domainFilter := createDomainFilter(cfg)

	pvdr, err := buildProvider(ctx, cfg, domainFilter)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	reg, err := createRegistry(cfg, pvdr)
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}

	dnsRecs, err := createAndApplyPlan(ctx, cfg, reg, endpointsSource, domainFilter)
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}

	return dnsRecs, nil
}

// RegexDomainFilter overrides DomainFilter
func createDomainFilter(cfg *externaldns.Config) *endpoint.DomainFilter {
	if cfg.RegexDomainFilter != nil && cfg.RegexDomainFilter.String() != "" {
		return endpoint.NewRegexDomainFilter(cfg.RegexDomainFilter, cfg.RegexDomainExclusion)
	} else {
		return endpoint.NewDomainFilterWithExclusions(cfg.DomainFilter, cfg.ExcludeDomains)
	}
}

// create and apply dns plan, If plan is successfully applied then returns dns record, which defines the desired records of the plan
func createAndApplyPlan(ctx context.Context, cfg *externaldns.Config, r registry.Registry, endpointSource source.Source, domainFilter *endpoint.DomainFilter) ([]api.DNSRecord, error) {

	records, err := r.Records(ctx)
	if err != nil {
		return nil, errors.New("failed to list records, " + err.Error())
	}

	ctx = context.WithValue(ctx, provider.RecordsContextKey, records)
	endpoints, err := endpointSource.Endpoints(ctx)
	if err != nil {
		return nil, errors.New("failed to list source endpoints, " + err.Error())
	}

	endpoints, err = r.AdjustEndpoints(endpoints)
	if err != nil {
		return nil, errors.New("failed to adjust source endpoints, " + err.Error())
	}

	registryFilter := r.GetDomainFilter()

	pln := &plan.Plan{
		Policies:       []plan.Policy{plan.Policies[cfg.Policy]},
		Current:        records,
		Desired:        endpoints,
		DomainFilter:   endpoint.MatchAllDomainFilters{domainFilter, registryFilter},
		ManagedRecords: cfg.ManagedDNSRecordTypes,
	}

	pln = pln.Calculate()
	klog.Info("Desired: ", pln.Desired)
	klog.Info("Current: ", pln.Current)

	dnsRecs := make([]api.DNSRecord, 0)

	if pln.Changes.HasChanges() {
		err = r.ApplyChanges(ctx, pln.Changes)
		if err != nil {
			klog.Error(err.Error())
			return nil, err
		}
		klog.Info("plan applied")
	} else {
		klog.Info("all records are already up to date")
	}

	managedRecordsTypes := sets.NewString()

	for _, dnsType := range cfg.ManagedDNSRecordTypes {
		managedRecordsTypes.Insert(dnsType)
	}

	for _, rec := range pln.Desired {
		if managedRecordsTypes.Has(rec.RecordType) {
			dnsRecs = append(dnsRecs, api.DNSRecord{Name: rec.DNSName, Target: rec.Targets.String()})
		}
	}
	return dnsRecs, nil
}

func convertEDNSObjectToCfg(edns *api.ExternalDNS) *externaldns.Config {
	config := defaultConfig

	if edns.Namespace != "" {
		config.Namespace = edns.Namespace
	}

	if edns.Spec.RequestTimeout != nil {
		config.RequestTimeout = *edns.Spec.RequestTimeout
	}

	// SOURCE
	var sources []string
	sources = append(sources, strings.ToLower(edns.Spec.Source.Type.Kind))
	// sources[] must contain strings that are lower cased
	config.Sources = sources

	if edns.Spec.OCRouterName != nil {
		config.OCPRouterName = *edns.Spec.OCRouterName
	}
	if edns.Spec.GatewayNamespace != nil {
		config.GatewayNamespace = *edns.Spec.GatewayNamespace
	}
	if edns.Spec.GatewayLabelFilter != nil {
		config.GatewayLabelFilter = *edns.Spec.GatewayLabelFilter
	}
	if edns.Spec.ManageDNSRecordTypes != nil {
		config.ManagedDNSRecordTypes = edns.Spec.ManageDNSRecordTypes
	}
	if edns.Spec.DefaultTargets != nil {
		config.DefaultTargets = edns.Spec.DefaultTargets
	}

	// For Node
	if edns.Spec.Source.Node != nil && edns.Spec.Source.Type.Kind == "Node" {
		config.FQDNTemplate = edns.Spec.Source.Node.FQDNTemplate
		if edns.Spec.Source.Node.AnnotationFilter != nil {
			config.AnnotationFilter = *edns.Spec.Source.Node.AnnotationFilter
		}
		if edns.Spec.Source.Node.LabelFilter != nil {
			config.LabelFilter = *edns.Spec.Source.Node.LabelFilter
		}
	}

	// For Service
	if edns.Spec.Source.Service != nil && edns.Spec.Source.Type.Kind == "Service" {
		if edns.Spec.Source.Service.LabelFilter != nil {
			config.LabelFilter = *edns.Spec.Source.Service.LabelFilter
		}
		if edns.Spec.Source.Service.Namespace != nil {
			config.Namespace = *edns.Spec.Source.Service.Namespace
		}
		if edns.Spec.Source.Service.AnnotationFilter != nil {
			config.AnnotationFilter = *edns.Spec.Source.Service.AnnotationFilter
		}
		if edns.Spec.Source.Service.FQDNTemplate != nil {
			config.FQDNTemplate = *edns.Spec.Source.Service.FQDNTemplate
		}
		if edns.Spec.Source.Service.CombineFQDNAndAnnotation != nil {
			config.CombineFQDNAndAnnotation = *edns.Spec.Source.Service.CombineFQDNAndAnnotation
		}
		if edns.Spec.Source.Service.Compatibility != nil {
			config.Compatibility = *edns.Spec.Source.Service.Compatibility
		}
		if edns.Spec.Source.Service.PublishInternal != nil {
			config.PublishInternal = *edns.Spec.Source.Service.PublishInternal
		}
		if edns.Spec.Source.Service.PublishHostIP != nil {
			config.PublishHostIP = *edns.Spec.Source.Service.PublishHostIP
		}
		if edns.Spec.Source.Service.AlwaysPublishNotReadyAddresses != nil {
			config.AlwaysPublishNotReadyAddresses = *edns.Spec.Source.Service.AlwaysPublishNotReadyAddresses
		}
		if edns.Spec.Source.Service.ServiceTypeFilter != nil {
			config.ServiceTypeFilter = edns.Spec.Source.Service.ServiceTypeFilter
		}
		if edns.Spec.Source.Service.IgnoreHostnameAnnotation != nil {
			config.IgnoreHostnameAnnotation = *edns.Spec.Source.Service.IgnoreHostnameAnnotation
		}
	}

	// For Ingress
	if edns.Spec.Source.Ingress != nil && edns.Spec.Source.Type.Kind == "Ingress" {
		if edns.Spec.Source.Ingress.IgnoreIngressRulesSpec != nil {
			config.IgnoreIngressRulesSpec = *edns.Spec.Source.Ingress.IgnoreIngressRulesSpec
		}
		if edns.Spec.Source.Ingress.IgnoreHostnameAnnotation != nil {
			config.IgnoreHostnameAnnotation = *edns.Spec.Source.Ingress.IgnoreHostnameAnnotation
		}
		if edns.Spec.Source.Ingress.FQDNTemplate != nil {
			config.FQDNTemplate = *edns.Spec.Source.Ingress.FQDNTemplate
		}
		if edns.Spec.Source.Ingress.Namespace != nil {
			config.Namespace = *edns.Spec.Source.Ingress.Namespace
		}
		if edns.Spec.Source.Ingress.AnnotationFilter != nil {
			config.AnnotationFilter = *edns.Spec.Source.Ingress.AnnotationFilter
		}
		if edns.Spec.Source.Ingress.CombineFQDNAndAnnotation != nil {
			config.CombineFQDNAndAnnotation = *edns.Spec.Source.Ingress.CombineFQDNAndAnnotation
		}
		if edns.Spec.Source.Ingress.IgnoreIngressTLSSpec != nil {
			config.IgnoreIngressTLSSpec = *edns.Spec.Source.Ingress.IgnoreIngressTLSSpec
		}
		if edns.Spec.Source.Ingress.LabelFilter != nil {
			config.LabelFilter = *edns.Spec.Source.Ingress.LabelFilter
		}
	}

	// PROVIDER
	config.Provider = edns.Spec.Provider.String()

	if edns.Spec.DomainFilter != nil {
		config.DomainFilter = edns.Spec.DomainFilter
	}
	if edns.Spec.ExcludeDomains != nil {
		config.ExcludeDomains = edns.Spec.ExcludeDomains
	}
	if edns.Spec.ZoneIDFilter != nil {
		config.ZoneIDFilter = edns.Spec.ZoneIDFilter
	}

	// for aws provider
	if edns.Spec.AWS != nil {
		if edns.Spec.AWS.ZoneTagFilter != nil {
			config.AWSZoneTagFilter = edns.Spec.AWS.ZoneTagFilter
		}
		if edns.Spec.AWS.ZoneType != nil {
			config.AWSZoneType = *edns.Spec.AWS.ZoneType
		}
		if edns.Spec.AWS.AssumeRole != nil {
			config.AWSAssumeRole = *edns.Spec.AWS.AssumeRole
		}
		if edns.Spec.AWS.BatchChangeSize != nil {
			config.AWSBatchChangeSize = *edns.Spec.AWS.BatchChangeSize
		}
		if edns.Spec.AWS.BatchChangeInterval != nil {
			config.AWSBatchChangeInterval = *edns.Spec.AWS.BatchChangeInterval
		}
		if edns.Spec.AWS.EvaluateTargetHealth != nil {
			config.AWSEvaluateTargetHealth = *edns.Spec.AWS.EvaluateTargetHealth
		}
		if edns.Spec.AWS.APIRetries != nil {
			config.AWSAPIRetries = *edns.Spec.AWS.APIRetries
		}
		if edns.Spec.AWS.PreferCNAME != nil {
			config.AWSPreferCNAME = *edns.Spec.AWS.PreferCNAME
		}
		if edns.Spec.AWS.ZoneCacheDuration != nil {
			config.AWSZoneCacheDuration = *edns.Spec.AWS.ZoneCacheDuration
		}
		if edns.Spec.AWS.SDServiceCleanup != nil {
			config.AWSSDServiceCleanup = *edns.Spec.AWS.SDServiceCleanup
		}
		if edns.Spec.AWS.SDCreateTag != nil {
			config.AWSSDCreateTag = *edns.Spec.AWS.SDCreateTag
		}
	}

	// for cloudflare provider
	if edns.Spec.Cloudflare != nil {
		if edns.Spec.Cloudflare.Proxied != nil {
			config.CloudflareProxied = *edns.Spec.Cloudflare.Proxied
		}
		if edns.Spec.Cloudflare.CustomHostnames != nil {
			config.CloudflareCustomHostnames = *edns.Spec.Cloudflare.CustomHostnames
		}
		if edns.Spec.Cloudflare.CustomHostnamesCertificateAuthority != nil {
			config.CloudflareCustomHostnamesCertificateAuthority = *edns.Spec.Cloudflare.CustomHostnamesCertificateAuthority
		}
		if edns.Spec.Cloudflare.CustomHostnamesMinTLSVersion != nil {
			config.CloudflareCustomHostnamesMinTLSVersion = *edns.Spec.Cloudflare.CustomHostnamesMinTLSVersion
		}
		if edns.Spec.Cloudflare.RegionalServices != nil {
			config.CloudflareRegionalServices = *edns.Spec.Cloudflare.RegionalServices
		}
		if edns.Spec.Cloudflare.RegionKey != nil {
			config.CloudflareRegionKey = *edns.Spec.Cloudflare.RegionKey
		}
	}

	// for azure provide
	if edns.Spec.Provider == api.ProviderAzure {
		// hard-code assignment of AzureConfigFile path
		config.AzureConfigFile = fmt.Sprintf("/tmp/%s-%s-credential", edns.Namespace, edns.Name)
	}
	if edns.Spec.Azure != nil {
		if edns.Spec.Azure.SubscriptionId != nil {
			config.AzureSubscriptionID = *edns.Spec.Azure.SubscriptionId
		}
		if edns.Spec.Azure.ResourceGroup != nil {
			config.AzureResourceGroup = *edns.Spec.Azure.ResourceGroup
		}
		if edns.Spec.Azure.UserAssignedIdentityClientID != nil {
			config.AzureUserAssignedIdentityClientID = *edns.Spec.Azure.UserAssignedIdentityClientID
		}
		if edns.Spec.Azure.ZonesCacheDuration != nil {
			config.AzureZonesCacheDuration = *edns.Spec.Azure.ZonesCacheDuration
		}
		if edns.Spec.Azure.MaxRetriesCount != nil {
			config.AzureMaxRetriesCount = *edns.Spec.Azure.MaxRetriesCount
		}
	}

	// for google dns provider
	if edns.Spec.Google != nil {
		if edns.Spec.Google.Project != nil {
			config.GoogleProject = *edns.Spec.Google.Project
		}

		if edns.Spec.Google.BatchChangeSize != nil {
			config.GoogleBatchChangeSize = *edns.Spec.Google.BatchChangeSize
		}

		if edns.Spec.Google.BatchChangeInterval != nil {
			config.GoogleBatchChangeInterval = *edns.Spec.Google.BatchChangeInterval
		}

		if edns.Spec.Google.ZoneVisibility != nil {
			config.GoogleZoneVisibility = *edns.Spec.Google.ZoneVisibility
		}
	}

	// POLICY

	if edns.Spec.Policy != nil {
		config.Policy = edns.Spec.Policy.String()
	}

	// REGISTRY
	if edns.Spec.Registry != nil {
		config.Registry = *edns.Spec.Registry
	}
	if edns.Spec.TXTOwnerID != nil {
		config.TXTOwnerID = *edns.Spec.TXTOwnerID
	}
	if edns.Spec.TXTPrefix != nil {
		config.TXTPrefix = *edns.Spec.TXTPrefix
	}
	if edns.Spec.TXTSuffix != nil {
		config.TXTSuffix = *edns.Spec.TXTSuffix
	}
	if edns.Spec.TXTWildcardReplacement != nil {
		config.TXTWildcardReplacement = *edns.Spec.TXTWildcardReplacement
	}

	return &config
}

func createEndpointsSource(ctx context.Context, cfg *externaldns.Config) (source.Source, error) {
	sourceCfg := source.NewSourceConfig(cfg)

	// Lookup all the selected sources by names and pass them the desired configuration.
	sources, err := source.ByNames(ctx, &source.SingletonClientGenerator{
		KubeConfig:   cfg.KubeConfig,
		APIServerURL: cfg.APIServerURL,
		RequestTimeout: func() time.Duration {
			if cfg.UpdateEvents {
				return 0
			}
			return cfg.RequestTimeout
		}(),
	}, cfg.Sources, sourceCfg)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	combinedSource := wrappers.NewDedupSource(wrappers.NewMultiSource(sources, sourceCfg.DefaultTargets, sourceCfg.ForceDefaultTargets))
	cfg.AddSourceWrapper("dedup")
	combinedSource = wrappers.NewNAT64Source(combinedSource, cfg.NAT64Networks)
	cfg.AddSourceWrapper("nat64")
	// Filter targets
	targetFilter := endpoint.NewTargetNetFilterWithExclusions(cfg.TargetNetFilter, cfg.ExcludeTargetNets)
	if targetFilter.IsEnabled() {
		combinedSource = wrappers.NewTargetFilterSource(combinedSource, targetFilter)
		cfg.AddSourceWrapper("target-filter")
	}

	return combinedSource, nil
}

func buildProvider(
	ctx context.Context,
	cfg *externaldns.Config,
	domainFilter *endpoint.DomainFilter,
) (provider.Provider, error) {
	var p provider.Provider
	var err error

	zoneNameFilter := endpoint.NewDomainFilter(cfg.ZoneNameFilter)
	zoneIDFilter := provider.NewZoneIDFilter(cfg.ZoneIDFilter)
	zoneTypeFilter := provider.NewZoneTypeFilter(cfg.AWSZoneType)
	zoneTagFilter := provider.NewZoneTagFilter(cfg.AWSZoneTagFilter)

	switch cfg.Provider {
	case "akamai":
		p, err = akamai.NewAkamaiProvider(
			akamai.AkamaiConfig{
				DomainFilter:          domainFilter,
				ZoneIDFilter:          zoneIDFilter,
				ServiceConsumerDomain: cfg.AkamaiServiceConsumerDomain,
				ClientToken:           cfg.AkamaiClientToken,
				ClientSecret:          cfg.AkamaiClientSecret,
				AccessToken:           cfg.AkamaiAccessToken,
				EdgercPath:            cfg.AkamaiEdgercPath,
				EdgercSection:         cfg.AkamaiEdgercSection,
				DryRun:                cfg.DryRun,
			}, nil)
	case "alibabacloud":
		p, err = alibabacloud.NewAlibabaCloudProvider(cfg.AlibabaCloudConfigFile, domainFilter, zoneIDFilter, cfg.AlibabaCloudZoneType, cfg.DryRun)
	case "aws":
		configs := aws.CreateV2Configs(cfg)
		clients := make(map[string]aws.Route53API, len(configs))
		for profile, config := range configs {
			clients[profile] = route53.NewFromConfig(config)
		}

		p, err = aws.NewAWSProvider(
			aws.AWSConfig{
				DomainFilter:          domainFilter,
				ZoneIDFilter:          zoneIDFilter,
				ZoneTypeFilter:        zoneTypeFilter,
				ZoneTagFilter:         zoneTagFilter,
				ZoneMatchParent:       cfg.AWSZoneMatchParent,
				BatchChangeSize:       cfg.AWSBatchChangeSize,
				BatchChangeSizeBytes:  cfg.AWSBatchChangeSizeBytes,
				BatchChangeSizeValues: cfg.AWSBatchChangeSizeValues,
				BatchChangeInterval:   cfg.AWSBatchChangeInterval,
				EvaluateTargetHealth:  cfg.AWSEvaluateTargetHealth,
				PreferCNAME:           cfg.AWSPreferCNAME,
				DryRun:                cfg.DryRun,
				ZoneCacheDuration:     cfg.AWSZoneCacheDuration,
			},
			clients,
		)
	case "aws-sd":
		// Check that only compatible Registry is used with AWS-SD
		if cfg.Registry != "noop" && cfg.Registry != "aws-sd" {
			log.Infof("Registry \"%s\" cannot be used with AWS Cloud Map. Switching to \"aws-sd\".", cfg.Registry)
			cfg.Registry = "aws-sd"
		}
		p, err = awssd.NewAWSSDProvider(domainFilter, cfg.AWSZoneType, cfg.DryRun, cfg.AWSSDServiceCleanup, cfg.TXTOwnerID, cfg.AWSSDCreateTag, sd.NewFromConfig(aws.CreateDefaultV2Config(cfg)))
	case "azure-dns", "azure":
		p, err = azure.NewAzureProvider(cfg.AzureConfigFile, domainFilter, zoneNameFilter, zoneIDFilter, cfg.AzureSubscriptionID, cfg.AzureResourceGroup, cfg.AzureUserAssignedIdentityClientID, cfg.AzureActiveDirectoryAuthorityHost, cfg.AzureZonesCacheDuration, cfg.AzureMaxRetriesCount, cfg.DryRun)
	case "azure-private-dns":
		p, err = azure.NewAzurePrivateDNSProvider(cfg.AzureConfigFile, domainFilter, zoneNameFilter, zoneIDFilter, cfg.AzureSubscriptionID, cfg.AzureResourceGroup, cfg.AzureUserAssignedIdentityClientID, cfg.AzureActiveDirectoryAuthorityHost, cfg.AzureZonesCacheDuration, cfg.AzureMaxRetriesCount, cfg.DryRun)
	case "civo":
		p, err = civo.NewCivoProvider(domainFilter, cfg.DryRun)
	case "cloudflare":
		p, err = cloudflare.NewCloudFlareProvider(
			domainFilter,
			zoneIDFilter,
			cfg.CloudflareProxied,
			cfg.DryRun,
			cloudflare.RegionalServicesConfig{
				Enabled:   cfg.CloudflareRegionalServices,
				RegionKey: cfg.CloudflareRegionKey,
			},
			cloudflare.CustomHostnamesConfig{
				Enabled:              cfg.CloudflareCustomHostnames,
				MinTLSVersion:        cfg.CloudflareCustomHostnamesMinTLSVersion,
				CertificateAuthority: cfg.CloudflareCustomHostnamesCertificateAuthority,
			},
			cloudflare.DNSRecordsConfig{
				PerPage: cfg.CloudflareDNSRecordsPerPage,
				Comment: cfg.CloudflareDNSRecordsComment,
			})
	case "google":
		p, err = google.NewGoogleProvider(ctx, cfg.GoogleProject, domainFilter, zoneIDFilter, cfg.GoogleBatchChangeSize, cfg.GoogleBatchChangeInterval, cfg.GoogleZoneVisibility, cfg.DryRun)
	case "digitalocean":
		p, err = digitalocean.NewDigitalOceanProvider(ctx, domainFilter, cfg.DryRun, cfg.DigitalOceanAPIPageSize)
	case "ovh":
		p, err = ovh.NewOVHProvider(ctx, domainFilter, cfg.OVHEndpoint, cfg.OVHApiRateLimit, cfg.OVHEnableCNAMERelative, cfg.DryRun)
	case "linode":
		p, err = linode.NewLinodeProvider(domainFilter, cfg.DryRun)
	case "dnsimple":
		p, err = dnsimple.NewDnsimpleProvider(domainFilter, zoneIDFilter, cfg.DryRun)
	case "coredns", "skydns":
		p, err = coredns.NewCoreDNSProvider(domainFilter, cfg.CoreDNSPrefix, cfg.DryRun)
	case "exoscale":
		p, err = exoscale.NewExoscaleProvider(
			cfg.ExoscaleAPIEnvironment,
			cfg.ExoscaleAPIZone,
			cfg.ExoscaleAPIKey,
			cfg.ExoscaleAPISecret,
			cfg.DryRun,
			exoscale.ExoscaleWithDomain(domainFilter),
			exoscale.ExoscaleWithLogging(),
		)
	case "inmemory":
		p, err = inmemory.NewInMemoryProvider(inmemory.InMemoryInitZones(cfg.InMemoryZones), inmemory.InMemoryWithDomain(domainFilter), inmemory.InMemoryWithLogging()), nil
	case "pdns":
		p, err = pdns.NewPDNSProvider(
			ctx,
			pdns.PDNSConfig{
				DomainFilter: domainFilter,
				DryRun:       cfg.DryRun,
				Server:       cfg.PDNSServer,
				ServerID:     cfg.PDNSServerID,
				APIKey:       cfg.PDNSAPIKey,
				TLSConfig: pdns.TLSConfig{
					SkipTLSVerify:         cfg.PDNSSkipTLSVerify,
					CAFilePath:            cfg.TLSCA,
					ClientCertFilePath:    cfg.TLSClientCert,
					ClientCertKeyFilePath: cfg.TLSClientCertKey,
				},
			},
		)
	case "oci":
		var config *oci.OCIConfig
		// if the instance-principals flag was set, and a compartment OCID was provided, then ignore the
		// OCI config file, and provide a config that uses instance principal authentication.
		if cfg.OCIAuthInstancePrincipal {
			if len(cfg.OCICompartmentOCID) == 0 {
				err = fmt.Errorf("instance principal authentication requested, but no compartment OCID provided")
			} else {
				authConfig := oci.OCIAuthConfig{UseInstancePrincipal: true}
				config = &oci.OCIConfig{Auth: authConfig, CompartmentID: cfg.OCICompartmentOCID}
			}
		} else {
			config, err = oci.LoadOCIConfig(cfg.OCIConfigFile)
		}
		config.ZoneCacheDuration = cfg.OCIZoneCacheDuration
		if err == nil {
			p, err = oci.NewOCIProvider(*config, domainFilter, zoneIDFilter, cfg.OCIZoneScope, cfg.DryRun)
		}
	case "rfc2136":
		tlsConfig := rfc2136.TLSConfig{
			UseTLS:                cfg.RFC2136UseTLS,
			SkipTLSVerify:         cfg.RFC2136SkipTLSVerify,
			CAFilePath:            cfg.TLSCA,
			ClientCertFilePath:    cfg.TLSClientCert,
			ClientCertKeyFilePath: cfg.TLSClientCertKey,
		}
		p, err = rfc2136.NewRfc2136Provider(cfg.RFC2136Host, cfg.RFC2136Port, cfg.RFC2136Zone, cfg.RFC2136Insecure, cfg.RFC2136TSIGKeyName, cfg.RFC2136TSIGSecret, cfg.RFC2136TSIGSecretAlg, cfg.RFC2136TAXFR, domainFilter, cfg.DryRun, cfg.RFC2136MinTTL, cfg.RFC2136CreatePTR, cfg.RFC2136GSSTSIG, cfg.RFC2136KerberosUsername, cfg.RFC2136KerberosPassword, cfg.RFC2136KerberosRealm, cfg.RFC2136BatchChangeSize, tlsConfig, cfg.RFC2136LoadBalancingStrategy, nil)
	case "ns1":
		p, err = ns1.NewNS1Provider(
			ns1.NS1Config{
				DomainFilter:  domainFilter,
				ZoneIDFilter:  zoneIDFilter,
				NS1Endpoint:   cfg.NS1Endpoint,
				NS1IgnoreSSL:  cfg.NS1IgnoreSSL,
				DryRun:        cfg.DryRun,
				MinTTLSeconds: cfg.NS1MinTTLSeconds,
			},
		)
	case "transip":
		p, err = transip.NewTransIPProvider(cfg.TransIPAccountName, cfg.TransIPPrivateKeyFile, domainFilter, cfg.DryRun)
	case "scaleway":
		p, err = scaleway.NewScalewayProvider(ctx, domainFilter, cfg.DryRun)
	case "godaddy":
		p, err = godaddy.NewGoDaddyProvider(ctx, domainFilter, cfg.GoDaddyTTL, cfg.GoDaddyAPIKey, cfg.GoDaddySecretKey, cfg.GoDaddyOTE, cfg.DryRun)
	case "gandi":
		p, err = gandi.NewGandiProvider(ctx, domainFilter, cfg.DryRun)
	case "pihole":
		p, err = pihole.NewPiholeProvider(
			pihole.PiholeConfig{
				Server:                cfg.PiholeServer,
				Password:              cfg.PiholePassword,
				TLSInsecureSkipVerify: cfg.PiholeTLSInsecureSkipVerify,
				DomainFilter:          domainFilter,
				DryRun:                cfg.DryRun,
				APIVersion:            cfg.PiholeApiVersion,
			},
		)
	case "plural":
		p, err = plural.NewPluralProvider(cfg.PluralCluster, cfg.PluralProvider)
	case "webhook":
		p, err = webhook.NewWebhookProvider(cfg.WebhookProviderURL)
	default:
		err = fmt.Errorf("unknown dns provider: %s", cfg.Provider)
	}
	if p != nil && cfg.ProviderCacheTime > 0 {
		p = provider.NewCachedProvider(
			p,
			cfg.ProviderCacheTime,
		)
	}
	return p, err
}

func createRegistry(cfg *externaldns.Config, p provider.Provider) (registry.Registry, error) {
	var r registry.Registry
	var err error
	switch cfg.Registry {
	case "dynamodb":
		var dynamodbOpts []func(*dynamodb.Options)
		if cfg.AWSDynamoDBRegion != "" {
			dynamodbOpts = []func(*dynamodb.Options){
				func(opts *dynamodb.Options) {
					opts.Region = cfg.AWSDynamoDBRegion
				},
			}
		}
		r, err = registry.NewDynamoDBRegistry(p, cfg.TXTOwnerID, dynamodb.NewFromConfig(aws.CreateDefaultV2Config(cfg), dynamodbOpts...), cfg.AWSDynamoDBTable, cfg.TXTPrefix, cfg.TXTSuffix, cfg.TXTWildcardReplacement, cfg.ManagedDNSRecordTypes, cfg.ExcludeDNSRecordTypes, []byte(cfg.TXTEncryptAESKey), cfg.TXTCacheInterval)
	case "noop":
		r, err = registry.NewNoopRegistry(p)
	case "txt":
		r, err = registry.NewTXTRegistry(p, cfg.TXTPrefix, cfg.TXTSuffix, cfg.TXTOwnerID, cfg.TXTCacheInterval, cfg.TXTWildcardReplacement, cfg.ManagedDNSRecordTypes, cfg.ExcludeDNSRecordTypes, cfg.TXTEncryptEnabled, []byte(cfg.TXTEncryptAESKey))
	case "aws-sd":
		r, err = registry.NewAWSSDRegistry(p, cfg.TXTOwnerID)
	default:
		log.Fatalf("unknown registry: %s", cfg.Registry)
	}
	return r, err
}
