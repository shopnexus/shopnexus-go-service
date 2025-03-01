// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.2
// source: refund.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	RefundMethod_PICK_UP                   RefundMethod = 1
	RefundMethod_DROP_OFF                  RefundMethod = 2
)

// Enum value maps for RefundMethod.
var (
	RefundMethod_name = map[int32]string{
		0: "REFUND_METHOD_UNSPECIFIED",
		1: "PICK_UP",
		2: "DROP_OFF",
	}
	RefundMethod_value = map[string]int32{
		"REFUND_METHOD_UNSPECIFIED": 0,
		"PICK_UP":                   1,
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

	RefundId        string       `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
	ProductSerialId string       `protobuf:"bytes,2,opt,name=product_serial_id,json=productSerialId,proto3" json:"product_serial_id,omitempty"`
	Description     string       `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Status          Status       `protobuf:"varint,4,opt,name=status,proto3,enum=shopnexus.Status" json:"status,omitempty"`
	RefundMethod    RefundMethod `protobuf:"varint,5,opt,name=refund_method,json=refundMethod,proto3,enum=shopnexus.RefundMethod" json:"refund_method,omitempty"`
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

func (x *GetRefundResponse) GetProductSerialId() string {
	if x != nil {
		return x.ProductSerialId
	}
	return ""
}

func (x *GetRefundResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *GetRefundResponse) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_STATUS_UNSPECIFIED
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

	PaymentId int64        `protobuf:"varint,1,opt,name=payment_id,json=paymentId,proto3" json:"payment_id,omitempty"`
	Method    RefundMethod `protobuf:"varint,2,opt,name=method,proto3,enum=shopnexus.RefundMethod" json:"method,omitempty"`
	Reason    string       `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	Address   *string      `protobuf:"bytes,4,opt,name=address,proto3,oneof" json:"address,omitempty"`
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

func (x *CreateRefundRequest) GetPaymentId() int64 {
	if x != nil {
		return x.PaymentId
	}
	return 0
}

func (x *CreateRefundRequest) GetMethod() RefundMethod {
	if x != nil {
		return x.Method
	}
	return RefundMethod_REFUND_METHOD_UNSPECIFIED
}

func (x *CreateRefundRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *CreateRefundRequest) GetAddress() string {
	if x != nil && x.Address != nil {
		return *x.Address
	}
	return ""
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

type CancelRefundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefundId string `protobuf:"bytes,1,opt,name=refund_id,json=refundId,proto3" json:"refund_id,omitempty"`
}

func (x *CancelRefundRequest) Reset() {
	*x = CancelRefundRequest{}
	mi := &file_refund_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CancelRefundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelRefundRequest) ProtoMessage() {}

func (x *CancelRefundRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use CancelRefundRequest.ProtoReflect.Descriptor instead.
func (*CancelRefundRequest) Descriptor() ([]byte, []int) {
	return file_refund_proto_rawDescGZIP(), []int{5}
}

func (x *CancelRefundRequest) GetRefundId() string {
	if x != nil {
		return x.RefundId
	}
	return ""
}

var File_refund_proto protoreflect.FileDescriptor

var file_refund_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x52, 0x65, 0x66, 0x75, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x66,
	0x75, 0x6e, 0x64, 0x49, 0x64, 0x22, 0xe7, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x52, 0x65, 0x66,
	0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x72,
	0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x53, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78,
	0x75, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x3c, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e,
	0x65, 0x78, 0x75, 0x73, 0x2e, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x22,
	0xa8, 0x01, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78,
	0x75, 0x73, 0x2e, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52,
	0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12,
	0x1d, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x88, 0x01, 0x01, 0x42, 0x0a,
	0x0a, 0x08, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x33, 0x0a, 0x14, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x22,
	0xa4, 0x01, 0x0a, 0x12, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c,
	0x52, 0x06, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x73, 0x12, 0x3c, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x17, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x52, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x22, 0x32, 0x0a, 0x13, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x2a, 0x48, 0x0a, 0x0c, 0x52, 0x65,
	0x66, 0x75, 0x6e, 0x64, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1d, 0x0a, 0x19, 0x52, 0x45,
	0x46, 0x55, 0x4e, 0x44, 0x5f, 0x4d, 0x45, 0x54, 0x48, 0x4f, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x49, 0x43,
	0x4b, 0x5f, 0x55, 0x50, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x52, 0x4f, 0x50, 0x5f, 0x4f,
	0x46, 0x46, 0x10, 0x02, 0x32, 0x9f, 0x02, 0x0a, 0x06, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x12,
	0x42, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78,
	0x75, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e,
	0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x40, 0x0a, 0x05, 0x50, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1d, 0x2e, 0x73, 0x68, 0x6f, 0x70,
	0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x66, 0x75, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x42, 0x0a, 0x06, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x12, 0x1e, 0x2e, 0x73,
	0x68, 0x6f, 0x70, 0x6e, 0x65, 0x78, 0x75, 0x73, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52,
	0x65, 0x66, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x1d, 0x5a, 0x1b, 0x73, 0x68, 0x6f, 0x70, 0x6e, 0x65,
	0x78, 0x75, 0x73, 0x2d, 0x67, 0x6f, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_refund_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_refund_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_refund_proto_goTypes = []any{
	(RefundMethod)(0),            // 0: shopnexus.RefundMethod
	(*GetRefundRequest)(nil),     // 1: shopnexus.GetRefundRequest
	(*GetRefundResponse)(nil),    // 2: shopnexus.GetRefundResponse
	(*CreateRefundRequest)(nil),  // 3: shopnexus.CreateRefundRequest
	(*CreateRefundResponse)(nil), // 4: shopnexus.CreateRefundResponse
	(*PatchRefundRequest)(nil),   // 5: shopnexus.PatchRefundRequest
	(*CancelRefundRequest)(nil),  // 6: shopnexus.CancelRefundRequest
	(Status)(0),                  // 7: shopnexus.Status
	(*emptypb.Empty)(nil),        // 8: google.protobuf.Empty
}
var file_refund_proto_depIdxs = []int32{
	7, // 0: shopnexus.GetRefundResponse.status:type_name -> shopnexus.Status
	0, // 1: shopnexus.GetRefundResponse.refund_method:type_name -> shopnexus.RefundMethod
	0, // 2: shopnexus.CreateRefundRequest.method:type_name -> shopnexus.RefundMethod
	0, // 3: shopnexus.PatchRefundRequest.refund_method:type_name -> shopnexus.RefundMethod
	1, // 4: shopnexus.Refund.Get:input_type -> shopnexus.GetRefundRequest
	3, // 5: shopnexus.Refund.Create:input_type -> shopnexus.CreateRefundRequest
	5, // 6: shopnexus.Refund.Patch:input_type -> shopnexus.PatchRefundRequest
	6, // 7: shopnexus.Refund.Cancel:input_type -> shopnexus.CancelRefundRequest
	2, // 8: shopnexus.Refund.Get:output_type -> shopnexus.GetRefundResponse
	4, // 9: shopnexus.Refund.Create:output_type -> shopnexus.CreateRefundResponse
	8, // 10: shopnexus.Refund.Patch:output_type -> google.protobuf.Empty
	8, // 11: shopnexus.Refund.Cancel:output_type -> google.protobuf.Empty
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_refund_proto_init() }
func file_refund_proto_init() {
	if File_refund_proto != nil {
		return
	}
	file_status_proto_init()
	file_refund_proto_msgTypes[2].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_refund_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
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
