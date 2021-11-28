// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: recordvalidator.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Scheme_Order int32

const (
	Scheme_RANDOM        Scheme_Order = 0
	Scheme_ORDER         Scheme_Order = 1
	Scheme_REVERSE_ORDER Scheme_Order = 2
)

// Enum value maps for Scheme_Order.
var (
	Scheme_Order_name = map[int32]string{
		0: "RANDOM",
		1: "ORDER",
		2: "REVERSE_ORDER",
	}
	Scheme_Order_value = map[string]int32{
		"RANDOM":        0,
		"ORDER":         1,
		"REVERSE_ORDER": 2,
	}
)

func (x Scheme_Order) Enum() *Scheme_Order {
	p := new(Scheme_Order)
	*p = x
	return p
}

func (x Scheme_Order) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Scheme_Order) Descriptor() protoreflect.EnumDescriptor {
	return file_recordvalidator_proto_enumTypes[0].Descriptor()
}

func (Scheme_Order) Type() protoreflect.EnumType {
	return &file_recordvalidator_proto_enumTypes[0]
}

func (x Scheme_Order) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Scheme_Order.Descriptor instead.
func (Scheme_Order) EnumDescriptor() ([]byte, []int) {
	return file_recordvalidator_proto_rawDescGZIP(), []int{1, 0}
}

type Schemes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schemes []*Scheme `protobuf:"bytes,1,rep,name=schemes,proto3" json:"schemes,omitempty"`
}

func (x *Schemes) Reset() {
	*x = Schemes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordvalidator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Schemes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schemes) ProtoMessage() {}

func (x *Schemes) ProtoReflect() protoreflect.Message {
	mi := &file_recordvalidator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Schemes.ProtoReflect.Descriptor instead.
func (*Schemes) Descriptor() ([]byte, []int) {
	return file_recordvalidator_proto_rawDescGZIP(), []int{0}
}

func (x *Schemes) GetSchemes() []*Scheme {
	if x != nil {
		return x.Schemes
	}
	return nil
}

type Scheme struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string          `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	StartTime    int64           `protobuf:"varint,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	InstanceIds  []int32         `protobuf:"varint,3,rep,packed,name=instance_ids,json=instanceIds,proto3" json:"instance_ids,omitempty"`
	CompletedIds []int32         `protobuf:"varint,4,rep,packed,name=completed_ids,json=completedIds,proto3" json:"completed_ids,omitempty"`
	CurrentPick  int32           `protobuf:"varint,5,opt,name=current_pick,json=currentPick,proto3" json:"current_pick,omitempty"`
	Order        Scheme_Order    `protobuf:"varint,6,opt,name=order,proto3,enum=recordvalidator.Scheme_Order" json:"order,omitempty"`
	Unbox        bool            `protobuf:"varint,7,opt,name=unbox,proto3" json:"unbox,omitempty"`
	Soft         bool            `protobuf:"varint,8,opt,name=soft,proto3" json:"soft,omitempty"`
	CompleteDate map[int32]int64 `protobuf:"bytes,9,rep,name=complete_date,json=completeDate,proto3" json:"complete_date,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *Scheme) Reset() {
	*x = Scheme{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordvalidator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Scheme) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Scheme) ProtoMessage() {}

func (x *Scheme) ProtoReflect() protoreflect.Message {
	mi := &file_recordvalidator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Scheme.ProtoReflect.Descriptor instead.
func (*Scheme) Descriptor() ([]byte, []int) {
	return file_recordvalidator_proto_rawDescGZIP(), []int{1}
}

func (x *Scheme) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Scheme) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *Scheme) GetInstanceIds() []int32 {
	if x != nil {
		return x.InstanceIds
	}
	return nil
}

func (x *Scheme) GetCompletedIds() []int32 {
	if x != nil {
		return x.CompletedIds
	}
	return nil
}

func (x *Scheme) GetCurrentPick() int32 {
	if x != nil {
		return x.CurrentPick
	}
	return 0
}

func (x *Scheme) GetOrder() Scheme_Order {
	if x != nil {
		return x.Order
	}
	return Scheme_RANDOM
}

func (x *Scheme) GetUnbox() bool {
	if x != nil {
		return x.Unbox
	}
	return false
}

func (x *Scheme) GetSoft() bool {
	if x != nil {
		return x.Soft
	}
	return false
}

func (x *Scheme) GetCompleteDate() map[int32]int64 {
	if x != nil {
		return x.CompleteDate
	}
	return nil
}

type GetSchemeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetSchemeRequest) Reset() {
	*x = GetSchemeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordvalidator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSchemeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSchemeRequest) ProtoMessage() {}

func (x *GetSchemeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recordvalidator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSchemeRequest.ProtoReflect.Descriptor instead.
func (*GetSchemeRequest) Descriptor() ([]byte, []int) {
	return file_recordvalidator_proto_rawDescGZIP(), []int{2}
}

func (x *GetSchemeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetSchemeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scheme *Scheme `protobuf:"bytes,1,opt,name=scheme,proto3" json:"scheme,omitempty"`
}

func (x *GetSchemeResponse) Reset() {
	*x = GetSchemeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recordvalidator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSchemeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSchemeResponse) ProtoMessage() {}

func (x *GetSchemeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recordvalidator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSchemeResponse.ProtoReflect.Descriptor instead.
func (*GetSchemeResponse) Descriptor() ([]byte, []int) {
	return file_recordvalidator_proto_rawDescGZIP(), []int{3}
}

func (x *GetSchemeResponse) GetScheme() *Scheme {
	if x != nil {
		return x.Scheme
	}
	return nil
}

var File_recordvalidator_proto protoreflect.FileDescriptor

var file_recordvalidator_proto_rawDesc = []byte{
	0x0a, 0x15, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x3c, 0x0a, 0x07, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x52, 0x07, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x65, 0x73, 0x22, 0xc9, 0x03, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0b, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0c,
	0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x49, 0x64, 0x73, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x69, 0x63, 0x6b, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x69, 0x63, 0x6b, 0x12,
	0x33, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d,
	0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72,
	0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x05, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x6e, 0x62, 0x6f, 0x78, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x75, 0x6e, 0x62, 0x6f, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f,
	0x66, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x73, 0x6f, 0x66, 0x74, 0x12, 0x4e,
	0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x2e, 0x43,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0c, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x1a, 0x3f,
	0x0a, 0x11, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x65, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x31, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x41, 0x4e, 0x44,
	0x4f, 0x4d, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x10, 0x01, 0x12,
	0x11, 0x0a, 0x0d, 0x52, 0x45, 0x56, 0x45, 0x52, 0x53, 0x45, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52,
	0x10, 0x02, 0x22, 0x26, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x44, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2f, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65,
	0x32, 0x6e, 0x0a, 0x16, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x09, 0x47, 0x65,
	0x74, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x21, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x72, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2f, 0x72, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_recordvalidator_proto_rawDescOnce sync.Once
	file_recordvalidator_proto_rawDescData = file_recordvalidator_proto_rawDesc
)

func file_recordvalidator_proto_rawDescGZIP() []byte {
	file_recordvalidator_proto_rawDescOnce.Do(func() {
		file_recordvalidator_proto_rawDescData = protoimpl.X.CompressGZIP(file_recordvalidator_proto_rawDescData)
	})
	return file_recordvalidator_proto_rawDescData
}

var file_recordvalidator_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_recordvalidator_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_recordvalidator_proto_goTypes = []interface{}{
	(Scheme_Order)(0),         // 0: recordvalidator.Scheme.Order
	(*Schemes)(nil),           // 1: recordvalidator.Schemes
	(*Scheme)(nil),            // 2: recordvalidator.Scheme
	(*GetSchemeRequest)(nil),  // 3: recordvalidator.GetSchemeRequest
	(*GetSchemeResponse)(nil), // 4: recordvalidator.GetSchemeResponse
	nil,                       // 5: recordvalidator.Scheme.CompleteDateEntry
}
var file_recordvalidator_proto_depIdxs = []int32{
	2, // 0: recordvalidator.Schemes.schemes:type_name -> recordvalidator.Scheme
	0, // 1: recordvalidator.Scheme.order:type_name -> recordvalidator.Scheme.Order
	5, // 2: recordvalidator.Scheme.complete_date:type_name -> recordvalidator.Scheme.CompleteDateEntry
	2, // 3: recordvalidator.GetSchemeResponse.scheme:type_name -> recordvalidator.Scheme
	3, // 4: recordvalidator.RecordValidatorService.GetScheme:input_type -> recordvalidator.GetSchemeRequest
	4, // 5: recordvalidator.RecordValidatorService.GetScheme:output_type -> recordvalidator.GetSchemeResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_recordvalidator_proto_init() }
func file_recordvalidator_proto_init() {
	if File_recordvalidator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_recordvalidator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Schemes); i {
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
		file_recordvalidator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Scheme); i {
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
		file_recordvalidator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSchemeRequest); i {
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
		file_recordvalidator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSchemeResponse); i {
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
			RawDescriptor: file_recordvalidator_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_recordvalidator_proto_goTypes,
		DependencyIndexes: file_recordvalidator_proto_depIdxs,
		EnumInfos:         file_recordvalidator_proto_enumTypes,
		MessageInfos:      file_recordvalidator_proto_msgTypes,
	}.Build()
	File_recordvalidator_proto = out.File
	file_recordvalidator_proto_rawDesc = nil
	file_recordvalidator_proto_goTypes = nil
	file_recordvalidator_proto_depIdxs = nil
}
