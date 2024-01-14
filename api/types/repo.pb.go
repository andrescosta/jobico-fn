// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: repo.proto

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

type File_FileType int32

const (
	File_NoType     File_FileType = 0
	File_JsonSchema File_FileType = 1
	File_Wasm       File_FileType = 2
)

// Enum value maps for File_FileType.
var (
	File_FileType_name = map[int32]string{
		0: "NoType",
		1: "JsonSchema",
		2: "Wasm",
	}
	File_FileType_value = map[string]int32{
		"NoType":     0,
		"JsonSchema": 1,
		"Wasm":       2,
	}
)

func (x File_FileType) Enum() *File_FileType {
	p := new(File_FileType)
	*p = x
	return p
}

func (x File_FileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (File_FileType) Descriptor() protoreflect.EnumDescriptor {
	return file_repo_proto_enumTypes[0].Descriptor()
}

func (File_FileType) Type() protoreflect.EnumType {
	return &file_repo_proto_enumTypes[0]
}

func (x File_FileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use File_FileType.Descriptor instead.
func (File_FileType) EnumDescriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{9, 0}
}

type UpdateToFileStrRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tenant string `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"` //not used
}

func (x *UpdateToFileStrRequest) Reset() {
	*x = UpdateToFileStrRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateToFileStrRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateToFileStrRequest) ProtoMessage() {}

func (x *UpdateToFileStrRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateToFileStrRequest.ProtoReflect.Descriptor instead.
func (*UpdateToFileStrRequest) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateToFileStrRequest) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

type UpdateToFileStrReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   UpdateType  `protobuf:"varint,1,opt,name=type,proto3,enum=UpdateType" json:"type,omitempty"`
	Object *TenantFile `protobuf:"bytes,2,opt,name=object,proto3" json:"object,omitempty"`
}

func (x *UpdateToFileStrReply) Reset() {
	*x = UpdateToFileStrReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateToFileStrReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateToFileStrReply) ProtoMessage() {}

func (x *UpdateToFileStrReply) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateToFileStrReply.ProtoReflect.Descriptor instead.
func (*UpdateToFileStrReply) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateToFileStrReply) GetType() UpdateType {
	if x != nil {
		return x.Type
	}
	return UpdateType_New
}

func (x *UpdateToFileStrReply) GetObject() *TenantFile {
	if x != nil {
		return x.Object
	}
	return nil
}

type GetAllFileNamesReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantFiles []*TenantFiles `protobuf:"bytes,1,rep,name=tenantFiles,proto3" json:"tenantFiles,omitempty"`
}

func (x *GetAllFileNamesReply) Reset() {
	*x = GetAllFileNamesReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllFileNamesReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllFileNamesReply) ProtoMessage() {}

func (x *GetAllFileNamesReply) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllFileNamesReply.ProtoReflect.Descriptor instead.
func (*GetAllFileNamesReply) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{2}
}

func (x *GetAllFileNamesReply) GetTenantFiles() []*TenantFiles {
	if x != nil {
		return x.TenantFiles
	}
	return nil
}

type AddFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantFile *TenantFile `protobuf:"bytes,1,opt,name=tenantFile,proto3" json:"tenantFile,omitempty"`
}

func (x *AddFileRequest) Reset() {
	*x = AddFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddFileRequest) ProtoMessage() {}

func (x *AddFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddFileRequest.ProtoReflect.Descriptor instead.
func (*AddFileRequest) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{3}
}

func (x *AddFileRequest) GetTenantFile() *TenantFile {
	if x != nil {
		return x.TenantFile
	}
	return nil
}

type AddFileReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *AddFileReply) Reset() {
	*x = AddFileReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddFileReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddFileReply) ProtoMessage() {}

func (x *AddFileReply) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddFileReply.ProtoReflect.Descriptor instead.
func (*AddFileReply) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{4}
}

func (x *AddFileReply) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type GetFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantFile *TenantFile `protobuf:"bytes,1,opt,name=tenantFile,proto3" json:"tenantFile,omitempty"`
}

func (x *GetFileRequest) Reset() {
	*x = GetFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileRequest) ProtoMessage() {}

func (x *GetFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileRequest.ProtoReflect.Descriptor instead.
func (*GetFileRequest) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{5}
}

func (x *GetFileRequest) GetTenantFile() *TenantFile {
	if x != nil {
		return x.TenantFile
	}
	return nil
}

type GetFileReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *File `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *GetFileReply) Reset() {
	*x = GetFileReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileReply) ProtoMessage() {}

func (x *GetFileReply) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileReply.ProtoReflect.Descriptor instead.
func (*GetFileReply) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{6}
}

func (x *GetFileReply) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

type TenantFiles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tenant string  `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
	Files  []*File `protobuf:"bytes,2,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *TenantFiles) Reset() {
	*x = TenantFiles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantFiles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantFiles) ProtoMessage() {}

func (x *TenantFiles) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantFiles.ProtoReflect.Descriptor instead.
func (*TenantFiles) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{7}
}

func (x *TenantFiles) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

func (x *TenantFiles) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type TenantFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tenant string `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
	File   *File  `protobuf:"bytes,2,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *TenantFile) Reset() {
	*x = TenantFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantFile) ProtoMessage() {}

func (x *TenantFile) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantFile.ProtoReflect.Descriptor instead.
func (*TenantFile) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{8}
}

func (x *TenantFile) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

func (x *TenantFile) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    File_FileType `protobuf:"varint,1,opt,name=type,proto3,enum=File_FileType" json:"type,omitempty"`
	Name    string        `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Content []byte        `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repo_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_repo_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_repo_proto_rawDescGZIP(), []int{9}
}

func (x *File) GetType() File_FileType {
	if x != nil {
		return x.Type
	}
	return File_NoType
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_repo_proto protoreflect.FileDescriptor

var file_repo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x30, 0x0a, 0x16, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x22, 0x5c, 0x0a, 0x14,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x1f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x46, 0x0a, 0x14, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x2e, 0x0a, 0x0b, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x0b, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x22, 0x3d, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x22, 0x28, 0x0a, 0x0c, 0x41, 0x64, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3d, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a,
	0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0b, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x0a,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x29, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x19, 0x0a, 0x04, 0x66, 0x69,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x42, 0x0a, 0x0b, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x05,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x3f, 0x0a, 0x0a, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12,
	0x19, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x8a, 0x01, 0x0a, 0x04, 0x46,
	0x69, 0x6c, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x30, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x6f, 0x54, 0x79, 0x70, 0x65, 0x10, 0x00, 0x12, 0x0e, 0x0a,
	0x0a, 0x4a, 0x73, 0x6f, 0x6e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x10, 0x01, 0x12, 0x08, 0x0a,
	0x04, 0x57, 0x61, 0x73, 0x6d, 0x10, 0x02, 0x32, 0xd4, 0x01, 0x0a, 0x04, 0x52, 0x65, 0x70, 0x6f,
	0x12, 0x29, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0f, 0x2e, 0x47, 0x65,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x47,
	0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x29, 0x0a, 0x07, 0x41,
	0x64, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0f, 0x2e, 0x41, 0x64, 0x64, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x41, 0x64, 0x64, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x45, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x54, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x12, 0x17, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x54, 0x6f, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x15, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x46, 0x69, 0x6c,
	0x65, 0x53, 0x74, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x30, 0x01, 0x12, 0x2f, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x12, 0x05, 0x2e, 0x56, 0x6f, 0x69, 0x64, 0x1a, 0x15, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c,
	0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x08,
	0x5a, 0x06, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_repo_proto_rawDescOnce sync.Once
	file_repo_proto_rawDescData = file_repo_proto_rawDesc
)

func file_repo_proto_rawDescGZIP() []byte {
	file_repo_proto_rawDescOnce.Do(func() {
		file_repo_proto_rawDescData = protoimpl.X.CompressGZIP(file_repo_proto_rawDescData)
	})
	return file_repo_proto_rawDescData
}

var file_repo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_repo_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_repo_proto_goTypes = []interface{}{
	(File_FileType)(0),             // 0: File.FileType
	(*UpdateToFileStrRequest)(nil), // 1: UpdateToFileStrRequest
	(*UpdateToFileStrReply)(nil),   // 2: UpdateToFileStrReply
	(*GetAllFileNamesReply)(nil),   // 3: GetAllFileNamesReply
	(*AddFileRequest)(nil),         // 4: AddFileRequest
	(*AddFileReply)(nil),           // 5: AddFileReply
	(*GetFileRequest)(nil),         // 6: GetFileRequest
	(*GetFileReply)(nil),           // 7: GetFileReply
	(*TenantFiles)(nil),            // 8: TenantFiles
	(*TenantFile)(nil),             // 9: TenantFile
	(*File)(nil),                   // 10: File
	(UpdateType)(0),                // 11: UpdateType
	(*Void)(nil),                   // 12: Void
}
var file_repo_proto_depIdxs = []int32{
	11, // 0: UpdateToFileStrReply.type:type_name -> UpdateType
	9,  // 1: UpdateToFileStrReply.object:type_name -> TenantFile
	8,  // 2: GetAllFileNamesReply.tenantFiles:type_name -> TenantFiles
	9,  // 3: AddFileRequest.tenantFile:type_name -> TenantFile
	9,  // 4: GetFileRequest.tenantFile:type_name -> TenantFile
	10, // 5: GetFileReply.file:type_name -> File
	10, // 6: TenantFiles.files:type_name -> File
	10, // 7: TenantFile.file:type_name -> File
	0,  // 8: File.type:type_name -> File.FileType
	6,  // 9: Repo.GetFile:input_type -> GetFileRequest
	4,  // 10: Repo.AddFile:input_type -> AddFileRequest
	1,  // 11: Repo.UpdateToFileStr:input_type -> UpdateToFileStrRequest
	12, // 12: Repo.GetAllFileNames:input_type -> Void
	7,  // 13: Repo.GetFile:output_type -> GetFileReply
	5,  // 14: Repo.AddFile:output_type -> AddFileReply
	2,  // 15: Repo.UpdateToFileStr:output_type -> UpdateToFileStrReply
	3,  // 16: Repo.GetAllFileNames:output_type -> GetAllFileNamesReply
	13, // [13:17] is the sub-list for method output_type
	9,  // [9:13] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_repo_proto_init() }
func file_repo_proto_init() {
	if File_repo_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_repo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateToFileStrRequest); i {
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
		file_repo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateToFileStrReply); i {
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
		file_repo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllFileNamesReply); i {
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
		file_repo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddFileRequest); i {
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
		file_repo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddFileReply); i {
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
		file_repo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileRequest); i {
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
		file_repo_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileReply); i {
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
		file_repo_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantFiles); i {
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
		file_repo_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantFile); i {
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
		file_repo_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
			RawDescriptor: file_repo_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_repo_proto_goTypes,
		DependencyIndexes: file_repo_proto_depIdxs,
		EnumInfos:         file_repo_proto_enumTypes,
		MessageInfos:      file_repo_proto_msgTypes,
	}.Build()
	File_repo_proto = out.File
	file_repo_proto_rawDesc = nil
	file_repo_proto_goTypes = nil
	file_repo_proto_depIdxs = nil
}
