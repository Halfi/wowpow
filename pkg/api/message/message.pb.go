// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: message.proto

package message

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message_Header int32

const (
	Message_close     Message_Header = 0
	Message_challenge Message_Header = 1
	Message_resource  Message_Header = 2
)

// Enum value maps for Message_Header.
var (
	Message_Header_name = map[int32]string{
		0: "close",
		1: "challenge",
		2: "resource",
	}
	Message_Header_value = map[string]int32{
		"close":     0,
		"challenge": 1,
		"resource":  2,
	}
)

func (x Message_Header) Enum() *Message_Header {
	p := new(Message_Header)
	*p = x
	return p
}

func (x Message_Header) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Message_Header) Descriptor() protoreflect.EnumDescriptor {
	return file_message_proto_enumTypes[0].Descriptor()
}

func (Message_Header) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[0]
}

func (x Message_Header) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Message_Header.Descriptor instead.
func (Message_Header) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0, 0}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header Message_Header `protobuf:"varint,1,opt,name=header,proto3,enum=message.Message_Header" json:"header,omitempty"`
	// Types that are assignable to Response:
	//	*Message_Hashcach
	//	*Message_Payload
	Response isMessage_Response `protobuf_oneof:"response"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetHeader() Message_Header {
	if x != nil {
		return x.Header
	}
	return Message_close
}

func (m *Message) GetResponse() isMessage_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (x *Message) GetHashcach() *Hashcach {
	if x, ok := x.GetResponse().(*Message_Hashcach); ok {
		return x.Hashcach
	}
	return nil
}

func (x *Message) GetPayload() string {
	if x, ok := x.GetResponse().(*Message_Payload); ok {
		return x.Payload
	}
	return ""
}

type isMessage_Response interface {
	isMessage_Response()
}

type Message_Hashcach struct {
	Hashcach *Hashcach `protobuf:"bytes,2,opt,name=hashcach,proto3,oneof"`
}

type Message_Payload struct {
	Payload string `protobuf:"bytes,3,opt,name=payload,proto3,oneof"`
}

func (*Message_Hashcach) isMessage_Response() {}

func (*Message_Payload) isMessage_Response() {}

type Hashcach struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version  int32                  `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Bits     int32                  `protobuf:"varint,2,opt,name=bits,proto3" json:"bits,omitempty"`
	Date     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=date,proto3" json:"date,omitempty"`
	Resource string                 `protobuf:"bytes,4,opt,name=resource,proto3" json:"resource,omitempty"`
	Ext      string                 `protobuf:"bytes,5,opt,name=ext,proto3" json:"ext,omitempty"`
	Rand     string                 `protobuf:"bytes,6,opt,name=rand,proto3" json:"rand,omitempty"`
	Counter  string                 `protobuf:"bytes,7,opt,name=counter,proto3" json:"counter,omitempty"`
}

func (x *Hashcach) Reset() {
	*x = Hashcach{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Hashcach) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hashcach) ProtoMessage() {}

func (x *Hashcach) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hashcach.ProtoReflect.Descriptor instead.
func (*Hashcach) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *Hashcach) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Hashcach) GetBits() int32 {
	if x != nil {
		return x.Bits
	}
	return 0
}

func (x *Hashcach) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *Hashcach) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

func (x *Hashcach) GetExt() string {
	if x != nil {
		return x.Ext
	}
	return ""
}

func (x *Hashcach) GetRand() string {
	if x != nil {
		return x.Rand
	}
	return ""
}

func (x *Hashcach) GetCounter() string {
	if x != nil {
		return x.Counter
	}
	return ""
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc5, 0x01, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x68, 0x63, 0x61,
	0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x63, 0x61, 0x63, 0x68, 0x48, 0x00, 0x52, 0x08, 0x68,
	0x61, 0x73, 0x68, 0x63, 0x61, 0x63, 0x68, 0x12, 0x1a, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0x30, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x09, 0x0a,
	0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x6c,
	0x6c, 0x65, 0x6e, 0x67, 0x65, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x10, 0x02, 0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0xc4, 0x01, 0x0a, 0x08, 0x48, 0x61, 0x73, 0x68, 0x63, 0x61, 0x63, 0x68, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x69, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x62, 0x69, 0x74, 0x73, 0x12, 0x2e, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x78, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x78, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x61,
	0x6e, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x42, 0x18, 0x5a, 0x16, 0x77, 0x6f, 0x77, 0x70,
	0x6f, 0x77, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_proto_goTypes = []interface{}{
	(Message_Header)(0),           // 0: message.Message.Header
	(*Message)(nil),               // 1: message.Message
	(*Hashcach)(nil),              // 2: message.Hashcach
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_message_proto_depIdxs = []int32{
	0, // 0: message.Message.header:type_name -> message.Message.Header
	2, // 1: message.Message.hashcach:type_name -> message.Hashcach
	3, // 2: message.Hashcach.date:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Hashcach); i {
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
	file_message_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_Hashcach)(nil),
		(*Message_Payload)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		EnumInfos:         file_message_proto_enumTypes,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}
