// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.proto

/*
Package Im is a generated protocol buffer package.

It is generated from these files:
	msg.proto

It has these top-level messages:
	Helloworld
*/
package Im

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

type Helloworld struct {
	Id  int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Str string `protobuf:"bytes,2,opt,name=str" json:"str,omitempty"`
	Opt int32  `protobuf:"varint,3,opt,name=opt" json:"opt,omitempty"`
}

func (m *Helloworld) Reset()                    { *m = Helloworld{} }
func (m *Helloworld) String() string            { return proto.CompactTextString(m) }
func (*Helloworld) ProtoMessage()               {}
func (*Helloworld) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Helloworld) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Helloworld) GetStr() string {
	if m != nil {
		return m.Str
	}
	return ""
}

func (m *Helloworld) GetOpt() int32 {
	if m != nil {
		return m.Opt
	}
	return 0
}

func init() {
	proto.RegisterType((*Helloworld)(nil), "Im.helloworld")
}

func init() { proto.RegisterFile("msg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 101 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcc, 0x2d, 0x4e, 0xd7,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xf2, 0xcc, 0x55, 0x72, 0xe0, 0xe2, 0xca, 0x48, 0xcd,
	0xc9, 0xc9, 0x2f, 0xcf, 0x2f, 0xca, 0x49, 0x11, 0xe2, 0xe3, 0x62, 0xca, 0x4c, 0x91, 0x60, 0x54,
	0x60, 0xd4, 0x60, 0x0d, 0x62, 0xca, 0x4c, 0x11, 0x12, 0xe0, 0x62, 0x2e, 0x2e, 0x29, 0x92, 0x60,
	0x52, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0x31, 0x41, 0x22, 0xf9, 0x05, 0x25, 0x12, 0xcc, 0x60, 0x25,
	0x20, 0x66, 0x12, 0x1b, 0xd8, 0x30, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdf, 0xb0, 0xe7,
	0x7a, 0x59, 0x00, 0x00, 0x00,
}
