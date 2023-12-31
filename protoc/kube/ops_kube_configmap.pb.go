// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.6.1
// source: ops_kube_configmap.proto

package kube

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

type ConfigMap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace  string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	CreateTime string `protobuf:"bytes,3,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	Data       []*Map `protobuf:"bytes,4,rep,name=Data,proto3" json:"Data,omitempty"`
	Update     bool   `protobuf:"varint,5,opt,name=update,proto3" json:"update,omitempty"`
}

func (x *ConfigMap) Reset() {
	*x = ConfigMap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ops_kube_configmap_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigMap) ProtoMessage() {}

func (x *ConfigMap) ProtoReflect() protoreflect.Message {
	mi := &file_ops_kube_configmap_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigMap.ProtoReflect.Descriptor instead.
func (*ConfigMap) Descriptor() ([]byte, []int) {
	return file_ops_kube_configmap_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigMap) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ConfigMap) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ConfigMap) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

func (x *ConfigMap) GetData() []*Map {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ConfigMap) GetUpdate() bool {
	if x != nil {
		return x.Update
	}
	return false
}

var File_ops_kube_configmap_proto protoreflect.FileDescriptor

var file_ops_kube_configmap_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6f, 0x70, 0x73, 0x5f, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x6d, 0x61, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6b, 0x75, 0x62, 0x65,
	0x1a, 0x14, 0x6f, 0x70, 0x73, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x4d, 0x61, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x2e, 0x4d, 0x61, 0x70, 0x52,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x07, 0x5a,
	0x05, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ops_kube_configmap_proto_rawDescOnce sync.Once
	file_ops_kube_configmap_proto_rawDescData = file_ops_kube_configmap_proto_rawDesc
)

func file_ops_kube_configmap_proto_rawDescGZIP() []byte {
	file_ops_kube_configmap_proto_rawDescOnce.Do(func() {
		file_ops_kube_configmap_proto_rawDescData = protoimpl.X.CompressGZIP(file_ops_kube_configmap_proto_rawDescData)
	})
	return file_ops_kube_configmap_proto_rawDescData
}

var file_ops_kube_configmap_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ops_kube_configmap_proto_goTypes = []interface{}{
	(*ConfigMap)(nil), // 0: kube.ConfigMap
	(*Map)(nil),       // 1: kube.Map
}
var file_ops_kube_configmap_proto_depIdxs = []int32{
	1, // 0: kube.ConfigMap.Data:type_name -> kube.Map
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ops_kube_configmap_proto_init() }
func file_ops_kube_configmap_proto_init() {
	if File_ops_kube_configmap_proto != nil {
		return
	}
	file_ops_collection_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ops_kube_configmap_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigMap); i {
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
			RawDescriptor: file_ops_kube_configmap_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ops_kube_configmap_proto_goTypes,
		DependencyIndexes: file_ops_kube_configmap_proto_depIdxs,
		MessageInfos:      file_ops_kube_configmap_proto_msgTypes,
	}.Build()
	File_ops_kube_configmap_proto = out.File
	file_ops_kube_configmap_proto_rawDesc = nil
	file_ops_kube_configmap_proto_goTypes = nil
	file_ops_kube_configmap_proto_depIdxs = nil
}
