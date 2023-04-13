// protoc --go_out=. --go-grpc_out=. ./*.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: inner.proto

package transport

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

type SingleRecv struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgId     uint32 `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`             // 消息id
	AreaId    int32  `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`          // 区域id
	ChannelId int32  `protobuf:"varint,3,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"` // 渠道id
	UserId    uint64 `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`          // 用户id
	Content   []byte `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`                       // 内容
}

func (x *SingleRecv) Reset() {
	*x = SingleRecv{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingleRecv) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingleRecv) ProtoMessage() {}

func (x *SingleRecv) ProtoReflect() protoreflect.Message {
	mi := &file_inner_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingleRecv.ProtoReflect.Descriptor instead.
func (*SingleRecv) Descriptor() ([]byte, []int) {
	return file_inner_proto_rawDescGZIP(), []int{0}
}

func (x *SingleRecv) GetMsgId() uint32 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

func (x *SingleRecv) GetAreaId() int32 {
	if x != nil {
		return x.AreaId
	}
	return 0
}

func (x *SingleRecv) GetChannelId() int32 {
	if x != nil {
		return x.ChannelId
	}
	return 0
}

func (x *SingleRecv) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *SingleRecv) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type SingleBack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // 用户id
	Content []byte `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`              // 内容
}

func (x *SingleBack) Reset() {
	*x = SingleBack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SingleBack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SingleBack) ProtoMessage() {}

func (x *SingleBack) ProtoReflect() protoreflect.Message {
	mi := &file_inner_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SingleBack.ProtoReflect.Descriptor instead.
func (*SingleBack) Descriptor() ([]byte, []int) {
	return file_inner_proto_rawDescGZIP(), []int{1}
}

func (x *SingleBack) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *SingleBack) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type MessageRouteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgs []*SingleRecv `protobuf:"bytes,1,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (x *MessageRouteReq) Reset() {
	*x = MessageRouteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageRouteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageRouteReq) ProtoMessage() {}

func (x *MessageRouteReq) ProtoReflect() protoreflect.Message {
	mi := &file_inner_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageRouteReq.ProtoReflect.Descriptor instead.
func (*MessageRouteReq) Descriptor() ([]byte, []int) {
	return file_inner_proto_rawDescGZIP(), []int{2}
}

func (x *MessageRouteReq) GetMsgs() []*SingleRecv {
	if x != nil {
		return x.Msgs
	}
	return nil
}

type MessageRouteRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgs []*SingleBack `protobuf:"bytes,1,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (x *MessageRouteRsp) Reset() {
	*x = MessageRouteRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageRouteRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageRouteRsp) ProtoMessage() {}

func (x *MessageRouteRsp) ProtoReflect() protoreflect.Message {
	mi := &file_inner_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageRouteRsp.ProtoReflect.Descriptor instead.
func (*MessageRouteRsp) Descriptor() ([]byte, []int) {
	return file_inner_proto_rawDescGZIP(), []int{3}
}

func (x *MessageRouteRsp) GetMsgs() []*SingleBack {
	if x != nil {
		return x.Msgs
	}
	return nil
}

var File_inner_proto protoreflect.FileDescriptor

var file_inner_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x8e, 0x01, 0x0a, 0x0a, 0x53, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x52, 0x65, 0x63, 0x76, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x61, 0x72, 0x65, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x61, 0x72, 0x65, 0x61, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3f, 0x0a, 0x0a, 0x53, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x42, 0x61, 0x63, 0x6b, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3c, 0x0a, 0x0f, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12, 0x29, 0x0a,
	0x04, 0x6d, 0x73, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x52, 0x65,
	0x63, 0x76, 0x52, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x22, 0x3c, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x73, 0x70, 0x12, 0x29, 0x0a, 0x04, 0x6d,
	0x73, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x42, 0x61, 0x63, 0x6b,
	0x52, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x32, 0x51, 0x0a, 0x05, 0x49, 0x6e, 0x6e, 0x65, 0x72, 0x12,
	0x48, 0x0a, 0x0c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12,
	0x1a, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x1a, 0x2e, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x52, 0x73, 0x70, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x3b, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inner_proto_rawDescOnce sync.Once
	file_inner_proto_rawDescData = file_inner_proto_rawDesc
)

func file_inner_proto_rawDescGZIP() []byte {
	file_inner_proto_rawDescOnce.Do(func() {
		file_inner_proto_rawDescData = protoimpl.X.CompressGZIP(file_inner_proto_rawDescData)
	})
	return file_inner_proto_rawDescData
}

var file_inner_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_inner_proto_goTypes = []interface{}{
	(*SingleRecv)(nil),      // 0: transport.SingleRecv
	(*SingleBack)(nil),      // 1: transport.SingleBack
	(*MessageRouteReq)(nil), // 2: transport.MessageRouteReq
	(*MessageRouteRsp)(nil), // 3: transport.MessageRouteRsp
}
var file_inner_proto_depIdxs = []int32{
	0, // 0: transport.MessageRouteReq.msgs:type_name -> transport.SingleRecv
	1, // 1: transport.MessageRouteRsp.msgs:type_name -> transport.SingleBack
	2, // 2: transport.Inner.MessageRoute:input_type -> transport.MessageRouteReq
	3, // 3: transport.Inner.MessageRoute:output_type -> transport.MessageRouteRsp
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_inner_proto_init() }
func file_inner_proto_init() {
	if File_inner_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_inner_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SingleRecv); i {
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
		file_inner_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SingleBack); i {
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
		file_inner_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageRouteReq); i {
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
		file_inner_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageRouteRsp); i {
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
			RawDescriptor: file_inner_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_inner_proto_goTypes,
		DependencyIndexes: file_inner_proto_depIdxs,
		MessageInfos:      file_inner_proto_msgTypes,
	}.Build()
	File_inner_proto = out.File
	file_inner_proto_rawDesc = nil
	file_inner_proto_goTypes = nil
	file_inner_proto_depIdxs = nil
}
