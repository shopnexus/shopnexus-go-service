// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v4.23.1
// source: refund.proto

package pb

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

type RefundMethod int32

const (
	RefundMethod_REFUND_METHOD_UNSPECIFIED RefundMethod = 0
	RefundMethod_PICKUP                    RefundMethod = 1
	RefundMethod_DROP_OFF                  RefundMethod = 2
)

// Enum value maps for RefundMethod.
var (
	RefundMethod_name = map[int32]string{
		0: "REFUND_METHOD_UNSPECIFIED",
		1: "PICKUP",
		2: "DROP_OFF",
	}
	RefundMethod_value = map[string]int32{
		"REFUND_METHOD_UNSPECIFIED": 0,
		"PICKUP":                    1,
		"DROP_OFF":                  2,
	}
)

func (x RefundMethod) Enum() *RefundMethod {
	p := new(RefundMethod)
	*p = x
	return p
}

func (x RefundMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RefundMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_refund_proto_enumTypes[0].Descriptor()
}

func (RefundMethod) Type() protoreflect.EnumType {
	return &file_refund_proto_enumTypes[0]
}

func (x RefundMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RefundMethod.Descriptor instead.
func (RefundMethod) EnumDescriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{0}
}

type RefundStatus int32

const (
	RefundStatus_REFUND_STATUS_UNSPECIFIED RefundStatus = 0
	RefundStatus_PENDING                   RefundStatus = 1
	RefundStatus_ACCEPTED                  RefundStatus = 2
	RefundStatus_REJECTED                  RefundStatus = 3
	RefundStatus_CANCELLED                 RefundStatus = 4
)

// Enum value maps for RefundStatus.
var (
	RefundStatus_name = map[int32]string{
		0: "REFUND_STATUS_UNSPECIFIED",
		1: "PENDING",
		2: "ACCEPTED",
		3: "REJECTED",
		4: "CANCELLED",
	}
	RefundStatus_value = map[string]int32{
		"REFUND_STATUS_UNSPECIFIED": 0,
		"PENDING":                   1,
		"ACCEPTED":                  2,
		"REJECTED":                  3,
		"CANCELLED":                 4,
	}
)

func (x RefundStatus) Enum() *RefundStatus {
	p := new(RefundStatus)
	*p = x
	return p
}

func (x RefundStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RefundStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_refund_proto_enumTypes[1].Descriptor()
}

func (RefundStatus) Type() protoreflect.EnumType {
	return &file_refund_proto_enumTypes[1]
}

func (x RefundStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RefundStatus.Descriptor instead.
func (RefundStatus) EnumDescriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{1}
}

type CreateRefundRequest_RefundMethod int32

const (
	CreateRefundRequest_REFUND_METHOD_UNSPECIFIED CreateRefundRequest_RefundMethod = 0
	CreateRefundRequest_PICKUP                    CreateRefundRequest_RefundMethod = 1
	CreateRefundRequest_DROP_OFF                  CreateRefundRequest_RefundMethod = 2
)

// Enum value maps for CreateRefundRequest_RefundMethod.
var (
	CreateRefundRequest_RefundMethod_name = map[int32]string{
		0: "REFUND_METHOD_UNSPECIFIED",
		1: "PICKUP",
		2: "DROP_OFF",
	}
	CreateRefundRequest_RefundMethod_value = map[string]int32{
		"REFUND_METHOD_UNSPECIFIED": 0,
		"PICKUP":                    1,
		"DROP_OFF":                  2,
	}
)

func (x CreateRefundRequest_RefundMethod) Enum() *CreateRefundRequest_RefundMethod {
	p := new(CreateRefundRequest_RefundMethod)
	*p = x
	return p
}

func (x CreateRefundRequest_RefundMethod) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CreateRefundRequest_RefundMethod) Descriptor() protoreflect.EnumDescriptor {
	return file_refund_proto_enumTypes[2].Descriptor()
}

func (CreateRefundRequest_RefundMethod) Type() protoreflect.EnumType {
	return &file_refund_proto_enumTypes[2]
}

func (x CreateRefundRequest_RefundMethod) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CreateRefundRequest_RefundMethod.Descriptor instead.
func (CreateRefundRequest_RefundMethod) EnumDescriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{2, 0}
}

type GetRefundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefundId string `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
}

func (x *GetRefundRequest) Reset() {
	*x = GetRefundRequest{}
	mi := &file_refund_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRefundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRefundRequest) ProtoMessage() {}

func (x *GetRefundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRefundRequest.ProtoReflect.Descriptor instead.
func (*GetRefundRequest) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{0}
}

func (x *GetRefundRequest) GetRefundId() string {
	if x != nil {
		return x.RefundId
	}
	return ""
}

type GetRefundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefundId     string       `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
	ProductId    []byte       `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Description  string       `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Images       [][]byte     `protobuf:"bytes,4,rep,name=images,proto3" json:"images,omitempty"`
	Videos       [][]byte     `protobuf:"bytes,5,rep,name=videos,proto3" json:"videos,omitempty"`
	Status       RefundStatus `protobuf:"varint,6,opt,name=status,proto3,enum=shopnexus.RefundStatus" json:"status,omitempty"`
	RefundMethod RefundMethod `protobuf:"varint,7,opt,name=refund_method,json=refundMethod,proto3,enum=shopnexus.RefundMethod" json:"refund_method,omitempty"`
}

func (x *GetRefundResponse) Reset() {
	*x = GetRefundResponse{}
	mi := &file_refund_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRefundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRefundResponse) ProtoMessage() {}

func (x *GetRefundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRefundResponse.ProtoReflect.Descriptor instead.
func (*GetRefundResponse) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{1}
}

func (x *GetRefundResponse) GetRefundId() string {
	if x != nil {
		return x.RefundId
	}
	return ""
}

func (x *GetRefundResponse) GetProductId() []byte {
	if x != nil {
		return x.ProductId
	}
	return nil
}

func (x *GetRefundResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *GetRefundResponse) GetImages() [][]byte {
	if x != nil {
		return x.Images
	}
	return nil
}

func (x *GetRefundResponse) GetVideos() [][]byte {
	if x != nil {
		return x.Videos
	}
	return nil
}

func (x *GetRefundResponse) GetStatus() RefundStatus {
	if x != nil {
		return x.Status
	}
	return RefundStatus_REFUND_STATUS_UNSPECIFIED
}

func (x *GetRefundResponse) GetRefundMethod() RefundMethod {
	if x != nil {
		return x.RefundMethod
	}
	return RefundMethod_REFUND_METHOD_UNSPECIFIED
}

type CreateRefundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId    []byte                           `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Description  string                           `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Images       [][]byte                         `protobuf:"bytes,3,rep,name=images,proto3" json:"images,omitempty"`
	Videos       [][]byte                         `protobuf:"bytes,4,rep,name=videos,proto3" json:"videos,omitempty"`
	RefundMethod CreateRefundRequest_RefundMethod `protobuf:"varint,5,opt,name=refund_method,json=refundMethod,proto3,enum=shopnexus.CreateRefundRequest_RefundMethod" json:"refund_method,omitempty"`
}

func (x *CreateRefundRequest) Reset() {
	*x = CreateRefundRequest{}
	mi := &file_refund_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRefundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRefundRequest) ProtoMessage() {}

func (x *CreateRefundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRefundRequest.ProtoReflect.Descriptor instead.
func (*CreateRefundRequest) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{2}
}

func (x *CreateRefundRequest) GetProductId() []byte {
	if x != nil {
		return x.ProductId
	}
	return nil
}

func (x *CreateRefundRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateRefundRequest) GetImages() [][]byte {
	if x != nil {
		return x.Images
	}
	return nil
}

func (x *CreateRefundRequest) GetVideos() [][]byte {
	if x != nil {
		return x.Videos
	}
	return nil
}

func (x *CreateRefundRequest) GetRefundMethod() CreateRefundRequest_RefundMethod {
	if x != nil {
		return x.RefundMethod
	}
	return CreateRefundRequest_REFUND_METHOD_UNSPECIFIED
}

type CreateRefundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefundId string `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
}

func (x *CreateRefundResponse) Reset() {
	*x = CreateRefundResponse{}
	mi := &file_refund_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRefundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRefundResponse) ProtoMessage() {}

func (x *CreateRefundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRefundResponse.ProtoReflect.Descriptor instead.
func (*CreateRefundResponse) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{3}
}

func (x *CreateRefundResponse) GetRefundId() string {
	if x != nil {
		return x.RefundId
	}
	return ""
}

type PatchRefundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// string refund_id = 1;
	Description  string       `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Images       [][]byte     `protobuf:"bytes,2,rep,name=images,proto3" json:"images,omitempty"`
	Videos       [][]byte     `protobuf:"bytes,3,rep,name=videos,proto3" json:"videos,omitempty"`
	RefundMethod RefundMethod `protobuf:"varint,4,opt,name=refund_method,json=refundMethod,proto3,enum=shopnexus.RefundMethod" json:"refund_method,omitempty"`
}

func (x *PatchRefundRequest) Reset() {
	*x = PatchRefundRequest{}
	mi := &file_refund_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchRefundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchRefundRequest) ProtoMessage() {}

func (x *PatchRefundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchRefundRequest.ProtoReflect.Descriptor instead.
func (*PatchRefundRequest) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{4}
}

func (x *PatchRefundRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PatchRefundRequest) GetImages() [][]byte {
	if x != nil {
		return x.Images
	}
	return nil
}

func (x *PatchRefundRequest) GetVideos() [][]byte {
	if x != nil {
		return x.Videos
	}
	return nil
}

func (x *PatchRefundRequest) GetRefundMethod() RefundMethod {
	if x != nil {
		return x.RefundMethod
	}
	return RefundMethod_REFUND_METHOD_UNSPECIFIED
}

type PatchRefundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefundId string `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
}

func (x *PatchRefundResponse) Reset() {
	*x = PatchRefundResponse{}
	mi := &file_refund_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchRefundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchRefundResponse) ProtoMessage() {}

func (x *PatchRefundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchRefundResponse.ProtoReflect.Descriptor instead.
func (*PatchRefundResponse) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{5}
}

func (x *PatchRefundResponse) GetRefundId() string {
	if x != nil {
		return x.RefundId
	}
	return ""
}

type CancelRefundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefundId string `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
}

func (x *CancelRefundRequest) Reset() {
	*x = CancelRefundRequest{}
	mi := &file_refund_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CancelRefundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelRefundRequest) ProtoMessage() {}

func (x *CancelRefundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelRefundRequest.ProtoReflect.Descriptor instead.
func (*CancelRefundRequest) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{6}
}

func (x *CancelRefundRequest) GetRefundId() string {
	if x != nil {
		return x.RefundId
	}
	return ""
}

type CancelRefundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefundId string `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
}

func (x *CancelRefundResponse) Reset() {
	*x = CancelRefundResponse{}
	mi := &file_refund_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CancelRefundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelRefundResponse) ProtoMessage() {}

func (x *CancelRefundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_refund_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelRefundResponse.ProtoReflect.Descriptor instead.
func (*CancelRefundResponse) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{7}
}

func (x *CancelRefundResponse) GetRefundId() string {
	if x != nil {
		return x.RefundId
	}
	return ""
}

var File_refund_proto protoreflect.FileDescriptor

var file_refund_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x22, 0x2f, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x22, 0x90, 0x02, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x12, 0x2f,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17,
	0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x52, 0x65, 0x66, 0x75, 0x6e,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x3c, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78,
	0x75, 0x73, 0x2e, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x0c, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x22, 0xa1, 0x02,
	0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06,
	0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x12, 0x50, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64,
	0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2b, 0x2e,
	0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65,
	0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x22, 0x47, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1d, 0x0a, 0x19, 0x52, 0x45, 0x46, 0x55,
	0x4e, 0x44, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x49, 0x43, 0x4b, 0x55,
	0x50, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x52, 0x4f, 0x50, 0x5f, 0x4f, 0x46, 0x46, 0x10,
	0x02, 0x22, 0x33, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x75, 0x6e,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x66,
	0x75, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x66, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x22, 0xa4, 0x01, 0x0a, 0x12, 0x50, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x16, 0x0a, 0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52,
	0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x12,
	0x3c, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78,
	0x75, 0x73, 0x2e, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x0c, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x22, 0x32, 0x0a,
	0x13, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49,
	0x64, 0x22, 0x32, 0x0a, 0x13, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52, 0x65, 0x66, 0x75, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x66,
	0x75, 0x6e, 0x64, 0x49, 0x64, 0x22, 0x33, 0x0a, 0x14, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52,
	0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x2a, 0x47, 0x0a, 0x0c, 0x52, 0x65,
	0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1d, 0x0a, 0x19, 0x52, 0x45,
	0x46, 0x55, 0x4e, 0x44, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x49, 0x43,
	0x4b, 0x55, 0x50, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x52, 0x4f, 0x50, 0x5f, 0x4f, 0x46,
	0x46, 0x10, 0x02, 0x2a, 0x65, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a, 0x19, 0x52, 0x45, 0x46, 0x55, 0x4e, 0x44, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12,
	0x0c, 0x0a, 0x08, 0x41, 0x43, 0x43, 0x45, 0x50, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0c, 0x0a,
	0x08, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x43,
	0x41, 0x4e, 0x43, 0x45, 0x4c, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x32, 0xb0, 0x02, 0x0a, 0x06, 0x52,
	0x65, 0x66, 0x75, 0x6e, 0x64, 0x12, 0x42, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x73,
	0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x73, 0x68, 0x6f, 0x70,
	0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x06, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x05, 0x50, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x1d, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x50, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x50, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x4b, 0x0a, 0x06, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x12, 0x1e, 0x2e, 0x73, 0x68, 0x6f,
	0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52, 0x65, 0x66,
	0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73, 0x68, 0x6f,
	0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52, 0x65, 0x66,
	0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1d, 0x5a,
	0x1b, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2d, 0x67, 0x6f, 0x2d, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_refund_proto_rawDescOnce sync.Once
	file_refund_proto_rawDescData = file_refund_proto_rawDesc
)

func file_refund_proto_rawDescGZIP() []byte {
	file_refund_proto_rawDescOnce.Do(func() {
		file_refund_proto_rawDescData = protoimpl.X.CompressGZIP(file_refund_proto_rawDescData)
	})
	return file_refund_proto_rawDescData
}

var file_refund_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_refund_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_refund_proto_goTypes = []any{
	(RefundMethod)(0),                     // 0: shopnexus.RefundMethod
	(RefundStatus)(0),                     // 1: shopnexus.RefundStatus
	(CreateRefundRequest_RefundMethod)(0), // 2: shopnexus.CreateRefundRequest.RefundMethod
	(*GetRefundRequest)(nil),              // 3: shopnexus.GetRefundRequest
	(*GetRefundResponse)(nil),             // 4: shopnexus.GetRefundResponse
	(*CreateRefundRequest)(nil),           // 5: shopnexus.CreateRefundRequest
	(*CreateRefundResponse)(nil),          // 6: shopnexus.CreateRefundResponse
	(*PatchRefundRequest)(nil),            // 7: shopnexus.PatchRefundRequest
	(*PatchRefundResponse)(nil),           // 8: shopnexus.PatchRefundResponse
	(*CancelRefundRequest)(nil),           // 9: shopnexus.CancelRefundRequest
	(*CancelRefundResponse)(nil),          // 10: shopnexus.CancelRefundResponse
}
var file_refund_proto_depIdxs = []int32{
	1,  // 0: shopnexus.GetRefundResponse.status:type_name -> shopnexus.RefundStatus
	0,  // 1: shopnexus.GetRefundResponse.refund_method:type_name -> shopnexus.RefundMethod
	2,  // 2: shopnexus.CreateRefundRequest.refund_method:type_name -> shopnexus.CreateRefundRequest.RefundMethod
	0,  // 3: shopnexus.PatchRefundRequest.refund_method:type_name -> shopnexus.RefundMethod
	3,  // 4: shopnexus.Refund.Get:input_type -> shopnexus.GetRefundRequest
	5,  // 5: shopnexus.Refund.Create:input_type -> shopnexus.CreateRefundRequest
	7,  // 6: shopnexus.Refund.Patch:input_type -> shopnexus.PatchRefundRequest
	9,  // 7: shopnexus.Refund.Cancel:input_type -> shopnexus.CancelRefundRequest
	4,  // 8: shopnexus.Refund.Get:output_type -> shopnexus.GetRefundResponse
	6,  // 9: shopnexus.Refund.Create:output_type -> shopnexus.CreateRefundResponse
	8,  // 10: shopnexus.Refund.Patch:output_type -> shopnexus.PatchRefundResponse
	10, // 11: shopnexus.Refund.Cancel:output_type -> shopnexus.CancelRefundResponse
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_refund_proto_init() }
func file_refund_proto_init() {
	if File_refund_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_refund_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_refund_proto_goTypes,
		DependencyIndexes: file_refund_proto_depIdxs,
		EnumInfos:         file_refund_proto_enumTypes,
		MessageInfos:      file_refund_proto_msgTypes,
	}.Build()
	File_refund_proto = out.File
	file_refund_proto_rawDesc = nil
	file_refund_proto_goTypes = nil
	file_refund_proto_depIdxs = nil
}
