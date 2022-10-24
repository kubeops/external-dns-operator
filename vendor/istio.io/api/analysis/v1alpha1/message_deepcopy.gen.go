// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: analysis/v1alpha1/message.proto

// Describes the structure of messages generated by Istio analyzers.

package v1alpha1

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// DeepCopyInto supports using AnalysisMessageBase within kubernetes types, where deepcopy-gen is used.
func (in *AnalysisMessageBase) DeepCopyInto(out *AnalysisMessageBase) {
	p := proto.Clone(in).(*AnalysisMessageBase)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnalysisMessageBase. Required by controller-gen.
func (in *AnalysisMessageBase) DeepCopy() *AnalysisMessageBase {
	if in == nil {
		return nil
	}
	out := new(AnalysisMessageBase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto supports using AnalysisMessageBase_Type within kubernetes types, where deepcopy-gen is used.
func (in *AnalysisMessageBase_Type) DeepCopyInto(out *AnalysisMessageBase_Type) {
	p := proto.Clone(in).(*AnalysisMessageBase_Type)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnalysisMessageBase_Type. Required by controller-gen.
func (in *AnalysisMessageBase_Type) DeepCopy() *AnalysisMessageBase_Type {
	if in == nil {
		return nil
	}
	out := new(AnalysisMessageBase_Type)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto supports using AnalysisMessageWeakSchema within kubernetes types, where deepcopy-gen is used.
func (in *AnalysisMessageWeakSchema) DeepCopyInto(out *AnalysisMessageWeakSchema) {
	p := proto.Clone(in).(*AnalysisMessageWeakSchema)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnalysisMessageWeakSchema. Required by controller-gen.
func (in *AnalysisMessageWeakSchema) DeepCopy() *AnalysisMessageWeakSchema {
	if in == nil {
		return nil
	}
	out := new(AnalysisMessageWeakSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto supports using AnalysisMessageWeakSchema_ArgType within kubernetes types, where deepcopy-gen is used.
func (in *AnalysisMessageWeakSchema_ArgType) DeepCopyInto(out *AnalysisMessageWeakSchema_ArgType) {
	p := proto.Clone(in).(*AnalysisMessageWeakSchema_ArgType)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnalysisMessageWeakSchema_ArgType. Required by controller-gen.
func (in *AnalysisMessageWeakSchema_ArgType) DeepCopy() *AnalysisMessageWeakSchema_ArgType {
	if in == nil {
		return nil
	}
	out := new(AnalysisMessageWeakSchema_ArgType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto supports using GenericAnalysisMessage within kubernetes types, where deepcopy-gen is used.
func (in *GenericAnalysisMessage) DeepCopyInto(out *GenericAnalysisMessage) {
	p := proto.Clone(in).(*GenericAnalysisMessage)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericAnalysisMessage. Required by controller-gen.
func (in *GenericAnalysisMessage) DeepCopy() *GenericAnalysisMessage {
	if in == nil {
		return nil
	}
	out := new(GenericAnalysisMessage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto supports using InternalErrorAnalysisMessage within kubernetes types, where deepcopy-gen is used.
func (in *InternalErrorAnalysisMessage) DeepCopyInto(out *InternalErrorAnalysisMessage) {
	p := proto.Clone(in).(*InternalErrorAnalysisMessage)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalErrorAnalysisMessage. Required by controller-gen.
func (in *InternalErrorAnalysisMessage) DeepCopy() *InternalErrorAnalysisMessage {
	if in == nil {
		return nil
	}
	out := new(InternalErrorAnalysisMessage)
	in.DeepCopyInto(out)
	return out
}
