package pkg

import (
	"context"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/pkg/apis/externaldns"
	"sigs.k8s.io/external-dns/plan"
	_ "sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/registry"
	"sigs.k8s.io/external-dns/source"
)

func CreateSinglePlanForCRD(cfg *externaldns.Config, r registry.Registry, ctx context.Context, source source.Source) (*plan.Plan, error) {
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

	endpoints, err := source.Endpoints(ctx)
	if err != nil {
		return nil, err
	}

	//missingRecords r.MissingRecords

	plan := plan.Plan{
		Policies:           []plan.Policy{plan.Policies[cfg.Policy]},
		Current:            records,
		Desired:            endpoints,
		DomainFilter:       domainFilter,
		PropertyComparator: r.PropertyValuesEqual,
		ManagedRecords:     cfg.ManagedDNSRecordTypes,
	}

	return &plan, nil
}
