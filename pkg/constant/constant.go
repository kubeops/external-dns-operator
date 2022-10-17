package constant

// kubebuilder:validation:Enum:=sync;upsert-only;create-only
type Policy string

const (
	//Policy
	PolicySync       Policy = "sync"
	PolicyUpsertOnly Policy = "upsert-only"
	PolicyCreateOnly Policy = "create-only"
)

func (p Policy) String() string {
	return string(p)
}

type ExternalDNSPhase string

const (
	//ExternalDNSPhase
	ExternalDNSPhaseCurrent    ExternalDNSPhase = "Current"
	ExternalDNSPhaseFailed     ExternalDNSPhase = "Failed"
	ExternalDNSPhaseInProgress ExternalDNSPhase = "InProgress"
)

// kubebuilder:validation:Enum:=aws;cloudflare
type Provider string

const (
	//Provider
	ProviderAWS        Provider = "aws"
	ProviderCloudflare Provider = "cloudflare"
)

func (p Provider) String() string {
	return string(p)
}

const (
	//ConditionType
	CreateAndRegisterWatcher = "CreateAndRegisterWatcher"
	CreateAndSetCredential   = "CreateAndSetCredential"
	CreateAndApplyPlan       = "CreateAndApplyPlan"
)
