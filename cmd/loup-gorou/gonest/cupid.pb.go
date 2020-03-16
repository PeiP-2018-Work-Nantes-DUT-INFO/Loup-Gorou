// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.1
// 	protoc        v3.6.1
// source: cupid.proto

package gonest

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

type CupidMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IpAddress1 string `protobuf:"bytes,1,opt,name=ipAddress1,proto3" json:"ipAddress1,omitempty"`
	IpAddress2 string `protobuf:"bytes,2,opt,name=ipAddress2,proto3" json:"ipAddress2,omitempty"`
}

func (x *CupidMessage) Reset() {
	*x = CupidMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cupid_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CupidMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CupidMessage) ProtoMessage() {}

func (x *CupidMessage) ProtoReflect() protoreflect.Message {
	mi := &file_cupid_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CupidMessage.ProtoReflect.Descriptor instead.
func (*CupidMessage) Descriptor() ([]byte, []int) {
	return file_cupid_proto_rawDescGZIP(), []int{0}
}

func (x *CupidMessage) GetIpAddress1() string {
	if x != nil {
		return x.IpAddress1
	}
	return ""
}

func (x *CupidMessage) GetIpAddress2() string {
	if x != nil {
		return x.IpAddress2
	}
	return ""
}

var File_cupid_proto protoreflect.FileDescriptor

var file_cupid_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x75, 0x70, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67,
	0x6f, 0x6e, 0x65, 0x73, 0x74, 0x22, 0x4e, 0x0a, 0x0c, 0x43, 0x75, 0x70, 0x69, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x70, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x31, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x70, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x32, 0x42, 0x08, 0x5a, 0x06, 0x67, 0x6f, 0x6e, 0x65, 0x73, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cupid_proto_rawDescOnce sync.Once
	file_cupid_proto_rawDescData = file_cupid_proto_rawDesc
)

func file_cupid_proto_rawDescGZIP() []byte {
	file_cupid_proto_rawDescOnce.Do(func() {
		file_cupid_proto_rawDescData = protoimpl.X.CompressGZIP(file_cupid_proto_rawDescData)
	})
	return file_cupid_proto_rawDescData
}

var file_cupid_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_cupid_proto_goTypes = []interface{}{
	(*CupidMessage)(nil), // 0: gonest.CupidMessage
}
var file_cupid_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cupid_proto_init() }
func file_cupid_proto_init() {
	if File_cupid_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cupid_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CupidMessage); i {
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
			RawDescriptor: file_cupid_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cupid_proto_goTypes,
		DependencyIndexes: file_cupid_proto_depIdxs,
		MessageInfos:      file_cupid_proto_msgTypes,
	}.Build()
	File_cupid_proto = out.File
	file_cupid_proto_rawDesc = nil
	file_cupid_proto_goTypes = nil
	file_cupid_proto_depIdxs = nil
}
