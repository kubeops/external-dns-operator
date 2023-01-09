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
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external/v1alpha1"

	"github.com/sirupsen/logrus"
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
	"sigs.k8s.io/external-dns/provider/bluecat"
	"sigs.k8s.io/external-dns/provider/cloudflare"
	"sigs.k8s.io/external-dns/provider/coredns"
	"sigs.k8s.io/external-dns/provider/designate"
	"sigs.k8s.io/external-dns/provider/digitalocean"
	"sigs.k8s.io/external-dns/provider/dnsimple"
	"sigs.k8s.io/external-dns/provider/dyn"
	"sigs.k8s.io/external-dns/provider/exoscale"
	"sigs.k8s.io/external-dns/provider/gandi"
	"sigs.k8s.io/external-dns/provider/godaddy"
	"sigs.k8s.io/external-dns/provider/google"
	"sigs.k8s.io/external-dns/provider/ibmcloud"
	"sigs.k8s.io/external-dns/provider/infoblox"
	"sigs.k8s.io/external-dns/provider/inmemory"
	"sigs.k8s.io/external-dns/provider/linode"
	"sigs.k8s.io/external-dns/provider/ns1"
	"sigs.k8s.io/external-dns/provider/oci"
	"sigs.k8s.io/external-dns/provider/ovh"
	"sigs.k8s.io/external-dns/provider/pdns"
	"sigs.k8s.io/external-dns/provider/rcode0"
	"sigs.k8s.io/external-dns/provider/rdns"
	"sigs.k8s.io/external-dns/provider/rfc2136"
	"sigs.k8s.io/external-dns/provider/safedns"
	"sigs.k8s.io/external-dns/provider/scaleway"
	"sigs.k8s.io/external-dns/provider/transip"
	"sigs.k8s.io/external-dns/provider/ultradns"
	"sigs.k8s.io/external-dns/provider/vinyldns"
	"sigs.k8s.io/external-dns/provider/vultr"
	"sigs.k8s.io/external-dns/registry"
	"sigs.k8s.io/external-dns/source"
)

const (
	providerAWSSD = "aws-sd"
)

var defaultConfig = externaldns.Config{
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
	AlibabaCloudConfigFile:      "/etc/kubernetes/alibaba-cloud.json",
	AWSZoneType:                 "",
	AWSZoneTagFilter:            []string{},
	AWSAssumeRole:               "",
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

func SetDNSRecords(ctx context.Context, edns *externaldnsv1alpha1.ExternalDNS) ([]externaldnsv1alpha1.DNSRecord, error) {
	cfg := convertEDNSObjectToCfg(edns)

	endpointsSource, err := createEndpointsSource(ctx, cfg)
	if err != nil {
		klog.Error("failed to create endpoints source.", err.Error())
		return nil, err
	}

	pvdr, err := createProviderFromCfg(ctx, cfg, endpointsSource)
	if err != nil {
		klog.Error("failed to create provider: ", err.Error())
		return nil, err
	}

	reg, err := createRegistry(cfg, *pvdr)
	if err != nil {
		klog.Errorf("failed to create Registry.", err.Error())
		return nil, err
	}

	dnsRecs, e := createAndApplyPlan(ctx, cfg, reg, endpointsSource)
	if e != nil {
		klog.Errorf("failed to create and apply plan: %s", err.Error())
		return nil, e
	}

	return dnsRecs, nil
}

// create and apply dns plan, If plan is successfully applied then returns dns record, which defines the desired records of the plan
func createAndApplyPlan(ctx context.Context, cfg *externaldns.Config, r registry.Registry, endpointSource source.Source) ([]externaldnsv1alpha1.DNSRecord, error) {
	var domainFilter endpoint.DomainFilter
	if cfg.RegexDomainFilter.String() != "" {
		domainFilter = endpoint.NewRegexDomainFilter(cfg.RegexDomainFilter, cfg.RegexDomainExclusion)
	} else {
		domainFilter = endpoint.NewDomainFilterWithExclusions(cfg.DomainFilter, cfg.ExcludeDomains)
	}

	records, err := r.Records(ctx)
	if err != nil {
		return nil, err
	}

	missingRecords := r.MissingRecords()

	ctx = context.WithValue(ctx, provider.RecordsContextKey, records)
	endpoints, err := endpointSource.Endpoints(ctx)
	if err != nil {
		return nil, err
	}

	endpoints = r.AdjustEndpoints(endpoints)

	if len(missingRecords) > 0 {
		missingRecordsPlan := &plan.Plan{
			Policies:           []plan.Policy{plan.Policies[cfg.Policy]},
			Missing:            missingRecords,
			DomainFilter:       domainFilter,
			PropertyComparator: r.PropertyValuesEqual,
			ManagedRecords:     cfg.ManagedDNSRecordTypes,
		}

		missingRecordsPlan = missingRecordsPlan.Calculate()
		if missingRecordsPlan.Changes.HasChanges() {
			err = r.ApplyChanges(ctx, missingRecordsPlan.Changes)
			if err != nil {
				return nil, err
			}
			klog.Info("all missing records are created")
		}
	}

	pln := &plan.Plan{
		Policies:           []plan.Policy{plan.Policies[cfg.Policy]},
		Current:            records,
		Desired:            endpoints,
		DomainFilter:       domainFilter,
		PropertyComparator: r.PropertyValuesEqual,
		ManagedRecords:     cfg.ManagedDNSRecordTypes,
	}

	pln = pln.Calculate()
	klog.Info("Desired: ", pln.Desired)
	klog.Info("Current: ", pln.Current)

	dnsRecs := make([]externaldnsv1alpha1.DNSRecord, 0)

	if pln.Changes.HasChanges() {
		err = r.ApplyChanges(ctx, pln.Changes)
		if err != nil {
			klog.Error("failed to apply plan")
			return nil, err
		}
		klog.Info("plan applied")

	} else {
		klog.Info("all records are already up to date")
	}

	for _, rec := range pln.Desired {
		dnsRecs = append(dnsRecs, externaldnsv1alpha1.DNSRecord{Name: rec.DNSName, Target: rec.Targets.String()})
	}
	return dnsRecs, nil
}

func convertEDNSObjectToCfg(edns *externaldnsv1alpha1.ExternalDNS) *externaldns.Config {
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
	}

	// for cloudflare provider
	if edns.Spec.Cloudflare != nil {

		if edns.Spec.Cloudflare.Proxied != nil {
			config.CloudflareProxied = *edns.Spec.Cloudflare.Proxied
		}

		if edns.Spec.Cloudflare.ZonesPerPage != nil {
			config.CloudflareZonesPerPage = *edns.Spec.Cloudflare.ZonesPerPage
		}
	}

	// for azure provide
	if edns.Spec.Provider.String() == externaldnsv1alpha1.ProviderAzure.String() {
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
	labelSelector, err := labels.Parse(cfg.LabelFilter)
	if err != nil {
		return nil, err
	}

	// Create a source.Config from the flags passed by the user.
	sourceCfg := &source.Config{
		Namespace:                      cfg.Namespace,
		AnnotationFilter:               cfg.AnnotationFilter,
		LabelFilter:                    labelSelector,
		FQDNTemplate:                   cfg.FQDNTemplate,
		CombineFQDNAndAnnotation:       cfg.CombineFQDNAndAnnotation,
		IgnoreHostnameAnnotation:       cfg.IgnoreHostnameAnnotation,
		IgnoreIngressTLSSpec:           cfg.IgnoreIngressTLSSpec,
		IgnoreIngressRulesSpec:         cfg.IgnoreIngressRulesSpec,
		GatewayNamespace:               cfg.GatewayNamespace,
		GatewayLabelFilter:             cfg.GatewayLabelFilter,
		Compatibility:                  cfg.Compatibility,
		PublishInternal:                cfg.PublishInternal,
		PublishHostIP:                  cfg.PublishHostIP,
		AlwaysPublishNotReadyAddresses: cfg.AlwaysPublishNotReadyAddresses,
		ConnectorServer:                cfg.ConnectorSourceServer,
		CRDSourceAPIVersion:            cfg.CRDSourceAPIVersion,
		CRDSourceKind:                  cfg.CRDSourceKind,
		KubeConfig:                     cfg.KubeConfig,
		APIServerURL:                   cfg.APIServerURL,
		ServiceTypeFilter:              cfg.ServiceTypeFilter,
		CFAPIEndpoint:                  cfg.CFAPIEndpoint,
		CFUsername:                     cfg.CFUsername,
		CFPassword:                     cfg.CFPassword,
		ContourLoadBalancerService:     cfg.ContourLoadBalancerService,
		GlooNamespace:                  cfg.GlooNamespace,
		SkipperRouteGroupVersion:       cfg.SkipperRouteGroupVersion,
		RequestTimeout:                 cfg.RequestTimeout,
		DefaultTargets:                 cfg.DefaultTargets,
		OCPRouterName:                  cfg.OCPRouterName,
	}

	// Lookup all the selected sources by names and pass them the desired configuration.
	sources, err := source.ByNames(ctx, &source.SingletonClientGenerator{
		KubeConfig:   cfg.KubeConfig,
		APIServerURL: cfg.APIServerURL,
		// If update events are enabled, disable timeout.
		RequestTimeout: func() time.Duration {
			if cfg.UpdateEvents {
				return 0
			}
			return cfg.RequestTimeout
		}(),
	}, cfg.Sources, sourceCfg)
	if err != nil {
		klog.Error("failed to get the source")
		return nil, err
	}

	// Combine multiple sources into a single, deduplicated source.
	endpointsSource := source.NewDedupSource(source.NewMultiSource(sources, sourceCfg.DefaultTargets))

	return endpointsSource, nil
}

func createProviderFromCfg(ctx context.Context, cfg *externaldns.Config, endpointsSource source.Source) (*provider.Provider, error) {
	var p provider.Provider
	var err error

	var domainFilter endpoint.DomainFilter
	if cfg.RegexDomainFilter.String() != "" {
		domainFilter = endpoint.NewRegexDomainFilter(cfg.RegexDomainFilter, cfg.RegexDomainExclusion)
	} else {
		domainFilter = endpoint.NewDomainFilterWithExclusions(cfg.DomainFilter, cfg.ExcludeDomains)
	}

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
		p, err = aws.NewAWSProvider(
			aws.AWSConfig{
				DomainFilter:         domainFilter,
				ZoneIDFilter:         zoneIDFilter,
				ZoneTypeFilter:       zoneTypeFilter,
				ZoneTagFilter:        zoneTagFilter,
				BatchChangeSize:      cfg.AWSBatchChangeSize,
				BatchChangeInterval:  cfg.AWSBatchChangeInterval,
				EvaluateTargetHealth: cfg.AWSEvaluateTargetHealth,
				AssumeRole:           cfg.AWSAssumeRole,
				APIRetries:           cfg.AWSAPIRetries,
				PreferCNAME:          cfg.AWSPreferCNAME,
				DryRun:               cfg.DryRun,
				ZoneCacheDuration:    cfg.AWSZoneCacheDuration,
			},
		)
	case providerAWSSD:
		if cfg.Registry != "noop" && cfg.Registry != providerAWSSD {
			cfg.Registry = providerAWSSD
		}
		p, err = awssd.NewAWSSDProvider(domainFilter, cfg.AWSZoneType, cfg.AWSAssumeRole, cfg.AWSAssumeRoleExternalID, cfg.DryRun, cfg.AWSSDServiceCleanup, cfg.TXTOwnerID)
	case "azure-dns", "azure":
		p, err = azure.NewAzureProvider(cfg.AzureConfigFile, domainFilter, zoneNameFilter, zoneIDFilter, cfg.AzureResourceGroup, cfg.AzureUserAssignedIdentityClientID, cfg.DryRun)
	case "azure-private-dns":
		p, err = azure.NewAzurePrivateDNSProvider(cfg.AzureConfigFile, domainFilter, zoneIDFilter, cfg.AzureResourceGroup, cfg.AzureUserAssignedIdentityClientID, cfg.DryRun)
	case "bluecat":
		p, err = bluecat.NewBluecatProvider(cfg.BluecatConfigFile, cfg.BluecatDNSConfiguration, cfg.BluecatDNSServerName, cfg.BluecatDNSDeployType, cfg.BluecatDNSView, cfg.BluecatGatewayHost, cfg.BluecatRootZone, cfg.TXTPrefix, cfg.TXTSuffix, domainFilter, zoneIDFilter, cfg.DryRun, cfg.BluecatSkipTLSVerify)
	case "vinyldns":
		p, err = vinyldns.NewVinylDNSProvider(domainFilter, zoneIDFilter, cfg.DryRun)
	case "vultr":
		p, err = vultr.NewVultrProvider(ctx, domainFilter, cfg.DryRun)
	case "ultradns":
		p, err = ultradns.NewUltraDNSProvider(domainFilter, cfg.DryRun)
	case "cloudflare":
		p, err = cloudflare.NewCloudFlareProvider(domainFilter, zoneIDFilter, cfg.CloudflareZonesPerPage, cfg.CloudflareProxied, cfg.DryRun)
	case "rcodezero":
		p, err = rcode0.NewRcodeZeroProvider(domainFilter, cfg.DryRun, cfg.RcodezeroTXTEncrypt)
	case "google":
		p, err = google.NewGoogleProvider(ctx, cfg.GoogleProject, domainFilter, zoneIDFilter, cfg.GoogleBatchChangeSize, cfg.GoogleBatchChangeInterval, cfg.GoogleZoneVisibility, cfg.DryRun)
	case "digitalocean":
		p, err = digitalocean.NewDigitalOceanProvider(ctx, domainFilter, cfg.DryRun, cfg.DigitalOceanAPIPageSize)
	case "ovh":
		p, err = ovh.NewOVHProvider(ctx, domainFilter, cfg.OVHEndpoint, cfg.OVHApiRateLimit, cfg.DryRun)
	case "linode":
		p, err = linode.NewLinodeProvider(domainFilter, cfg.DryRun, externaldns.Version)
	case "dnsimple":
		p, err = dnsimple.NewDnsimpleProvider(domainFilter, zoneIDFilter, cfg.DryRun)
	case "infoblox":
		p, err = infoblox.NewInfobloxProvider(
			infoblox.StartupConfig{
				DomainFilter:  domainFilter,
				ZoneIDFilter:  zoneIDFilter,
				Host:          cfg.InfobloxGridHost,
				Port:          cfg.InfobloxWapiPort,
				Username:      cfg.InfobloxWapiUsername,
				Password:      cfg.InfobloxWapiPassword,
				Version:       cfg.InfobloxWapiVersion,
				SSLVerify:     cfg.InfobloxSSLVerify,
				View:          cfg.InfobloxView,
				MaxResults:    cfg.InfobloxMaxResults,
				DryRun:        cfg.DryRun,
				FQDNRexEx:     cfg.InfobloxFQDNRegEx,
				CreatePTR:     cfg.InfobloxCreatePTR,
				CacheDuration: cfg.InfobloxCacheDuration,
			},
		)
	case "dyn":
		p, err = dyn.NewDynProvider(
			dyn.DynConfig{
				DomainFilter:  domainFilter,
				ZoneIDFilter:  zoneIDFilter,
				DryRun:        cfg.DryRun,
				CustomerName:  cfg.DynCustomerName,
				Username:      cfg.DynUsername,
				Password:      cfg.DynPassword,
				MinTTLSeconds: cfg.DynMinTTLSeconds,
				AppVersion:    externaldns.Version,
			},
		)
	case "coredns", "skydns":
		p, err = coredns.NewCoreDNSProvider(domainFilter, cfg.CoreDNSPrefix, cfg.DryRun)
	case "rdns":
		p, err = rdns.NewRDNSProvider(
			rdns.RDNSConfig{
				DomainFilter: domainFilter,
				DryRun:       cfg.DryRun,
			},
		)
	case "exoscale":
		p, err = exoscale.NewExoscaleProvider(cfg.ExoscaleEndpoint, cfg.ExoscaleAPIKey, cfg.ExoscaleAPISecret, cfg.DryRun, exoscale.ExoscaleWithDomain(domainFilter), exoscale.ExoscaleWithLogging()), nil
	case "inmemory":
		p, err = inmemory.NewInMemoryProvider(inmemory.InMemoryInitZones(cfg.InMemoryZones), inmemory.InMemoryWithDomain(domainFilter), inmemory.InMemoryWithLogging()), nil
	case "designate":
		p, err = designate.NewDesignateProvider(domainFilter, cfg.DryRun)
	case "pdns":
		p, err = pdns.NewPDNSProvider(
			ctx,
			pdns.PDNSConfig{
				DomainFilter: domainFilter,
				DryRun:       cfg.DryRun,
				Server:       cfg.PDNSServer,
				APIKey:       cfg.PDNSAPIKey,
				TLSConfig: pdns.TLSConfig{
					TLSEnabled:            cfg.PDNSTLSEnabled,
					CAFilePath:            cfg.TLSCA,
					ClientCertFilePath:    cfg.TLSClientCert,
					ClientCertKeyFilePath: cfg.TLSClientCertKey,
				},
			},
		)
	case "oci":
		var config *oci.OCIConfig
		config, err = oci.LoadOCIConfig(cfg.OCIConfigFile)
		if err == nil {
			p, err = oci.NewOCIProvider(*config, domainFilter, zoneIDFilter, cfg.DryRun)
		}
	case "rfc2136":
		p, err = rfc2136.NewRfc2136Provider(cfg.RFC2136Host, cfg.RFC2136Port, cfg.RFC2136Zone, cfg.RFC2136Insecure, cfg.RFC2136TSIGKeyName, cfg.RFC2136TSIGSecret, cfg.RFC2136TSIGSecretAlg, cfg.RFC2136TAXFR, domainFilter, cfg.DryRun, cfg.RFC2136MinTTL, cfg.RFC2136GSSTSIG, cfg.RFC2136KerberosUsername, cfg.RFC2136KerberosPassword, cfg.RFC2136KerberosRealm, cfg.RFC2136BatchChangeSize, nil)
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
	case "ibmcloud":
		p, err = ibmcloud.NewIBMCloudProvider(cfg.IBMCloudConfigFile, domainFilter, zoneIDFilter, endpointsSource, cfg.IBMCloudProxied, cfg.DryRun)
	case "safedns":
		p, err = safedns.NewSafeDNSProvider(domainFilter, cfg.DryRun)
	default:
		log.Fatalf("unknown dns provider: %s", cfg.Provider)
	}

	return &p, err
}

func createRegistry(cfg *externaldns.Config, p provider.Provider) (registry.Registry, error) {
	var r registry.Registry
	var err error

	switch cfg.Registry {
	case "noop":
		r, err = registry.NewNoopRegistry(p)
	case "txt":
		r, err = registry.NewTXTRegistry(p, cfg.TXTPrefix, cfg.TXTSuffix, cfg.TXTOwnerID, cfg.TXTCacheInterval, cfg.TXTWildcardReplacement, cfg.ManagedDNSRecordTypes)
	case "aws-sd":
		r, err = registry.NewAWSSDRegistry(p.(*awssd.AWSSDProvider), cfg.TXTOwnerID)
	default:
		err = fmt.Errorf("unknown registry: %s", cfg.Registry)
	}

	return r, err
}
