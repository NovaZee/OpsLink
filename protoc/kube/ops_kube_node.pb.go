// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.6.1
// source: ops_kube_node.proto

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

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Ip         string        `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	HostName   string        `protobuf:"bytes,3,opt,name=host_name,json=hostName,proto3" json:"host_name,omitempty"`
	Labels     []string      `protobuf:"bytes,6,rep,name=labels,proto3" json:"labels,omitempty"`
	Taints     []string      `protobuf:"bytes,7,rep,name=taints,proto3" json:"taints,omitempty"`
	Capacity   *NodeCapacity `protobuf:"bytes,8,opt,name=capacity,proto3" json:"capacity,omitempty"`
	Usage      *NodeUsage    `protobuf:"bytes,9,opt,name=usage,proto3" json:"usage,omitempty"`
	CreateTime string        `protobuf:"bytes,10,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ops_kube_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_ops_kube_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_ops_kube_node_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Node) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Node) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *Node) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Node) GetTaints() []string {
	if x != nil {
		return x.Taints
	}
	return nil
}

func (x *Node) GetCapacity() *NodeCapacity {
	if x != nil {
		return x.Capacity
	}
	return nil
}

func (x *Node) GetUsage() *NodeUsage {
	if x != nil {
		return x.Usage
	}
	return nil
}

func (x *Node) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

type NodeUsage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pods   int32   `protobuf:"varint,1,opt,name=pods,proto3" json:"pods,omitempty"`
	Cpu    float64 `protobuf:"fixed64,2,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory float64 `protobuf:"fixed64,3,opt,name=memory,proto3" json:"memory,omitempty"`
}

func (x *NodeUsage) Reset() {
	*x = NodeUsage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ops_kube_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeUsage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeUsage) ProtoMessage() {}

func (x *NodeUsage) ProtoReflect() protoreflect.Message {
	mi := &file_ops_kube_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeUsage.ProtoReflect.Descriptor instead.
func (*NodeUsage) Descriptor() ([]byte, []int) {
	return file_ops_kube_node_proto_rawDescGZIP(), []int{1}
}

func (x *NodeUsage) GetPods() int32 {
	if x != nil {
		return x.Pods
	}
	return 0
}

func (x *NodeUsage) GetCpu() float64 {
	if x != nil {
		return x.Cpu
	}
	return 0
}

func (x *NodeUsage) GetMemory() float64 {
	if x != nil {
		return x.Memory
	}
	return 0
}

type NodeCapacity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cpu    int64 `protobuf:"varint,1,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory int64 `protobuf:"varint,2,opt,name=memory,proto3" json:"memory,omitempty"`
	Pods   int64 `protobuf:"varint,3,opt,name=pods,proto3" json:"pods,omitempty"`
}

func (x *NodeCapacity) Reset() {
	*x = NodeCapacity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ops_kube_node_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeCapacity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeCapacity) ProtoMessage() {}

func (x *NodeCapacity) ProtoReflect() protoreflect.Message {
	mi := &file_ops_kube_node_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeCapacity.ProtoReflect.Descriptor instead.
func (*NodeCapacity) Descriptor() ([]byte, []int) {
	return file_ops_kube_node_proto_rawDescGZIP(), []int{2}
}

func (x *NodeCapacity) GetCpu() int64 {
	if x != nil {
		return x.Cpu
	}
	return 0
}

func (x *NodeCapacity) GetMemory() int64 {
	if x != nil {
		return x.Memory
	}
	return 0
}

func (x *NodeCapacity) GetPods() int64 {
	if x != nil {
		return x.Pods
	}
	return 0
}

type Taint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value  string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Effect string `protobuf:"bytes,3,opt,name=effect,proto3" json:"effect,omitempty"`
}

func (x *Taint) Reset() {
	*x = Taint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ops_kube_node_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Taint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Taint) ProtoMessage() {}

func (x *Taint) ProtoReflect() protoreflect.Message {
	mi := &file_ops_kube_node_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Taint.ProtoReflect.Descriptor instead.
func (*Taint) Descriptor() ([]byte, []int) {
	return file_ops_kube_node_proto_rawDescGZIP(), []int{3}
}

func (x *Taint) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Taint) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Taint) GetEffect() string {
	if x != nil {
		return x.Effect
	}
	return ""
}

type FrontNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Labels []*Map   `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty"`
	Taints []*Taint `protobuf:"bytes,3,rep,name=taints,proto3" json:"taints,omitempty"`
}

func (x *FrontNode) Reset() {
	*x = FrontNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ops_kube_node_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FrontNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FrontNode) ProtoMessage() {}

func (x *FrontNode) ProtoReflect() protoreflect.Message {
	mi := &file_ops_kube_node_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FrontNode.ProtoReflect.Descriptor instead.
func (*FrontNode) Descriptor() ([]byte, []int) {
	return file_ops_kube_node_proto_rawDescGZIP(), []int{4}
}

func (x *FrontNode) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FrontNode) GetLabels() []*Map {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *FrontNode) GetTaints() []*Taint {
	if x != nil {
		return x.Taints
	}
	return nil
}

var File_ops_kube_node_proto protoreflect.FileDescriptor

var file_ops_kube_node_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6f, 0x70, 0x73, 0x5f, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6b, 0x75, 0x62, 0x65, 0x1a, 0x14, 0x6f, 0x70, 0x73,
	0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xef, 0x01, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x1b,
	0x0a, 0x09, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x2e, 0x0a, 0x08, 0x63,
	0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x6b, 0x75, 0x62, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74,
	0x79, 0x52, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12, 0x25, 0x0a, 0x05, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6b, 0x75, 0x62,
	0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x05, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x22, 0x49, 0x0a, 0x09, 0x4e, 0x6f, 0x64, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x6f, 0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x03, 0x63, 0x70, 0x75, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x22, 0x4c,
	0x0a, 0x0c, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x63, 0x70, 0x75, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x63, 0x70, 0x75,
	0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x64, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x6f, 0x64, 0x73, 0x22, 0x47, 0x0a, 0x05,
	0x54, 0x61, 0x69, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65,
	0x66, 0x66, 0x65, 0x63, 0x74, 0x22, 0x67, 0x0a, 0x09, 0x46, 0x72, 0x6f, 0x6e, 0x74, 0x4e, 0x6f,
	0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x2e, 0x4d, 0x61,
	0x70, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x23, 0x0a, 0x06, 0x74, 0x61, 0x69,
	0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x6b, 0x75, 0x62, 0x65,
	0x2e, 0x54, 0x61, 0x69, 0x6e, 0x74, 0x52, 0x06, 0x74, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x42, 0x07,
	0x5a, 0x05, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ops_kube_node_proto_rawDescOnce sync.Once
	file_ops_kube_node_proto_rawDescData = file_ops_kube_node_proto_rawDesc
)

func file_ops_kube_node_proto_rawDescGZIP() []byte {
	file_ops_kube_node_proto_rawDescOnce.Do(func() {
		file_ops_kube_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_ops_kube_node_proto_rawDescData)
	})
	return file_ops_kube_node_proto_rawDescData
}

var file_ops_kube_node_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_ops_kube_node_proto_goTypes = []interface{}{
	(*Node)(nil),         // 0: kube.node
	(*NodeUsage)(nil),    // 1: kube.NodeUsage
	(*NodeCapacity)(nil), // 2: kube.NodeCapacity
	(*Taint)(nil),        // 3: kube.Taint
	(*FrontNode)(nil),    // 4: kube.FrontNode
	(*Map)(nil),          // 5: kube.Map
}
var file_ops_kube_node_proto_depIdxs = []int32{
	2, // 0: kube.node.capacity:type_name -> kube.NodeCapacity
	1, // 1: kube.node.usage:type_name -> kube.NodeUsage
	5, // 2: kube.FrontNode.labels:type_name -> kube.Map
	3, // 3: kube.FrontNode.taints:type_name -> kube.Taint
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_ops_kube_node_proto_init() }
func file_ops_kube_node_proto_init() {
	if File_ops_kube_node_proto != nil {
		return
	}
	file_ops_collection_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ops_kube_node_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_ops_kube_node_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeUsage); i {
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
		file_ops_kube_node_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeCapacity); i {
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
		file_ops_kube_node_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Taint); i {
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
		file_ops_kube_node_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FrontNode); i {
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
			RawDescriptor: file_ops_kube_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ops_kube_node_proto_goTypes,
		DependencyIndexes: file_ops_kube_node_proto_depIdxs,
		MessageInfos:      file_ops_kube_node_proto_msgTypes,
	}.Build()
	File_ops_kube_node_proto = out.File
	file_ops_kube_node_proto_rawDesc = nil
	file_ops_kube_node_proto_goTypes = nil
	file_ops_kube_node_proto_depIdxs = nil
}
