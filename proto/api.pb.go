// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.3
// source: proto/api.proto

package proto

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

type GetContext struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token   string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	GuildID string `protobuf:"bytes,2,opt,name=guildID,proto3" json:"guildID,omitempty"`
}

func (x *GetContext) Reset() {
	*x = GetContext{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContext) ProtoMessage() {}

func (x *GetContext) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetContext.ProtoReflect.Descriptor instead.
func (*GetContext) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{0}
}

func (x *GetContext) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *GetContext) GetGuildID() string {
	if x != nil {
		return x.GuildID
	}
	return ""
}

type Args struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=Description,proto3" json:"Description,omitempty"`
	Required    bool   `protobuf:"varint,3,opt,name=Required,proto3" json:"Required,omitempty"`
}

func (x *Args) Reset() {
	*x = Args{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Args) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Args) ProtoMessage() {}

func (x *Args) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Args.ProtoReflect.Descriptor instead.
func (*Args) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{1}
}

func (x *Args) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Args) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Args) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

type Commands struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string  `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Description string  `protobuf:"bytes,2,opt,name=Description,proto3" json:"Description,omitempty"`
	Args        []*Args `protobuf:"bytes,3,rep,name=Args,proto3" json:"Args,omitempty"`
}

func (x *Commands) Reset() {
	*x = Commands{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commands) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commands) ProtoMessage() {}

func (x *Commands) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commands.ProtoReflect.Descriptor instead.
func (*Commands) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{2}
}

func (x *Commands) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Commands) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Commands) GetArgs() []*Args {
	if x != nil {
		return x.Args
	}
	return nil
}

type BuildFunc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string      `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Runtime     string      `protobuf:"bytes,2,opt,name=Runtime,proto3" json:"Runtime,omitempty"`
	GuildID     string      `protobuf:"bytes,3,opt,name=GuildID,proto3" json:"GuildID,omitempty"`
	Version     string      `protobuf:"bytes,4,opt,name=Version,proto3" json:"Version,omitempty"`
	Description string      `protobuf:"bytes,5,opt,name=Description,proto3" json:"Description,omitempty"`
	Commands    []*Commands `protobuf:"bytes,6,rep,name=Commands,proto3" json:"Commands,omitempty"`
}

func (x *BuildFunc) Reset() {
	*x = BuildFunc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildFunc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildFunc) ProtoMessage() {}

func (x *BuildFunc) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildFunc.ProtoReflect.Descriptor instead.
func (*BuildFunc) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{3}
}

func (x *BuildFunc) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BuildFunc) GetRuntime() string {
	if x != nil {
		return x.Runtime
	}
	return ""
}

func (x *BuildFunc) GetGuildID() string {
	if x != nil {
		return x.GuildID
	}
	return ""
}

func (x *BuildFunc) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *BuildFunc) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *BuildFunc) GetCommands() []*Commands {
	if x != nil {
		return x.Commands
	}
	return nil
}

type ContextResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientID       string `protobuf:"bytes,1,opt,name=ClientID,proto3" json:"ClientID,omitempty"`
	GuildID        string `protobuf:"bytes,2,opt,name=GuildID,proto3" json:"GuildID,omitempty"`
	GuildName      string `protobuf:"bytes,3,opt,name=GuildName,proto3" json:"GuildName,omitempty"`
	ServerURL      string `protobuf:"bytes,4,opt,name=ServerURL,proto3" json:"ServerURL,omitempty"`
	CurrentContext bool   `protobuf:"varint,5,opt,name=CurrentContext,proto3" json:"CurrentContext,omitempty"`
	JWToken        string `protobuf:"bytes,6,opt,name=JWToken,proto3" json:"JWToken,omitempty"`
}

func (x *ContextResp) Reset() {
	*x = ContextResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContextResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContextResp) ProtoMessage() {}

func (x *ContextResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContextResp.ProtoReflect.Descriptor instead.
func (*ContextResp) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{4}
}

func (x *ContextResp) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *ContextResp) GetGuildID() string {
	if x != nil {
		return x.GuildID
	}
	return ""
}

func (x *ContextResp) GetGuildName() string {
	if x != nil {
		return x.GuildName
	}
	return ""
}

func (x *ContextResp) GetServerURL() string {
	if x != nil {
		return x.ServerURL
	}
	return ""
}

func (x *ContextResp) GetCurrentContext() bool {
	if x != nil {
		return x.CurrentContext
	}
	return false
}

func (x *ContextResp) GetJWToken() string {
	if x != nil {
		return x.JWToken
	}
	return ""
}

type Wrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token   string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	GuildID string `protobuf:"bytes,2,opt,name=guildID,proto3" json:"guildID,omitempty"`
	// Types that are assignable to Payload:
	//
	//	*Wrapper_BuildImage
	Payload isWrapper_Payload `protobuf_oneof:"Payload"`
}

func (x *Wrapper) Reset() {
	*x = Wrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Wrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Wrapper) ProtoMessage() {}

func (x *Wrapper) ProtoReflect() protoreflect.Message {
	mi := &file_proto_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Wrapper.ProtoReflect.Descriptor instead.
func (*Wrapper) Descriptor() ([]byte, []int) {
	return file_proto_api_proto_rawDescGZIP(), []int{5}
}

func (x *Wrapper) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Wrapper) GetGuildID() string {
	if x != nil {
		return x.GuildID
	}
	return ""
}

func (m *Wrapper) GetPayload() isWrapper_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *Wrapper) GetBuildImage() *BuildFunc {
	if x, ok := x.GetPayload().(*Wrapper_BuildImage); ok {
		return x.BuildImage
	}
	return nil
}

type isWrapper_Payload interface {
	isWrapper_Payload()
}

type Wrapper_BuildImage struct {
	BuildImage *BuildFunc `protobuf:"bytes,3,opt,name=BuildImage,proto3,oneof"`
}

func (*Wrapper_BuildImage) isWrapper_Payload() {}

var File_proto_api_proto protoreflect.FileDescriptor

var file_proto_api_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3c, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x44,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x44, 0x22,
	0x58, 0x0a, 0x04, 0x41, 0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x5b, 0x0a, 0x08, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x04, 0x41,
	0x72, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x41, 0x72, 0x67, 0x73,
	0x52, 0x04, 0x41, 0x72, 0x67, 0x73, 0x22, 0xb6, 0x01, 0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x46, 0x75, 0x6e, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x52, 0x75, 0x6e, 0x74, 0x69,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x44, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x08, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x52, 0x08, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x22,
	0xc1, 0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12,
	0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x47,
	0x75, 0x69, 0x6c, 0x64, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x47, 0x75,
	0x69, 0x6c, 0x64, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x47, 0x75, 0x69, 0x6c, 0x64, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x55, 0x52, 0x4c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x55, 0x52,
	0x4c, 0x12, 0x26, 0x0a, 0x0e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x43, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4a, 0x57, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4a, 0x57, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x72, 0x0a, 0x07, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x44, 0x12, 0x2c,
	0x0a, 0x0a, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x46, 0x75, 0x6e, 0x63, 0x48, 0x00,
	0x52, 0x0a, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x42, 0x09, 0x0a, 0x07,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x34, 0x74, 0x61, 0x6c, 0x2d, 0x65, 0x72, 0x72, 0x30,
	0x72, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x66, 0x61, 0x61, 0x73, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_api_proto_rawDescOnce sync.Once
	file_proto_api_proto_rawDescData = file_proto_api_proto_rawDesc
)

func file_proto_api_proto_rawDescGZIP() []byte {
	file_proto_api_proto_rawDescOnce.Do(func() {
		file_proto_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_api_proto_rawDescData)
	})
	return file_proto_api_proto_rawDescData
}

var file_proto_api_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_api_proto_goTypes = []any{
	(*GetContext)(nil),  // 0: GetContext
	(*Args)(nil),        // 1: Args
	(*Commands)(nil),    // 2: Commands
	(*BuildFunc)(nil),   // 3: BuildFunc
	(*ContextResp)(nil), // 4: ContextResp
	(*Wrapper)(nil),     // 5: Wrapper
}
var file_proto_api_proto_depIdxs = []int32{
	1, // 0: Commands.Args:type_name -> Args
	2, // 1: BuildFunc.Commands:type_name -> Commands
	3, // 2: Wrapper.BuildImage:type_name -> BuildFunc
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_api_proto_init() }
func file_proto_api_proto_init() {
	if File_proto_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_api_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GetContext); i {
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
		file_proto_api_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Args); i {
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
		file_proto_api_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Commands); i {
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
		file_proto_api_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*BuildFunc); i {
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
		file_proto_api_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ContextResp); i {
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
		file_proto_api_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*Wrapper); i {
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
	file_proto_api_proto_msgTypes[5].OneofWrappers = []any{
		(*Wrapper_BuildImage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_api_proto_goTypes,
		DependencyIndexes: file_proto_api_proto_depIdxs,
		MessageInfos:      file_proto_api_proto_msgTypes,
	}.Build()
	File_proto_api_proto = out.File
	file_proto_api_proto_rawDesc = nil
	file_proto_api_proto_goTypes = nil
	file_proto_api_proto_depIdxs = nil
}
