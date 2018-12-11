// Code generated by protoc-gen-go. DO NOT EDIT.
// source: something.proto

package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Item struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp            int64    `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Item) Reset()         { *m = Item{} }
func (m *Item) String() string { return proto.CompactTextString(m) }
func (*Item) ProtoMessage()    {}
func (*Item) Descriptor() ([]byte, []int) {
	return fileDescriptor_80f5c795273d78da, []int{0}
}
func (m *Item) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Item.Unmarshal(m, b)
}
func (m *Item) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Item.Marshal(b, m, deterministic)
}
func (dst *Item) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Item.Merge(dst, src)
}
func (m *Item) XXX_Size() int {
	return xxx_messageInfo_Item.Size(m)
}
func (m *Item) XXX_DiscardUnknown() {
	xxx_messageInfo_Item.DiscardUnknown(m)
}

var xxx_messageInfo_Item proto.InternalMessageInfo

func (m *Item) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Item) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Item) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type Changelog struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Timestamp            int64    `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	ObjectID             string   `protobuf:"bytes,3,opt,name=objectID,proto3" json:"objectID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Changelog) Reset()         { *m = Changelog{} }
func (m *Changelog) String() string { return proto.CompactTextString(m) }
func (*Changelog) ProtoMessage()    {}
func (*Changelog) Descriptor() ([]byte, []int) {
	return fileDescriptor_80f5c795273d78da, []int{1}
}
func (m *Changelog) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Changelog.Unmarshal(m, b)
}
func (m *Changelog) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Changelog.Marshal(b, m, deterministic)
}
func (dst *Changelog) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Changelog.Merge(dst, src)
}
func (m *Changelog) XXX_Size() int {
	return xxx_messageInfo_Changelog.Size(m)
}
func (m *Changelog) XXX_DiscardUnknown() {
	xxx_messageInfo_Changelog.DiscardUnknown(m)
}

var xxx_messageInfo_Changelog proto.InternalMessageInfo

func (m *Changelog) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Changelog) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Changelog) GetObjectID() string {
	if m != nil {
		return m.ObjectID
	}
	return ""
}

func init() {
	proto.RegisterType((*Item)(nil), "Item")
	proto.RegisterType((*Changelog)(nil), "Changelog")
}

func init() { proto.RegisterFile("something.proto", fileDescriptor_80f5c795273d78da) }

var fileDescriptor_80f5c795273d78da = []byte{
	// 157 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xce, 0xcf, 0x4d,
	0x2d, 0xc9, 0xc8, 0xcc, 0x4b, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xf2, 0xe2, 0x62, 0xf1,
	0x2c, 0x49, 0xcd, 0x15, 0xe2, 0xe3, 0x62, 0xca, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c,
	0x62, 0xca, 0x4c, 0x11, 0x12, 0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60, 0x52, 0x60,
	0xd4, 0xe0, 0x09, 0x82, 0x70, 0x84, 0x64, 0xb8, 0x38, 0x4b, 0x32, 0x73, 0x53, 0x8b, 0x4b, 0x12,
	0x73, 0x0b, 0x24, 0x98, 0x15, 0x18, 0x35, 0x98, 0x83, 0x10, 0x02, 0x4a, 0xa1, 0x5c, 0x9c, 0xce,
	0x19, 0x89, 0x79, 0xe9, 0xa9, 0x39, 0xf9, 0xe9, 0x18, 0x06, 0xa2, 0x68, 0x65, 0x42, 0xd3, 0x2a,
	0x24, 0xc5, 0xc5, 0x91, 0x9f, 0x94, 0x95, 0x9a, 0x5c, 0xe2, 0xe9, 0x02, 0x36, 0x97, 0x33, 0x08,
	0xce, 0x77, 0xe2, 0x8a, 0xe2, 0x00, 0xbb, 0x35, 0xa9, 0x34, 0x2d, 0x89, 0x0d, 0xcc, 0x32, 0x06,
	0x04, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x44, 0xc7, 0xab, 0xc8, 0x00, 0x00, 0x00,
}
