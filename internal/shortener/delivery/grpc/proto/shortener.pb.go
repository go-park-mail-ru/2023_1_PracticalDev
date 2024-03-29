// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shortener.proto

package shortener

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type StringMessage struct {
	URL                  string   `protobuf:"bytes,1,opt,name=URL,proto3" json:"URL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StringMessage) Reset()         { *m = StringMessage{} }
func (m *StringMessage) String() string { return proto.CompactTextString(m) }
func (*StringMessage) ProtoMessage()    {}
func (*StringMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{0}
}

func (m *StringMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StringMessage.Unmarshal(m, b)
}
func (m *StringMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StringMessage.Marshal(b, m, deterministic)
}
func (m *StringMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringMessage.Merge(m, src)
}
func (m *StringMessage) XXX_Size() int {
	return xxx_messageInfo_StringMessage.Size(m)
}
func (m *StringMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StringMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StringMessage proto.InternalMessageInfo

func (m *StringMessage) GetURL() string {
	if m != nil {
		return m.URL
	}
	return ""
}

func init() {
	proto.RegisterType((*StringMessage)(nil), "shortener.StringMessage")
}

func init() {
	proto.RegisterFile("shortener.proto", fileDescriptor_6a64040fb43d257f)
}

var fileDescriptor_6a64040fb43d257f = []byte{
	// 132 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xce, 0xc8, 0x2f,
	0x2a, 0x49, 0xcd, 0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x84, 0x0b, 0x28,
	0x29, 0x72, 0xf1, 0x06, 0x97, 0x14, 0x65, 0xe6, 0xa5, 0xfb, 0xa6, 0x16, 0x17, 0x27, 0xa6, 0xa7,
	0x0a, 0x09, 0x70, 0x31, 0x87, 0x06, 0xf9, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0x98,
	0x46, 0x1d, 0x8c, 0x5c, 0x9c, 0xc1, 0x30, 0x0d, 0x42, 0xd6, 0x5c, 0xcc, 0xee, 0xa9, 0x25, 0x42,
	0x12, 0x7a, 0x08, 0x43, 0x51, 0x0c, 0x90, 0xc2, 0x29, 0xa3, 0xc4, 0x20, 0x64, 0xc7, 0xc5, 0xe6,
	0x5c, 0x94, 0x9a, 0x58, 0x92, 0x4a, 0x9e, 0x7e, 0x27, 0xbe, 0x28, 0x1e, 0x3d, 0x7d, 0x6b, 0xb8,
	0x7c, 0x12, 0x1b, 0xd8, 0x3f, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xea, 0xdd, 0x23, 0x78,
	0xe2, 0x00, 0x00, 0x00,
}
