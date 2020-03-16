// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.1
// 	protoc        v3.6.1
// source: itsHim.proto

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

type ItsHimMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RightNodeIpAddress string `protobuf:"bytes,1,opt,name=rightNodeIpAddress,proto3" json:"rightNodeIpAddress,omitempty"`
}

func (x *ItsHimMessage) Reset() {
	*x = ItsHimMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_itsHim_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ItsHimMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ItsHimMessage) ProtoMessage() {}

func (x *ItsHimMessage) ProtoReflect() protoreflect.Message {
	mi := &file_itsHim_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ItsHimMessage.ProtoReflect.Descriptor instead.
func (*ItsHimMessage) Descriptor() ([]byte, []int) {
	return file_itsHim_proto_rawDescGZIP(), []int{0}
}

func (x *ItsHimMessage) GetRightNodeIpAddress() string {
	if x != nil {
		return x.RightNodeIpAddress
	}
	return ""
}

var File_itsHim_proto protoreflect.FileDescriptor

var file_itsHim_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x69, 0x74, 0x73, 0x48, 0x69, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x67, 0x6f, 0x6e, 0x65, 0x73, 0x74, 0x22, 0x3f, 0x0a, 0x0d, 0x49, 0x74, 0x73, 0x48, 0x69, 0x6d,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2e, 0x0a, 0x12, 0x72, 0x69, 0x67, 0x68, 0x74,
	0x4e, 0x6f, 0x64, 0x65, 0x49, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x12, 0x72, 0x69, 0x67, 0x68, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x70,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x08, 0x5a, 0x06, 0x67, 0x6f, 0x6e, 0x65, 0x73,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_itsHim_proto_rawDescOnce sync.Once
	file_itsHim_proto_rawDescData = file_itsHim_proto_rawDesc
)

func file_itsHim_proto_rawDescGZIP() []byte {
	file_itsHim_proto_rawDescOnce.Do(func() {
		file_itsHim_proto_rawDescData = protoimpl.X.CompressGZIP(file_itsHim_proto_rawDescData)
	})
	return file_itsHim_proto_rawDescData
}

var file_itsHim_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_itsHim_proto_goTypes = []interface{}{
	(*ItsHimMessage)(nil), // 0: gonest.ItsHimMessage
}
var file_itsHim_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_itsHim_proto_init() }
func file_itsHim_proto_init() {
	if File_itsHim_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_itsHim_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ItsHimMessage); i {
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
			RawDescriptor: file_itsHim_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_itsHim_proto_goTypes,
		DependencyIndexes: file_itsHim_proto_depIdxs,
		MessageInfos:      file_itsHim_proto_msgTypes,
	}.Build()
	File_itsHim_proto = out.File
	file_itsHim_proto_rawDesc = nil
	file_itsHim_proto_goTypes = nil
	file_itsHim_proto_depIdxs = nil
}