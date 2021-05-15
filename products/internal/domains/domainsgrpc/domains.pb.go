// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: domains.proto

package domainsgrpc

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

type DiscountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId string `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	UserId    string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *DiscountRequest) Reset() {
	*x = DiscountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domains_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscountRequest) ProtoMessage() {}

func (x *DiscountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_domains_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscountRequest.ProtoReflect.Descriptor instead.
func (*DiscountRequest) Descriptor() ([]byte, []int) {
	return file_domains_proto_rawDescGZIP(), []int{0}
}

func (x *DiscountRequest) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *DiscountRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type DiscountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId    string  `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	UserId       string  `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Percentage   float32 `protobuf:"fixed32,3,opt,name=percentage,proto3" json:"percentage,omitempty"`
	ValueInCents int64   `protobuf:"varint,4,opt,name=value_in_cents,json=valueInCents,proto3" json:"value_in_cents,omitempty"`
}

func (x *DiscountResponse) Reset() {
	*x = DiscountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domains_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscountResponse) ProtoMessage() {}

func (x *DiscountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_domains_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscountResponse.ProtoReflect.Descriptor instead.
func (*DiscountResponse) Descriptor() ([]byte, []int) {
	return file_domains_proto_rawDescGZIP(), []int{1}
}

func (x *DiscountResponse) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *DiscountResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DiscountResponse) GetPercentage() float32 {
	if x != nil {
		return x.Percentage
	}
	return 0
}

func (x *DiscountResponse) GetValueInCents() int64 {
	if x != nil {
		return x.ValueInCents
	}
	return 0
}

var File_domains_proto protoreflect.FileDescriptor

var file_domains_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x22, 0x49, 0x0a, 0x0f, 0x44, 0x69, 0x73, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x90, 0x01, 0x0a, 0x10, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65,
	0x12, 0x24, 0x0a, 0x0e, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x63, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x49,
	0x6e, 0x43, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x54, 0x48, 0x01, 0x5a, 0x50, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x6c, 0x6d, 0x61, 0x72, 0x63, 0x6f, 0x67,
	0x64, 0x2f, 0x6d, 0x6f, 0x62, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x73, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x67, 0x72, 0x70, 0x63,
	0x3b, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_domains_proto_rawDescOnce sync.Once
	file_domains_proto_rawDescData = file_domains_proto_rawDesc
)

func file_domains_proto_rawDescGZIP() []byte {
	file_domains_proto_rawDescOnce.Do(func() {
		file_domains_proto_rawDescData = protoimpl.X.CompressGZIP(file_domains_proto_rawDescData)
	})
	return file_domains_proto_rawDescData
}

var file_domains_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_domains_proto_goTypes = []interface{}{
	(*DiscountRequest)(nil),  // 0: domains.DiscountRequest
	(*DiscountResponse)(nil), // 1: domains.DiscountResponse
}
var file_domains_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_domains_proto_init() }
func file_domains_proto_init() {
	if File_domains_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_domains_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscountRequest); i {
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
		file_domains_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscountResponse); i {
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
			RawDescriptor: file_domains_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_domains_proto_goTypes,
		DependencyIndexes: file_domains_proto_depIdxs,
		MessageInfos:      file_domains_proto_msgTypes,
	}.Build()
	File_domains_proto = out.File
	file_domains_proto_rawDesc = nil
	file_domains_proto_goTypes = nil
	file_domains_proto_depIdxs = nil
}
