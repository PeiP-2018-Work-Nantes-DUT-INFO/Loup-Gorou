// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.1
// 	protoc        v3.6.1
// source: leaderelection.proto

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

//Message envoyé par le leader, pour assurer que tout le monde est le meme leader.
type LeaderElectionMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Leader string `protobuf:"bytes,1,opt,name=leader,proto3" json:"leader,omitempty"`
}

func (x *LeaderElectionMessage) Reset() {
	*x = LeaderElectionMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_leaderelection_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaderElectionMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaderElectionMessage) ProtoMessage() {}

func (x *LeaderElectionMessage) ProtoReflect() protoreflect.Message {
	mi := &file_leaderelection_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaderElectionMessage.ProtoReflect.Descriptor instead.
func (*LeaderElectionMessage) Descriptor() ([]byte, []int) {
	return file_leaderelection_proto_rawDescGZIP(), []int{0}
}

func (x *LeaderElectionMessage) GetLeader() string {
	if x != nil {
		return x.Leader
	}
	return ""
}

var File_leaderelection_proto protoreflect.FileDescriptor

var file_leaderelection_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x6f, 0x6e, 0x65, 0x73, 0x74, 0x22, 0x2f,
	0x0a, 0x15, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x42,
	0x08, 0x5a, 0x06, 0x67, 0x6f, 0x6e, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_leaderelection_proto_rawDescOnce sync.Once
	file_leaderelection_proto_rawDescData = file_leaderelection_proto_rawDesc
)

func file_leaderelection_proto_rawDescGZIP() []byte {
	file_leaderelection_proto_rawDescOnce.Do(func() {
		file_leaderelection_proto_rawDescData = protoimpl.X.CompressGZIP(file_leaderelection_proto_rawDescData)
	})
	return file_leaderelection_proto_rawDescData
}

var file_leaderelection_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_leaderelection_proto_goTypes = []interface{}{
	(*LeaderElectionMessage)(nil), // 0: gonest.LeaderElectionMessage
}
var file_leaderelection_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_leaderelection_proto_init() }
func file_leaderelection_proto_init() {
	if File_leaderelection_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_leaderelection_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaderElectionMessage); i {
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
			RawDescriptor: file_leaderelection_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_leaderelection_proto_goTypes,
		DependencyIndexes: file_leaderelection_proto_depIdxs,
		MessageInfos:      file_leaderelection_proto_msgTypes,
	}.Build()
	File_leaderelection_proto = out.File
	file_leaderelection_proto_rawDesc = nil
	file_leaderelection_proto_goTypes = nil
	file_leaderelection_proto_depIdxs = nil
}
