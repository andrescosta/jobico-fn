// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.0
// source: control.proto

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Schema string `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Event) GetSchema() string {
	if x != nil {
		return x.Schema
	}
	return ""
}

type ProcessorDefinition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event   *Event `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	Package string `protobuf:"bytes,2,opt,name=package,proto3" json:"package,omitempty"`
}

func (x *ProcessorDefinition) Reset() {
	*x = ProcessorDefinition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessorDefinition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessorDefinition) ProtoMessage() {}

func (x *ProcessorDefinition) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessorDefinition.ProtoReflect.Descriptor instead.
func (*ProcessorDefinition) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessorDefinition) GetEvent() *Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *ProcessorDefinition) GetPackage() string {
	if x != nil {
		return x.Package
	}
	return ""
}

type GetEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *GetEventRequest) Reset() {
	*x = GetEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventRequest) ProtoMessage() {}

func (x *GetEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventRequest.ProtoReflect.Descriptor instead.
func (*GetEventRequest) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{2}
}

func (x *GetEventRequest) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type GetEventReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *Result  `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	Events []*Event `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *GetEventReply) Reset() {
	*x = GetEventReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventReply) ProtoMessage() {}

func (x *GetEventReply) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventReply.ProtoReflect.Descriptor instead.
func (*GetEventReply) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{3}
}

func (x *GetEventReply) GetResult() *Result {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *GetEventReply) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type AddEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event []*Event `protobuf:"bytes,1,rep,name=event,proto3" json:"event,omitempty"`
}

func (x *AddEventRequest) Reset() {
	*x = AddEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddEventRequest) ProtoMessage() {}

func (x *AddEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddEventRequest.ProtoReflect.Descriptor instead.
func (*AddEventRequest) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{4}
}

func (x *AddEventRequest) GetEvent() []*Event {
	if x != nil {
		return x.Event
	}
	return nil
}

type AddEventReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *Result `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *AddEventReply) Reset() {
	*x = AddEventReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddEventReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddEventReply) ProtoMessage() {}

func (x *AddEventReply) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddEventReply.ProtoReflect.Descriptor instead.
func (*AddEventReply) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{5}
}

func (x *AddEventReply) GetResult() *Result {
	if x != nil {
		return x.Result
	}
	return nil
}

type AddProcessorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Definition []*ProcessorDefinition `protobuf:"bytes,1,rep,name=definition,proto3" json:"definition,omitempty"`
}

func (x *AddProcessorRequest) Reset() {
	*x = AddProcessorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddProcessorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddProcessorRequest) ProtoMessage() {}

func (x *AddProcessorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddProcessorRequest.ProtoReflect.Descriptor instead.
func (*AddProcessorRequest) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{6}
}

func (x *AddProcessorRequest) GetDefinition() []*ProcessorDefinition {
	if x != nil {
		return x.Definition
	}
	return nil
}

type AddProcessorReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *Result `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *AddProcessorReply) Reset() {
	*x = AddProcessorReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddProcessorReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddProcessorReply) ProtoMessage() {}

func (x *AddProcessorReply) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddProcessorReply.ProtoReflect.Descriptor instead.
func (*AddProcessorReply) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{7}
}

func (x *AddProcessorReply) GetResult() *Result {
	if x != nil {
		return x.Result
	}
	return nil
}

type GetProcessorsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Definitions []*ProcessorDefinition `protobuf:"bytes,1,rep,name=definitions,proto3" json:"definitions,omitempty"`
}

func (x *GetProcessorsRequest) Reset() {
	*x = GetProcessorsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProcessorsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProcessorsRequest) ProtoMessage() {}

func (x *GetProcessorsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProcessorsRequest.ProtoReflect.Descriptor instead.
func (*GetProcessorsRequest) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{8}
}

func (x *GetProcessorsRequest) GetDefinitions() []*ProcessorDefinition {
	if x != nil {
		return x.Definitions
	}
	return nil
}

type GetProcessorsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Definitions []*ProcessorDefinition `protobuf:"bytes,1,rep,name=definitions,proto3" json:"definitions,omitempty"`
}

func (x *GetProcessorsReply) Reset() {
	*x = GetProcessorsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProcessorsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProcessorsReply) ProtoMessage() {}

func (x *GetProcessorsReply) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProcessorsReply.ProtoReflect.Descriptor instead.
func (*GetProcessorsReply) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{9}
}

func (x *GetProcessorsReply) GetDefinitions() []*ProcessorDefinition {
	if x != nil {
		return x.Definitions
	}
	return nil
}

var File_control_proto protoreflect.FileDescriptor

var file_control_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x15, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x22, 0x4d, 0x0a, 0x13, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x6f, 0x72, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c,
	0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x22, 0x31, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x06, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x50, 0x0a, 0x0d, 0x47, 0x65, 0x74,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1f, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1e, 0x0a, 0x06, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x2f, 0x0a, 0x0f, 0x41,
	0x64, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c,
	0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x30, 0x0a, 0x0d,
	0x41, 0x64, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1f, 0x0a,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x4b,
	0x0a, 0x13, 0x41, 0x64, 0x64, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0a, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x6f, 0x72, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0a, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x34, 0x0a, 0x11, 0x41,
	0x64, 0x64, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x1f, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x4e, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f,
	0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x0b, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x22, 0x4c, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f,
	0x72, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x36, 0x0a, 0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x32,
	0xe6, 0x01, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x2f, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x10, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x47, 0x65, 0x74,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x09,
	0x41, 0x64, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x10, 0x2e, 0x41, 0x64, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x41, 0x64,
	0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3a, 0x0a,
	0x0c, 0x41, 0x64, 0x64, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x14, 0x2e,
	0x41, 0x64, 0x64, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x41, 0x64, 0x64, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0d, 0x47, 0x65, 0x74,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x12, 0x15, 0x2e, 0x47, 0x65, 0x74,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_control_proto_rawDescOnce sync.Once
	file_control_proto_rawDescData = file_control_proto_rawDesc
)

func file_control_proto_rawDescGZIP() []byte {
	file_control_proto_rawDescOnce.Do(func() {
		file_control_proto_rawDescData = protoimpl.X.CompressGZIP(file_control_proto_rawDescData)
	})
	return file_control_proto_rawDescData
}

var file_control_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_control_proto_goTypes = []interface{}{
	(*Event)(nil),                // 0: Event
	(*ProcessorDefinition)(nil),  // 1: ProcessorDefinition
	(*GetEventRequest)(nil),      // 2: GetEventRequest
	(*GetEventReply)(nil),        // 3: GetEventReply
	(*AddEventRequest)(nil),      // 4: AddEventRequest
	(*AddEventReply)(nil),        // 5: AddEventReply
	(*AddProcessorRequest)(nil),  // 6: AddProcessorRequest
	(*AddProcessorReply)(nil),    // 7: AddProcessorReply
	(*GetProcessorsRequest)(nil), // 8: GetProcessorsRequest
	(*GetProcessorsReply)(nil),   // 9: GetProcessorsReply
	(*Result)(nil),               // 10: Result
}
var file_control_proto_depIdxs = []int32{
	0,  // 0: ProcessorDefinition.event:type_name -> Event
	0,  // 1: GetEventRequest.events:type_name -> Event
	10, // 2: GetEventReply.result:type_name -> Result
	0,  // 3: GetEventReply.events:type_name -> Event
	0,  // 4: AddEventRequest.event:type_name -> Event
	10, // 5: AddEventReply.result:type_name -> Result
	1,  // 6: AddProcessorRequest.definition:type_name -> ProcessorDefinition
	10, // 7: AddProcessorReply.result:type_name -> Result
	1,  // 8: GetProcessorsRequest.definitions:type_name -> ProcessorDefinition
	1,  // 9: GetProcessorsReply.definitions:type_name -> ProcessorDefinition
	2,  // 10: Control.GetEvents:input_type -> GetEventRequest
	4,  // 11: Control.AddEvents:input_type -> AddEventRequest
	6,  // 12: Control.AddProcessor:input_type -> AddProcessorRequest
	8,  // 13: Control.GetProcessors:input_type -> GetProcessorsRequest
	3,  // 14: Control.GetEvents:output_type -> GetEventReply
	5,  // 15: Control.AddEvents:output_type -> AddEventReply
	7,  // 16: Control.AddProcessor:output_type -> AddProcessorReply
	9,  // 17: Control.GetProcessors:output_type -> GetProcessorsReply
	14, // [14:18] is the sub-list for method output_type
	10, // [10:14] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_control_proto_init() }
func file_control_proto_init() {
	if File_control_proto != nil {
		return
	}
	file_common_messages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_control_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessorDefinition); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddEventReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddProcessorRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddProcessorReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProcessorsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_control_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProcessorsReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_control_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_control_proto_goTypes,
		DependencyIndexes: file_control_proto_depIdxs,
		MessageInfos:      file_control_proto_msgTypes,
	}.Build()
	File_control_proto = out.File
	file_control_proto_rawDesc = nil
	file_control_proto_goTypes = nil
	file_control_proto_depIdxs = nil
}
