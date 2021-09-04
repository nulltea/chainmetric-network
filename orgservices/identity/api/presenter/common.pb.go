// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: identity/api/presenter/common.proto

package presenter

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Status int32

const (
	Status_UNKNOWN    Status = 0
	Status_OK         Status = 200
	Status_ACCEPTED   Status = 202
	Status_NO_CONTENT Status = 204
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0:   "UNKNOWN",
		200: "OK",
		202: "ACCEPTED",
		204: "NO_CONTENT",
	}
	Status_value = map[string]int32{
		"UNKNOWN":    0,
		"OK":         200,
		"ACCEPTED":   202,
		"NO_CONTENT": 204,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_identity_api_presenter_common_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_identity_api_presenter_common_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_identity_api_presenter_common_proto_rawDescGZIP(), []int{0}
}

type StatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status Status `protobuf:"varint,1,opt,name=status,proto3,enum=chainmetric.identity.presenter.Status" json:"status,omitempty"`
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_identity_api_presenter_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_api_presenter_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_identity_api_presenter_common_proto_rawDescGZIP(), []int{0}
}

func (x *StatusResponse) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

var File_identity_api_presenter_common_proto protoreflect.FileDescriptor

var file_identity_api_presenter_common_proto_rawDesc = []byte{
	0x0a, 0x23, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x65, 0x73,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x22, 0x50, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x3e, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x07,
	0x0a, 0x02, 0x4f, 0x4b, 0x10, 0xc8, 0x01, 0x12, 0x0d, 0x0a, 0x08, 0x41, 0x43, 0x43, 0x45, 0x50,
	0x54, 0x45, 0x44, 0x10, 0xca, 0x01, 0x12, 0x0f, 0x0a, 0x0a, 0x4e, 0x4f, 0x5f, 0x43, 0x4f, 0x4e,
	0x54, 0x45, 0x4e, 0x54, 0x10, 0xcc, 0x01, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6d, 0x6f, 0x74, 0x68, 0x2d, 0x79, 0x2f, 0x63,
	0x68, 0x61, 0x69, 0x6e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x73, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_identity_api_presenter_common_proto_rawDescOnce sync.Once
	file_identity_api_presenter_common_proto_rawDescData = file_identity_api_presenter_common_proto_rawDesc
)

func file_identity_api_presenter_common_proto_rawDescGZIP() []byte {
	file_identity_api_presenter_common_proto_rawDescOnce.Do(func() {
		file_identity_api_presenter_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_identity_api_presenter_common_proto_rawDescData)
	})
	return file_identity_api_presenter_common_proto_rawDescData
}

var file_identity_api_presenter_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_identity_api_presenter_common_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_identity_api_presenter_common_proto_goTypes = []interface{}{
	(Status)(0),            // 0: chainmetric.identity.presenter.Status
	(*StatusResponse)(nil), // 1: chainmetric.identity.presenter.StatusResponse
}
var file_identity_api_presenter_common_proto_depIdxs = []int32{
	0, // 0: chainmetric.identity.presenter.StatusResponse.status:type_name -> chainmetric.identity.presenter.Status
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_identity_api_presenter_common_proto_init() }
func file_identity_api_presenter_common_proto_init() {
	if File_identity_api_presenter_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_identity_api_presenter_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse); i {
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
			RawDescriptor: file_identity_api_presenter_common_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_identity_api_presenter_common_proto_goTypes,
		DependencyIndexes: file_identity_api_presenter_common_proto_depIdxs,
		EnumInfos:         file_identity_api_presenter_common_proto_enumTypes,
		MessageInfos:      file_identity_api_presenter_common_proto_msgTypes,
	}.Build()
	File_identity_api_presenter_common_proto = out.File
	file_identity_api_presenter_common_proto_rawDesc = nil
	file_identity_api_presenter_common_proto_goTypes = nil
	file_identity_api_presenter_common_proto_depIdxs = nil
}