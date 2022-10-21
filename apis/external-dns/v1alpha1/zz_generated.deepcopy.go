//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	"kmodules.xyz/client-go/api/v1"
	timex "time"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSProvider) DeepCopyInto(out *AWSProvider) {
	*out = *in
	if in.ZoneType != nil {
		in, out := &in.ZoneType, &out.ZoneType
		*out = new(string)
		**out = **in
	}
	if in.ZoneTagFilter != nil {
		in, out := &in.ZoneTagFilter, &out.ZoneTagFilter
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AssumeRole != nil {
		in, out := &in.AssumeRole, &out.AssumeRole
		*out = new(string)
		**out = **in
	}
	if in.BatchChangeSize != nil {
		in, out := &in.BatchChangeSize, &out.BatchChangeSize
		*out = new(int)
		**out = **in
	}
	if in.BatchChangeInterval != nil {
		in, out := &in.BatchChangeInterval, &out.BatchChangeInterval
		*out = new(timex.Duration)
		**out = **in
	}
	if in.EvaluateTargetHealth != nil {
		in, out := &in.EvaluateTargetHealth, &out.EvaluateTargetHealth
		*out = new(bool)
		**out = **in
	}
	if in.APIRetries != nil {
		in, out := &in.APIRetries, &out.APIRetries
		*out = new(int)
		**out = **in
	}
	if in.PreferCNAME != nil {
		in, out := &in.PreferCNAME, &out.PreferCNAME
		*out = new(bool)
		**out = **in
	}
	if in.ZoneCacheDuration != nil {
		in, out := &in.ZoneCacheDuration, &out.ZoneCacheDuration
		*out = new(timex.Duration)
		**out = **in
	}
	if in.SDServiceCleanup != nil {
		in, out := &in.SDServiceCleanup, &out.SDServiceCleanup
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSProvider.
func (in *AWSProvider) DeepCopy() *AWSProvider {
	if in == nil {
		return nil
	}
	out := new(AWSProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureProvider) DeepCopyInto(out *AzureProvider) {
	*out = *in
	if in.ResourceGroup != nil {
		in, out := &in.ResourceGroup, &out.ResourceGroup
		*out = new(string)
		**out = **in
	}
	if in.SubscriptionId != nil {
		in, out := &in.SubscriptionId, &out.SubscriptionId
		*out = new(string)
		**out = **in
	}
	if in.UserAssignedIdentityClientID != nil {
		in, out := &in.UserAssignedIdentityClientID, &out.UserAssignedIdentityClientID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureProvider.
func (in *AzureProvider) DeepCopy() *AzureProvider {
	if in == nil {
		return nil
	}
	out := new(AzureProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudflareProvider) DeepCopyInto(out *CloudflareProvider) {
	*out = *in
	if in.Proxied != nil {
		in, out := &in.Proxied, &out.Proxied
		*out = new(bool)
		**out = **in
	}
	if in.ZonesPerPage != nil {
		in, out := &in.ZonesPerPage, &out.ZonesPerPage
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudflareProvider.
func (in *CloudflareProvider) DeepCopy() *CloudflareProvider {
	if in == nil {
		return nil
	}
	out := new(CloudflareProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNS) DeepCopyInto(out *ExternalDNS) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNS.
func (in *ExternalDNS) DeepCopy() *ExternalDNS {
	if in == nil {
		return nil
	}
	out := new(ExternalDNS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExternalDNS) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSList) DeepCopyInto(out *ExternalDNSList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ExternalDNS, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSList.
func (in *ExternalDNSList) DeepCopy() *ExternalDNSList {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExternalDNSList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSSpec) DeepCopyInto(out *ExternalDNSSpec) {
	*out = *in
	out.ProviderSecretRef = in.ProviderSecretRef
	if in.RequestTimeout != nil {
		in, out := &in.RequestTimeout, &out.RequestTimeout
		*out = new(timex.Duration)
		**out = **in
	}
	in.Source.DeepCopyInto(&out.Source)
	if in.OCRouterName != nil {
		in, out := &in.OCRouterName, &out.OCRouterName
		*out = new(string)
		**out = **in
	}
	if in.GatewayNamespace != nil {
		in, out := &in.GatewayNamespace, &out.GatewayNamespace
		*out = new(string)
		**out = **in
	}
	if in.GatewayLabelFilter != nil {
		in, out := &in.GatewayLabelFilter, &out.GatewayLabelFilter
		*out = new(string)
		**out = **in
	}
	if in.ConnectorSourceServer != nil {
		in, out := &in.ConnectorSourceServer, &out.ConnectorSourceServer
		*out = new(string)
		**out = **in
	}
	if in.ManageDNSRecordTypes != nil {
		in, out := &in.ManageDNSRecordTypes, &out.ManageDNSRecordTypes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DefaultTargets != nil {
		in, out := &in.DefaultTargets, &out.DefaultTargets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DomainFilter != nil {
		in, out := &in.DomainFilter, &out.DomainFilter
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExcludeDomains != nil {
		in, out := &in.ExcludeDomains, &out.ExcludeDomains
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ZoneIDFilter != nil {
		in, out := &in.ZoneIDFilter, &out.ZoneIDFilter
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AWS != nil {
		in, out := &in.AWS, &out.AWS
		*out = new(AWSProvider)
		(*in).DeepCopyInto(*out)
	}
	if in.Cloudflare != nil {
		in, out := &in.Cloudflare, &out.Cloudflare
		*out = new(CloudflareProvider)
		(*in).DeepCopyInto(*out)
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(AzureProvider)
		(*in).DeepCopyInto(*out)
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(Policy)
		**out = **in
	}
	if in.Registry != nil {
		in, out := &in.Registry, &out.Registry
		*out = new(string)
		**out = **in
	}
	if in.TXTOwnerID != nil {
		in, out := &in.TXTOwnerID, &out.TXTOwnerID
		*out = new(string)
		**out = **in
	}
	if in.TXTPrefix != nil {
		in, out := &in.TXTPrefix, &out.TXTPrefix
		*out = new(string)
		**out = **in
	}
	if in.TXTSuffix != nil {
		in, out := &in.TXTSuffix, &out.TXTSuffix
		*out = new(string)
		**out = **in
	}
	if in.TXTWildcardReplacement != nil {
		in, out := &in.TXTWildcardReplacement, &out.TXTWildcardReplacement
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSSpec.
func (in *ExternalDNSSpec) DeepCopy() *ExternalDNSSpec {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSStatus) DeepCopyInto(out *ExternalDNSStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSStatus.
func (in *ExternalDNSStatus) DeepCopy() *ExternalDNSStatus {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressConfig) DeepCopyInto(out *IngressConfig) {
	*out = *in
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
	if in.IgnoreHostnameAnnotation != nil {
		in, out := &in.IgnoreHostnameAnnotation, &out.IgnoreHostnameAnnotation
		*out = new(bool)
		**out = **in
	}
	if in.CombineFQDNAndAnnotation != nil {
		in, out := &in.CombineFQDNAndAnnotation, &out.CombineFQDNAndAnnotation
		*out = new(bool)
		**out = **in
	}
	if in.AnnotationFilter != nil {
		in, out := &in.AnnotationFilter, &out.AnnotationFilter
		*out = new(string)
		**out = **in
	}
	if in.LabelFilter != nil {
		in, out := &in.LabelFilter, &out.LabelFilter
		*out = new(string)
		**out = **in
	}
	if in.FQDNTemplate != nil {
		in, out := &in.FQDNTemplate, &out.FQDNTemplate
		*out = new(string)
		**out = **in
	}
	if in.IgnoreIngressTLSSpec != nil {
		in, out := &in.IgnoreIngressTLSSpec, &out.IgnoreIngressTLSSpec
		*out = new(bool)
		**out = **in
	}
	if in.IgnoreIngressRulesSpec != nil {
		in, out := &in.IgnoreIngressRulesSpec, &out.IgnoreIngressRulesSpec
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressConfig.
func (in *IngressConfig) DeepCopy() *IngressConfig {
	if in == nil {
		return nil
	}
	out := new(IngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfig) DeepCopyInto(out *NodeConfig) {
	*out = *in
	if in.AnnotationFilter != nil {
		in, out := &in.AnnotationFilter, &out.AnnotationFilter
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfig.
func (in *NodeConfig) DeepCopy() *NodeConfig {
	if in == nil {
		return nil
	}
	out := new(NodeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceConfig) DeepCopyInto(out *ServiceConfig) {
	*out = *in
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
	if in.IgnoreHostnameAnnotation != nil {
		in, out := &in.IgnoreHostnameAnnotation, &out.IgnoreHostnameAnnotation
		*out = new(bool)
		**out = **in
	}
	if in.CombineFQDNAndAnnotation != nil {
		in, out := &in.CombineFQDNAndAnnotation, &out.CombineFQDNAndAnnotation
		*out = new(bool)
		**out = **in
	}
	if in.AnnotationFilter != nil {
		in, out := &in.AnnotationFilter, &out.AnnotationFilter
		*out = new(string)
		**out = **in
	}
	if in.LabelFilter != nil {
		in, out := &in.LabelFilter, &out.LabelFilter
		*out = new(string)
		**out = **in
	}
	if in.FQDNTemplate != nil {
		in, out := &in.FQDNTemplate, &out.FQDNTemplate
		*out = new(string)
		**out = **in
	}
	if in.Compatibility != nil {
		in, out := &in.Compatibility, &out.Compatibility
		*out = new(string)
		**out = **in
	}
	if in.PublishInternal != nil {
		in, out := &in.PublishInternal, &out.PublishInternal
		*out = new(bool)
		**out = **in
	}
	if in.PublishHostIP != nil {
		in, out := &in.PublishHostIP, &out.PublishHostIP
		*out = new(bool)
		**out = **in
	}
	if in.AlwaysPublishNotReadyAddresses != nil {
		in, out := &in.AlwaysPublishNotReadyAddresses, &out.AlwaysPublishNotReadyAddresses
		*out = new(bool)
		**out = **in
	}
	if in.ServiceTypeFilter != nil {
		in, out := &in.ServiceTypeFilter, &out.ServiceTypeFilter
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceConfig.
func (in *ServiceConfig) DeepCopy() *ServiceConfig {
	if in == nil {
		return nil
	}
	out := new(ServiceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SourceConfig) DeepCopyInto(out *SourceConfig) {
	*out = *in
	out.Type = in.Type
	if in.Node != nil {
		in, out := &in.Node, &out.Node
		*out = new(NodeConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(ServiceConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SourceConfig.
func (in *SourceConfig) DeepCopy() *SourceConfig {
	if in == nil {
		return nil
	}
	out := new(SourceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TypeInfo) DeepCopyInto(out *TypeInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TypeInfo.
func (in *TypeInfo) DeepCopy() *TypeInfo {
	if in == nil {
		return nil
	}
	out := new(TypeInfo)
	in.DeepCopyInto(out)
	return out
}
