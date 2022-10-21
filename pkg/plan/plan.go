package plan

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"log"
	"regexp"
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
	"strings"
	"time"
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

func createAndApplyPlan(ctx context.Context, cfg *externaldns.Config, r registry.Registry, endpointSource source.Source) (string, error) {

	var domainFilter endpoint.DomainFilter
	if cfg.RegexDomainFilter.String() != "" {
		domainFilter = endpoint.NewRegexDomainFilter(cfg.RegexDomainFilter, cfg.RegexDomainExclusion)
	} else {
		domainFilter = endpoint.NewDomainFilterWithExclusions(cfg.DomainFilter, cfg.ExcludeDomains)
	}

	records, err := r.Records(ctx)
	if err != nil {
		return "", err
	}

	missingRecords := r.MissingRecords()

	ctx = context.WithValue(ctx, provider.RecordsContextKey, records)
	endpoints, err := endpointSource.Endpoints(ctx)
	if err != nil {
		return "", err
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
				return "", err
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

	var successMsg = ""

	if pln.Changes.HasChanges() {
		err = r.ApplyChanges(ctx, pln.Changes)
		if err != nil {
			klog.Error("failed to apply plan")
			return "", err
		}
		klog.Info("plan applied")
		successMsg = "plan applied"
	} else {
		klog.Info("all records are already up to date")
		successMsg = "all records are already up to date"
	}

	return successMsg, nil
}

func convertEDNSObjectToCfg(crd *externaldnsv1alpha1.ExternalDNS) (*externaldns.Config, error) {

	// Create a config file for single record
	c := defaultConfig

	if crd.Namespace != "" {
		c.Namespace = crd.Namespace
	}

	if crd.Spec.RequestTimeout != nil {
		c.RequestTimeout = *crd.Spec.RequestTimeout
	}

	s := crd.Spec

	//SOURCE
	var sources []string
	sources = append(sources, strings.ToLower(s.Source.Type.Kind))
	// sources[] must contain strings that are lower cased
	c.Sources = sources

	if s.OCRouterName != nil {
		c.OCPRouterName = *s.OCRouterName
	}
	if s.GatewayNamespace != nil {
		c.GatewayNamespace = *s.GatewayNamespace
	}
	if s.GatewayLabelFilter != nil {
		c.GatewayLabelFilter = *s.GatewayLabelFilter
	}
	if s.ManageDNSRecordTypes != nil {
		c.ManagedDNSRecordTypes = s.ManageDNSRecordTypes
	}
	if s.DefaultTargets != nil {
		c.DefaultTargets = s.DefaultTargets
	}

	ss := s.Source

	// For Node
	if ss.Node != nil && ss.Type.Kind == "Node" {
		ssn := ss.Node
		c.FQDNTemplate = ssn.FQDNTemplate
		if ssn.AnnotationFilter != nil {
			c.AnnotationFilter = *ssn.AnnotationFilter
		}
	}

	// For Service
	if ss.Service != nil && ss.Type.Kind == "Service" {
		sss := ss.Service
		if sss.LabelFilter != nil {
			c.LabelFilter = *sss.LabelFilter
		}
		if sss.Namespace != nil {
			c.Namespace = *sss.Namespace
		}
		if sss.AnnotationFilter != nil {
			c.AnnotationFilter = *sss.AnnotationFilter
		}
		if sss.FQDNTemplate != nil {
			c.FQDNTemplate = *sss.FQDNTemplate
		}
		if sss.CombineFQDNAndAnnotation != nil {
			c.CombineFQDNAndAnnotation = *sss.CombineFQDNAndAnnotation
		}
		if sss.Compatibility != nil {
			c.Compatibility = *sss.Compatibility
		}
		if sss.PublishInternal != nil {
			c.PublishInternal = *sss.PublishInternal
		}
		if sss.PublishHostIP != nil {
			c.PublishHostIP = *sss.PublishHostIP
		}
		if sss.AlwaysPublishNotReadyAddresses != nil {
			c.AlwaysPublishNotReadyAddresses = *sss.AlwaysPublishNotReadyAddresses
		}
		if sss.ServiceTypeFilter != nil {
			c.ServiceTypeFilter = sss.ServiceTypeFilter
		}
		if sss.IgnoreHostnameAnnotation != nil {
			c.IgnoreHostnameAnnotation = *sss.IgnoreHostnameAnnotation
		}
	}

	// For Ingress
	if ss.Ingress != nil && ss.Type.Kind == "Ingress" {
		ssi := ss.Ingress
		if ssi.IgnoreIngressRulesSpec != nil {
			c.IgnoreIngressRulesSpec = *ssi.IgnoreIngressRulesSpec
		}
		if ssi.IgnoreHostnameAnnotation != nil {
			c.IgnoreHostnameAnnotation = *ssi.IgnoreHostnameAnnotation
		}
		if ssi.FQDNTemplate != nil {
			c.FQDNTemplate = *ssi.FQDNTemplate
		}
		if ssi.Namespace != nil {
			c.Namespace = *ssi.Namespace
		}
		if ssi.AnnotationFilter != nil {
			c.AnnotationFilter = *ssi.AnnotationFilter
		}
		if ssi.CombineFQDNAndAnnotation != nil {
			c.CombineFQDNAndAnnotation = *ssi.CombineFQDNAndAnnotation
		}
		if ssi.IgnoreIngressTLSSpec != nil {
			c.IgnoreIngressTLSSpec = *ssi.IgnoreIngressTLSSpec
		}
		if ssi.LabelFilter != nil {
			c.LabelFilter = *ssi.LabelFilter
		}
	}

	// PROVIDER
	c.Provider = s.Provider.String()

	if s.DomainFilter != nil {
		c.DomainFilter = s.DomainFilter
	}
	if s.ExcludeDomains != nil {
		c.ExcludeDomains = s.ExcludeDomains
	}
	if s.ZoneIDFilter != nil {
		c.ZoneIDFilter = s.ZoneIDFilter
	}

	// for aws provider
	if s.AWS != nil {

		aw := s.AWS
		if aw.ZoneTagFilter != nil {
			c.AWSZoneTagFilter = aw.ZoneTagFilter
		}
		if aw.ZoneType != nil {
			c.AWSZoneType = *aw.ZoneType
		}
		if aw.AssumeRole != nil {
			c.AWSAssumeRole = *aw.AssumeRole
		}
		if aw.BatchChangeSize != nil {
			c.AWSBatchChangeSize = *aw.BatchChangeSize
		}
		if aw.BatchChangeInterval != nil {
			c.AWSBatchChangeInterval = *aw.BatchChangeInterval
		}
		if aw.EvaluateTargetHealth != nil {
			c.AWSEvaluateTargetHealth = *aw.EvaluateTargetHealth
		}
		if aw.APIRetries != nil {
			c.AWSAPIRetries = *aw.APIRetries
		}
		if aw.PreferCNAME != nil {
			c.AWSPreferCNAME = *aw.PreferCNAME
		}
		if aw.ZoneCacheDuration != nil {
			c.AWSZoneCacheDuration = *aw.ZoneCacheDuration
		}
		if aw.SDServiceCleanup != nil {
			c.AWSSDServiceCleanup = *aw.SDServiceCleanup
		}
	}

	// for cloudflare provider
	if s.Cloudflare != nil {

		cfl := s.Cloudflare
		if cfl.Proxied != nil {
			c.CloudflareProxied = *cfl.Proxied
		}

		if cfl.ZonesPerPage != nil {
			c.CloudflareZonesPerPage = *cfl.ZonesPerPage
		}
	}

	// for azure provide
	if s.Azure != nil {
		az := s.Azure

		// AzureConfigFile is only for Azure provider, not for Azure-Private-DNS
		c.AzureConfigFile = fmt.Sprintf("/tmp/%s-%s-credential")

		if az.SubscriptionId != nil {
			c.AzureSubscriptionID = *az.SubscriptionId
		}
		if az.ResourceGroup != nil {
			c.AzureResourceGroup = *az.ResourceGroup
		}
		if az.UserAssignedIdentityClientID != nil {
			c.AzureUserAssignedIdentityClientID = *az.UserAssignedIdentityClientID
		}
	}

	// POLICY

	if s.Policy != nil {
		c.Policy = s.Policy.String()
	}

	// REGISTRY
	if s.Registry != nil {
		c.Registry = *s.Registry
	}
	if s.TXTOwnerID != nil {
		c.TXTOwnerID = *s.TXTOwnerID
	}
	if s.TXTPrefix != nil {
		c.TXTPrefix = *s.TXTPrefix
	}
	if s.TXTSuffix != nil {
		c.TXTSuffix = *s.TXTSuffix
	}
	if s.TXTWildcardReplacement != nil {
		c.TXTWildcardReplacement = *s.TXTWildcardReplacement
	}

	return &c, nil
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
	case "aws-sd":
		// Check that only compatible Registry is used with AWS-SD
		if cfg.Registry != "noop" && cfg.Registry != "aws-sd" {
			// removed the log notification
			//log.Infof("Registry \"%s\" cannot be used with AWS Cloud Map. Switching to \"aws-sd\".", cfg.Registry)
			cfg.Registry = "aws-sd"
		}
		p, err = awssd.NewAWSSDProvider(domainFilter, cfg.AWSZoneType, cfg.AWSAssumeRole, cfg.DryRun, cfg.AWSSDServiceCleanup, cfg.TXTOwnerID)
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
	if err != nil {
		err = errors.New(fmt.Sprintf("unknown dns provider: %s", cfg.Provider))
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
		err = errors.New(fmt.Sprintf("unknown registry: %s", cfg.Registry))
	}

	return r, err
}

func SetDNSRecords(edns *externaldnsv1alpha1.ExternalDNS, ctx context.Context) (string, error) {

	cfg, err := convertEDNSObjectToCfg(edns)
	if err != nil {
		klog.Error("failed to convert crd into cfg.", err.Error())
		return "", err
	}
	endpointsSource, err := createEndpointsSource(ctx, cfg)
	if err != nil {
		klog.Error("failed to create endpoints source.", err.Error())
		return "", err
	}

	pvdr, err := createProviderFromCfg(ctx, cfg, endpointsSource)
	if err != nil {
		klog.Error("failed to create provider.", err.Error())
		return "", err
	}

	reg, err := createRegistry(cfg, *pvdr)
	if err != nil {
		klog.Errorf("failed to create Registry.", err.Error())
		return "", err
	}

	var successMsg string
	successMsg, err = createAndApplyPlan(ctx, cfg, reg, endpointsSource)
	if err != nil {
		klog.Errorf("failed to create and apply plan.", err.Error())
		return "", err
	}

	return successMsg, nil
}
