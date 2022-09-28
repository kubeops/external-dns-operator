package plan

import (
	"context"
	"fmt"
	"k8s.io/klog/v2"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/pkg/apis/externaldns"
	"sigs.k8s.io/external-dns/plan"
	_ "sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
	"sigs.k8s.io/external-dns/registry"
	"sigs.k8s.io/external-dns/source"
)

func CreateAndApplySinglePlanForCRD(ctx context.Context, cfg *externaldns.Config, r registry.Registry, endpointSource source.Source) error {

	var domainFilter endpoint.DomainFilter
	if cfg.RegexDomainFilter.String() != "" {
		domainFilter = endpoint.NewRegexDomainFilter(cfg.RegexDomainFilter, cfg.RegexDomainExclusion)
	} else {
		domainFilter = endpoint.NewDomainFilterWithExclusions(cfg.DomainFilter, cfg.ExcludeDomains)
	}

	records, err := r.Records(ctx)
	if err != nil {
		return err
	}

	missingRecords := r.MissingRecords()

	ctx = context.WithValue(ctx, provider.RecordsContextKey, records)
	endpoints, err := endpointSource.Endpoints(ctx)
	if err != nil {
		return err
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
				return err
			}
			klog.Info("all missing records are created")
		}
	}

	plan := &plan.Plan{
		Policies:           []plan.Policy{plan.Policies[cfg.Policy]},
		Current:            records,
		Desired:            endpoints,
		DomainFilter:       domainFilter,
		PropertyComparator: r.PropertyValuesEqual,
		ManagedRecords:     cfg.ManagedDNSRecordTypes,
	}

	/*
			plan.Desired = []*endpoint.Endpoint{
				{
					DNSName:       "api.superm4n.tk",
					Targets:       []string{"127.0.0.1"},
					RecordType:    "A",
					SetIdentifier: "",
					RecordTTL:     0,
					Labels:        nil,
					ProviderSpecific: endpoint.ProviderSpecific{
						{
							Name:  "Region",
							Value: "us-east-1",
						},
					},
				},
			}


		for _, d := range plan.Desired {
			klog.Info("Desird =====================> ", d.String())
		}

		for _, e := range plan.Current {
			klog.Info("Current ================> ", e.String())
		}
	*/

	plan = plan.Calculate()

	if plan.Changes.HasChanges() {
		err = r.ApplyChanges(ctx, plan.Changes)
		if err != nil {
			klog.Info(err.Error())
			fmt.Println("failed to apply changes for plan")
			return err
		}
		fmt.Println("Plan Applied")
	} else {
		fmt.Println("all records are already up to date")
	}

	return nil
}
