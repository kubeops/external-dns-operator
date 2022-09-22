package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	externaldnsv1alpha1 "kubeops.dev/external-dns-operator/apis/external-dns/v1alpha1"
	"log"
	"regexp"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/pkg/apis/externaldns"
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
	"time"
)

var defaultConfig = &externaldns.Config{
	APIServerURL:                "",
	KubeConfig:                  "/home/rasel/Downloads/rasel-kubeconfig.yaml",
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

func ConvertCRDtoCfg(crd externaldnsv1alpha1.ExternalDNS) (*[]externaldns.Config, error) {

	var configs []externaldns.Config
	for _, entry := range *crd.Spec.Entries {

		// Create a config file for single record
		c := defaultConfig

		if crd.Namespace != "" {
			c.Namespace = crd.Namespace
		}

		if crd.Spec.Kubeconfig != nil {
			c.KubeConfig = *crd.Spec.Kubeconfig
		}
		if crd.Spec.APIServerURL != nil {
			c.APIServerURL = *crd.Spec.APIServerURL
		}
		if crd.Spec.RequestTimeout != nil {
			c.RequestTimeout = *crd.Spec.RequestTimeout
		}

		//Source
		s := entry.Sources
		if s.Names != nil {
			c.Sources = *s.Names
		}
		if s.OCRouterName != nil {
			c.OCPRouterName = *s.OCRouterName
		}
		if s.Namespace != nil {
			c.Namespace = *s.Namespace
		}
		if s.AnnotationFilter != nil {
			c.AnnotationFilter = *s.AnnotationFilter
		}
		if s.LabelFilter != nil {
			c.LabelFilter = *s.LabelFilter
		}
		if s.FQDNTemplate != nil {
			c.FQDNTemplate = *s.FQDNTemplate
		}
		if s.CombineFQDNAndAnnotation != nil {
			c.CombineFQDNAndAnnotation = *s.CombineFQDNAndAnnotation
		}
		if s.IgnoreHostnameAnnotation != nil {
			c.IgnoreHostnameAnnotation = *s.IgnoreHostnameAnnotation
		}
		if s.IgnoreIngressTLSSpec != nil {
			c.IgnoreIngressTLSSpec = *s.IgnoreIngressTLSSpec
		}
		if s.IgnoreIngressRulesSpec != nil {
			c.IgnoreIngressRulesSpec = *s.IgnoreIngressRulesSpec
		}
		if s.GatewayNamespace != nil {
			c.GatewayNamespace = *s.GatewayNamespace
		}
		if s.GatewayLabelFilter != nil {
			c.GatewayLabelFilter = *s.GatewayLabelFilter
		}
		if s.Compatibility != nil {
			c.Compatibility = *s.Compatibility
		}
		if s.PublishInternal != nil {
			c.PublishInternal = *s.PublishInternal
		}
		if s.PublishHostIP != nil {
			c.PublishHostIP = *s.PublishHostIP
		}
		if s.AlwaysPublishNotReadyAddresses != nil {
			c.AlwaysPublishNotReadyAddresses = *s.AlwaysPublishNotReadyAddresses
		}
		if s.ConnectorSourceServer != nil {
			c.ConnectorSourceServer = *s.ConnectorSourceServer
		}
		if s.ServiceTypeFilter != nil {
			c.ServiceTypeFilter = *s.ServiceTypeFilter
		}
		if s.ManageDNSRecordTypes != nil {
			c.ManagedDNSRecordTypes = *s.ManageDNSRecordTypes
		}
		if s.DefaultTargets != nil {
			c.DefaultTargets = *s.DefaultTargets
		}

		// PROVIDER
		p := entry.Provider

		if p.Name != nil {
			c.Provider = *p.Name
		}
		if p.DomainFilter != nil {
			c.DomainFilter = *p.DomainFilter
		}
		if p.ExcludeDomains != nil {
			c.ExcludeDomains = *p.ExcludeDomains
		}
		/*
			if p.RegexDomainFilter != nil {
				c.RegexDomainFilter = p.RegexDomainFilter
			}
			if p.RegexDomainExclusion != nil {
				c.RegexDomainExclusion = p.RegexDomainExclusion
			}
		*/
		if p.ZoneIDFilter != nil {
			c.ZoneIDFilter = *p.ZoneIDFilter
		}

		// For AWS Provider
		aw := p.AWS
		if aw.AWSZoneTagFilter != nil {
			c.AWSZoneTagFilter = *aw.AWSZoneTagFilter
		}
		if aw.AWSZoneType != nil {
			c.AWSZoneType = *aw.AWSZoneType
		}
		if aw.AWSAssumeRole != nil {
			c.AWSAssumeRole = *aw.AWSAssumeRole
		}
		if aw.AWSBatchChangeSize != nil {
			c.AWSBatchChangeSize = *aw.AWSBatchChangeSize
		}
		if aw.AWSBatchChangeInterval != nil {
			c.AWSBatchChangeInterval = *aw.AWSBatchChangeInterval
		}
		if aw.AWSEvaluateTargetHealth != nil {
			c.AWSEvaluateTargetHealth = *aw.AWSEvaluateTargetHealth
		}
		if aw.AWSAPIRetries != nil {
			c.AWSAPIRetries = *aw.AWSAPIRetries
		}
		if aw.AWSPreferCNAME != nil {
			c.AWSPreferCNAME = *aw.AWSPreferCNAME
		}
		if aw.AWSZoneCacheDuration != nil {
			c.AWSZoneCacheDuration = *aw.AWSZoneCacheDuration
		}
		if aw.AWSSDServiceCleanup != nil {
			c.AWSSDServiceCleanup = *aw.AWSSDServiceCleanup
		}

		// For Cloudflare provider
		cfl := p.Cloudflare
		if cfl.CloudflareProxied != nil {
			c.CloudflareProxied = *cfl.CloudflareProxied
		}
		if cfl.CloudflareZonesPerPage != nil {
			c.CloudflareZonesPerPage = *cfl.CloudflareZonesPerPage
		}

		configs = append(configs, *c)
	}

	return &configs, nil
}

func CreateEndpointsSource(ctx context.Context, cfg *externaldns.Config) (source.Source, error) {

	// error is explicitly ignored because the filter is already validated in validation.ValidateConfig
	labelSelector, _ := labels.Parse(cfg.LabelFilter)

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
		return nil, err
	}

	// Combine multiple sources into a single, deduplicated source.
	endpointsSource := source.NewDedupSource(source.NewMultiSource(sources, sourceCfg.DefaultTargets))

	return endpointsSource, nil
}

func CreateProviderFromCfg(ctx context.Context, cfg *externaldns.Config, endpointsSource source.Source) (*provider.Provider, error) {
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

func CreateRegistry(cfg *externaldns.Config, p provider.Provider) (registry.Registry, error) {

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
		//fmt.Println("unknown registry")
		//log.Fatalf("unknown registry: %s", cfg.Registry)
		err = errors.New(fmt.Sprintf("unknown registry: %s", cfg.Registry))
	}

	return r, err
}

/*
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

cfg.Registry -> spec.Registry.Type
cfg.Provider -> spec.Provider.Name
cfg.Sources -> ... SourceInfo.Names
*/
