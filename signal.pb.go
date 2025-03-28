// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: signal.proto

package ophelia_ci

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CommitRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CommitHash    string                 `protobuf:"bytes,1,opt,name=commit_hash,json=commitHash,proto3" json:"commit_hash,omitempty"`
	Branch        string                 `protobuf:"bytes,2,opt,name=branch,proto3" json:"branch,omitempty"`
	Repository    string                 `protobuf:"bytes,3,opt,name=repository,proto3" json:"repository,omitempty"`
	Tag           string                 `protobuf:"bytes,4,opt,name=tag,proto3" json:"tag,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommitRequest) Reset() {
	*x = CommitRequest{}
	mi := &file_signal_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitRequest) ProtoMessage() {}

func (x *CommitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_signal_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitRequest.ProtoReflect.Descriptor instead.
func (*CommitRequest) Descriptor() ([]byte, []int) {
	return file_signal_proto_rawDescGZIP(), []int{0}
}

func (x *CommitRequest) GetCommitHash() string {
	if x != nil {
		return x.CommitHash
	}
	return ""
}

func (x *CommitRequest) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *CommitRequest) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

func (x *CommitRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

var File_signal_proto protoreflect.FileDescriptor

var file_signal_proto_rawDesc = string([]byte{
	0x0a, 0x0c, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1e,
	0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67,
	0x32, 0x3f, 0x0a, 0x07, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x73, 0x12, 0x34, 0x0a, 0x0c, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x12, 0x15, 0x2e, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x45, 0x64, 0x6d, 0x69, 0x6c, 0x73, 0x6f, 0x6e, 0x52, 0x6f, 0x64, 0x72, 0x69, 0x67, 0x75, 0x65,
	0x73, 0x2f, 0x6f, 0x70, 0x68, 0x65, 0x6c, 0x69, 0x61, 0x2d, 0x63, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_signal_proto_rawDescOnce sync.Once
	file_signal_proto_rawDescData []byte
)

func file_signal_proto_rawDescGZIP() []byte {
	file_signal_proto_rawDescOnce.Do(func() {
		file_signal_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_signal_proto_rawDesc), len(file_signal_proto_rawDesc)))
	})
	return file_signal_proto_rawDescData
}

var file_signal_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_signal_proto_goTypes = []any{
	(*CommitRequest)(nil), // 0: signal.CommitRequest
	(*Empty)(nil),         // 1: common.Empty
}
var file_signal_proto_depIdxs = []int32{
	0, // 0: signal.Signals.CommitSignal:input_type -> signal.CommitRequest
	1, // 1: signal.Signals.CommitSignal:output_type -> common.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_signal_proto_init() }
func file_signal_proto_init() {
	if File_signal_proto != nil {
		return
	}
	file_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_signal_proto_rawDesc), len(file_signal_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_signal_proto_goTypes,
		DependencyIndexes: file_signal_proto_depIdxs,
		MessageInfos:      file_signal_proto_msgTypes,
	}.Build()
	File_signal_proto = out.File
	file_signal_proto_goTypes = nil
	file_signal_proto_depIdxs = nil
}
