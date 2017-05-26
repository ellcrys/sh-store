// Code generated by protoc-gen-gogo.
// source: server.proto
// DO NOT EDIT!

/*
Package proto_rpc is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
	DBSession
	CreateObjectsMsg
	CreateIdentityMsg
	ObjectResponse
	MultiObjectResponse
	ObjectResponseData
	Object
	Refs
*/
package proto_rpc

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// DBSession represents a database session request body
type DBSession struct {
	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty"`
}

func (m *DBSession) Reset()                    { *m = DBSession{} }
func (m *DBSession) String() string            { return proto.CompactTextString(m) }
func (*DBSession) ProtoMessage()               {}
func (*DBSession) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{0} }

func (m *DBSession) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

// CreateObjects represents a collection of objects to be created
type CreateObjectsMsg struct {
	Objects []*Object `protobuf:"bytes,1,rep,name=objects" json:"objects,omitempty" structs:"objects,omitempty" mapstructure:"objects,omitempty"`
}

func (m *CreateObjectsMsg) Reset()                    { *m = CreateObjectsMsg{} }
func (m *CreateObjectsMsg) String() string            { return proto.CompactTextString(m) }
func (*CreateObjectsMsg) ProtoMessage()               {}
func (*CreateObjectsMsg) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{1} }

func (m *CreateObjectsMsg) GetObjects() []*Object {
	if m != nil {
		return m.Objects
	}
	return nil
}

// CreateIdentityMsg represents an identity of a person or organization
type CreateIdentityMsg struct {
	Email     string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" structs:"email,omitempty" mapstructure:"email,omitempty"`
	Password  string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty" structs:"password,omitempty" mapstructure:"password,omitempty"`
	Developer bool   `protobuf:"varint,3,opt,name=developer,proto3" json:"developer,omitempty" structs:"developer,omitempty" mapstructure:"developer,omitempty"`
}

func (m *CreateIdentityMsg) Reset()                    { *m = CreateIdentityMsg{} }
func (m *CreateIdentityMsg) String() string            { return proto.CompactTextString(m) }
func (*CreateIdentityMsg) ProtoMessage()               {}
func (*CreateIdentityMsg) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{2} }

func (m *CreateIdentityMsg) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateIdentityMsg) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *CreateIdentityMsg) GetDeveloper() bool {
	if m != nil {
		return m.Developer
	}
	return false
}

// ObjectResponse defines a response about a single object
type ObjectResponse struct {
	Data  *ObjectResponseData `protobuf:"bytes,1,opt,name=data" json:"data,omitempty" structs:"data,omitempty" mapstructure:"data,omitempty"`
	Links map[string]string   `protobuf:"bytes,2,rep,name=links" json:"links,omitempty" structs:"links,omitempty" mapstructure:"links,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *ObjectResponse) Reset()                    { *m = ObjectResponse{} }
func (m *ObjectResponse) String() string            { return proto.CompactTextString(m) }
func (*ObjectResponse) ProtoMessage()               {}
func (*ObjectResponse) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{3} }

func (m *ObjectResponse) GetData() *ObjectResponseData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ObjectResponse) GetLinks() map[string]string {
	if m != nil {
		return m.Links
	}
	return nil
}

// MultiObjectResponse describes a response of multiple objects
type MultiObjectResponse struct {
	Data  []*ObjectResponseData `protobuf:"bytes,1,rep,name=data" json:"data,omitempty" structs:"data,omitempty" mapstructure:"data,omitempty"`
	Links map[string]string     `protobuf:"bytes,2,rep,name=links" json:"links,omitempty" structs:"links,omitempty" mapstructure:"links,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *MultiObjectResponse) Reset()                    { *m = MultiObjectResponse{} }
func (m *MultiObjectResponse) String() string            { return proto.CompactTextString(m) }
func (*MultiObjectResponse) ProtoMessage()               {}
func (*MultiObjectResponse) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{4} }

func (m *MultiObjectResponse) GetData() []*ObjectResponseData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *MultiObjectResponse) GetLinks() map[string]string {
	if m != nil {
		return m.Links
	}
	return nil
}

// ObjectResponse describes an object and related information
type ObjectResponseData struct {
	Type       string  `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty" structs:"type,omitempty" mapstructure:"type,omitempty"`
	ID         string  `protobuf:"bytes,2,opt,name=ID,proto3" json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty"`
	Attributes *Object `protobuf:"bytes,3,opt,name=attributes" json:"attributes,omitempty" structs:"attributes,omitempty" mapstructure:"attributes,omitempty"`
}

func (m *ObjectResponseData) Reset()                    { *m = ObjectResponseData{} }
func (m *ObjectResponseData) String() string            { return proto.CompactTextString(m) }
func (*ObjectResponseData) ProtoMessage()               {}
func (*ObjectResponseData) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{5} }

func (m *ObjectResponseData) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ObjectResponseData) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *ObjectResponseData) GetAttributes() *Object {
	if m != nil {
		return m.Attributes
	}
	return nil
}

// Object describes an object
type Object struct {
	ID          string `protobuf:"bytes,1,opt,name=ID,proto3" json:"id,omitempty" structs:"id,omitempty" mapstructure:"id,omitempty"`
	OwnerID     string `protobuf:"bytes,2,opt,name=ownerID,proto3" json:"owner_id,omitempty" structs:"ownerId,omitempty" mapstructure:"ownerId,omitempty"`
	CreatorID   string `protobuf:"bytes,3,opt,name=creatorID,proto3" json:"creator_id,omitempty" structs:"creatorId,omitempty" mapstructure:"creatorId,omitempty"`
	PartitionID string `protobuf:"bytes,4,opt,name=partitionID,proto3" json:"partition_id,omitempty" structs:"partitionId,omitempty" mapstructure:"partitionId,omitempty"`
	Key         string `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty" structs:"key,omitempty" mapstructure:"key,omitempty"`
	Value       string `protobuf:"bytes,6,opt,name=value,proto3" json:"value,omitempty" structs:"value,omitempty" mapstructure:"value,omitempty"`
	Protected   bool   `protobuf:"varint,7,opt,name=protected,proto3" json:"protected,omitempty" structs:"protected,omitempty" mapstructure:"protected,omitempty"`
	RefOnly     bool   `protobuf:"varint,8,opt,name=refOnly,proto3" json:"ref_only,omitempty" structs:"refOnly,omitempty" mapstructure:"refOnly,omitempty"`
	Timestamp   int64  `protobuf:"varint,9,opt,name=timestamp,proto3" json:"timestamp,omitempty" structs:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`
	PrevHash    string `protobuf:"bytes,10,opt,name=prevHash,proto3" json:"prev_hash,omitempty" structs:"prevHash,omitempty" mapstructure:"prevHash,omitempty"`
	Hash        string `protobuf:"bytes,11,opt,name=hash,proto3" json:"hash,omitempty" structs:"hash,omitempty" mapstructure:"hash,omitempty"`
	Refs        *Refs  `protobuf:"bytes,12,opt,name=Refs" json:"refs,omitempty" structs:"refs,omitempty" mapstructure:"refs,omitempty"`
}

func (m *Object) Reset()                    { *m = Object{} }
func (m *Object) String() string            { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()               {}
func (*Object) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{6} }

func (m *Object) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Object) GetOwnerID() string {
	if m != nil {
		return m.OwnerID
	}
	return ""
}

func (m *Object) GetCreatorID() string {
	if m != nil {
		return m.CreatorID
	}
	return ""
}

func (m *Object) GetPartitionID() string {
	if m != nil {
		return m.PartitionID
	}
	return ""
}

func (m *Object) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Object) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Object) GetProtected() bool {
	if m != nil {
		return m.Protected
	}
	return false
}

func (m *Object) GetRefOnly() bool {
	if m != nil {
		return m.RefOnly
	}
	return false
}

func (m *Object) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Object) GetPrevHash() string {
	if m != nil {
		return m.PrevHash
	}
	return ""
}

func (m *Object) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *Object) GetRefs() *Refs {
	if m != nil {
		return m.Refs
	}
	return nil
}

type Refs struct {
	Ref1  string `protobuf:"bytes,12,opt,name=ref1,proto3" json:"ref1,omitempty" structs:"ref1,omitempty" mapstructure:"ref1,omitempty"`
	Ref2  string `protobuf:"bytes,13,opt,name=ref2,proto3" json:"ref2,omitempty" structs:"ref2,omitempty" mapstructure:"ref2,omitempty"`
	Ref3  string `protobuf:"bytes,14,opt,name=ref3,proto3" json:"ref3,omitempty" structs:"ref3,omitempty" mapstructure:"ref3,omitempty"`
	Ref4  string `protobuf:"bytes,15,opt,name=ref4,proto3" json:"ref4,omitempty" structs:"ref4,omitempty" mapstructure:"ref4,omitempty"`
	Ref5  string `protobuf:"bytes,16,opt,name=ref5,proto3" json:"ref5,omitempty" structs:"ref5,omitempty" mapstructure:"ref5,omitempty"`
	Ref6  string `protobuf:"bytes,17,opt,name=ref6,proto3" json:"ref6,omitempty" structs:"ref6,omitempty" mapstructure:"ref6,omitempty"`
	Ref7  string `protobuf:"bytes,18,opt,name=ref7,proto3" json:"ref7,omitempty" structs:"ref7,omitempty" mapstructure:"ref7,omitempty"`
	Ref8  string `protobuf:"bytes,19,opt,name=ref8,proto3" json:"ref8,omitempty" structs:"ref8,omitempty" mapstructure:"ref8,omitempty"`
	Ref9  string `protobuf:"bytes,20,opt,name=ref9,proto3" json:"ref9,omitempty" structs:"ref9,omitempty" mapstructure:"ref9,omitempty"`
	Ref10 string `protobuf:"bytes,21,opt,name=ref10,proto3" json:"ref10,omitempty" structs:"ref10,omitempty" mapstructure:"ref10,omitempty"`
}

func (m *Refs) Reset()                    { *m = Refs{} }
func (m *Refs) String() string            { return proto.CompactTextString(m) }
func (*Refs) ProtoMessage()               {}
func (*Refs) Descriptor() ([]byte, []int) { return fileDescriptorServer, []int{7} }

func (m *Refs) GetRef1() string {
	if m != nil {
		return m.Ref1
	}
	return ""
}

func (m *Refs) GetRef2() string {
	if m != nil {
		return m.Ref2
	}
	return ""
}

func (m *Refs) GetRef3() string {
	if m != nil {
		return m.Ref3
	}
	return ""
}

func (m *Refs) GetRef4() string {
	if m != nil {
		return m.Ref4
	}
	return ""
}

func (m *Refs) GetRef5() string {
	if m != nil {
		return m.Ref5
	}
	return ""
}

func (m *Refs) GetRef6() string {
	if m != nil {
		return m.Ref6
	}
	return ""
}

func (m *Refs) GetRef7() string {
	if m != nil {
		return m.Ref7
	}
	return ""
}

func (m *Refs) GetRef8() string {
	if m != nil {
		return m.Ref8
	}
	return ""
}

func (m *Refs) GetRef9() string {
	if m != nil {
		return m.Ref9
	}
	return ""
}

func (m *Refs) GetRef10() string {
	if m != nil {
		return m.Ref10
	}
	return ""
}

func init() {
	proto.RegisterType((*DBSession)(nil), "proto_rpc.DBSession")
	proto.RegisterType((*CreateObjectsMsg)(nil), "proto_rpc.CreateObjectsMsg")
	proto.RegisterType((*CreateIdentityMsg)(nil), "proto_rpc.CreateIdentityMsg")
	proto.RegisterType((*ObjectResponse)(nil), "proto_rpc.ObjectResponse")
	proto.RegisterType((*MultiObjectResponse)(nil), "proto_rpc.MultiObjectResponse")
	proto.RegisterType((*ObjectResponseData)(nil), "proto_rpc.ObjectResponseData")
	proto.RegisterType((*Object)(nil), "proto_rpc.Object")
	proto.RegisterType((*Refs)(nil), "proto_rpc.Refs")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for API service

type APIClient interface {
	CreateIdentity(ctx context.Context, in *CreateIdentityMsg, opts ...grpc.CallOption) (*ObjectResponse, error)
	CreateObjects(ctx context.Context, in *CreateObjectsMsg, opts ...grpc.CallOption) (*MultiObjectResponse, error)
	CreateDBSession(ctx context.Context, in *DBSession, opts ...grpc.CallOption) (*DBSession, error)
}

type aPIClient struct {
	cc *grpc.ClientConn
}

func NewAPIClient(cc *grpc.ClientConn) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) CreateIdentity(ctx context.Context, in *CreateIdentityMsg, opts ...grpc.CallOption) (*ObjectResponse, error) {
	out := new(ObjectResponse)
	err := grpc.Invoke(ctx, "/proto_rpc.API/CreateIdentity", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) CreateObjects(ctx context.Context, in *CreateObjectsMsg, opts ...grpc.CallOption) (*MultiObjectResponse, error) {
	out := new(MultiObjectResponse)
	err := grpc.Invoke(ctx, "/proto_rpc.API/CreateObjects", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) CreateDBSession(ctx context.Context, in *DBSession, opts ...grpc.CallOption) (*DBSession, error) {
	out := new(DBSession)
	err := grpc.Invoke(ctx, "/proto_rpc.API/CreateDBSession", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for API service

type APIServer interface {
	CreateIdentity(context.Context, *CreateIdentityMsg) (*ObjectResponse, error)
	CreateObjects(context.Context, *CreateObjectsMsg) (*MultiObjectResponse, error)
	CreateDBSession(context.Context, *DBSession) (*DBSession, error)
}

func RegisterAPIServer(s *grpc.Server, srv APIServer) {
	s.RegisterService(&_API_serviceDesc, srv)
}

func _API_CreateIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateIdentityMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).CreateIdentity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_rpc.API/CreateIdentity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).CreateIdentity(ctx, req.(*CreateIdentityMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_CreateObjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateObjectsMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).CreateObjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_rpc.API/CreateObjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).CreateObjects(ctx, req.(*CreateObjectsMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_CreateDBSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DBSession)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).CreateDBSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_rpc.API/CreateDBSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).CreateDBSession(ctx, req.(*DBSession))
	}
	return interceptor(ctx, in, info, handler)
}

var _API_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto_rpc.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateIdentity",
			Handler:    _API_CreateIdentity_Handler,
		},
		{
			MethodName: "CreateObjects",
			Handler:    _API_CreateObjects_Handler,
		},
		{
			MethodName: "CreateDBSession",
			Handler:    _API_CreateDBSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}

func init() { proto.RegisterFile("server.proto", fileDescriptorServer) }

var fileDescriptorServer = []byte{
	// 1131 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x98, 0x4d, 0x8f, 0xeb, 0x34,
	0x17, 0xc7, 0x95, 0x74, 0xde, 0xea, 0x79, 0xf7, 0xcc, 0xf3, 0xc8, 0x0c, 0x2f, 0x19, 0x55, 0x2c,
	0x06, 0x09, 0xe6, 0xde, 0xe9, 0x7b, 0xaf, 0x00, 0x41, 0x29, 0x88, 0x4a, 0x73, 0x35, 0xc8, 0x88,
	0x2b, 0xd0, 0x45, 0x0c, 0x69, 0xeb, 0x4e, 0xc3, 0xb4, 0x4d, 0x48, 0xdc, 0x5e, 0x95, 0x0d, 0x12,
	0x62, 0xc7, 0x06, 0x56, 0x7c, 0x2a, 0x3e, 0x00, 0x0b, 0x82, 0x04, 0xbb, 0xb2, 0xe3, 0x13, 0x20,
	0xdb, 0x69, 0x52, 0x3b, 0x71, 0x56, 0x20, 0x58, 0x4d, 0x73, 0xfe, 0xc7, 0xfe, 0x1d, 0x9f, 0xd8,
	0xe7, 0xc4, 0x03, 0xf6, 0x02, 0xe2, 0xcf, 0x89, 0x7f, 0xe9, 0xf9, 0x2e, 0x75, 0x61, 0x91, 0xff,
	0xb9, 0xf5, 0xbd, 0xfe, 0xd9, 0x6b, 0x77, 0x0e, 0x1d, 0xcd, 0x7a, 0x97, 0x7d, 0x77, 0xf2, 0xe0,
	0xce, 0xbd, 0x73, 0x1f, 0x70, 0xa9, 0x37, 0x1b, 0xf2, 0x27, 0xfe, 0xc0, 0x7f, 0x89, 0x91, 0xa5,
	0x3e, 0x28, 0x76, 0xda, 0x1f, 0x92, 0x20, 0x70, 0xdc, 0x29, 0x7c, 0x02, 0xcc, 0x6e, 0x07, 0x19,
	0xe7, 0xc6, 0x45, 0xb1, 0xfd, 0xde, 0x32, 0xb4, 0xf6, 0x9c, 0xc1, 0xab, 0xee, 0xc4, 0xa1, 0x64,
	0xe2, 0xd1, 0xc5, 0x9f, 0xa1, 0x55, 0x0e, 0xa8, 0x3f, 0xeb, 0xd3, 0xe0, 0x51, 0x69, 0x5d, 0x28,
	0x9d, 0x4f, 0x6c, 0x4f, 0x28, 0x33, 0x9f, 0x28, 0x1a, 0x36, 0xbb, 0x9d, 0xd2, 0x8f, 0x06, 0x38,
	0x7a, 0xc7, 0x27, 0x36, 0x25, 0x37, 0xbd, 0x2f, 0x48, 0x9f, 0x06, 0x8f, 0x83, 0x3b, 0xf8, 0x8d,
	0x01, 0xb6, 0x5d, 0xf1, 0x88, 0x8c, 0xf3, 0xc2, 0xc5, 0x6e, 0xf9, 0xf8, 0x32, 0x5e, 0xc6, 0xa5,
	0x70, 0x6c, 0x7f, 0xb4, 0x0c, 0xad, 0xe3, 0xc8, 0x4b, 0x0a, 0xe5, 0xf5, 0x38, 0x94, 0x94, 0xaa,
	0xc6, 0x93, 0x76, 0xc0, 0x2b, 0x70, 0xe9, 0x0f, 0x13, 0x1c, 0x8b, 0xc8, 0xba, 0x03, 0x32, 0xa5,
	0x0e, 0x5d, 0xb0, 0xd0, 0x08, 0xd8, 0x24, 0x13, 0xdb, 0x19, 0x47, 0xa9, 0xb8, 0x59, 0x86, 0xd6,
	0x21, 0x37, 0x48, 0x21, 0x34, 0xe3, 0x10, 0x14, 0x4d, 0x0d, 0x40, 0x95, 0xb1, 0x98, 0x1d, 0x52,
	0xb0, 0xe3, 0xd9, 0x41, 0xf0, 0xcc, 0xf5, 0x07, 0xc8, 0xe4, 0xa4, 0x8f, 0x97, 0xa1, 0x05, 0x57,
	0x36, 0x09, 0xf6, 0x66, 0x0c, 0x4b, 0xcb, 0x2a, 0x2f, 0xc3, 0x03, 0xc7, 0x24, 0xb8, 0x00, 0xc5,
	0x01, 0x99, 0x93, 0xb1, 0xeb, 0x11, 0x1f, 0x15, 0xce, 0x8d, 0x8b, 0x9d, 0xf6, 0xd3, 0x65, 0x68,
	0x9d, 0xc4, 0x46, 0x89, 0xfb, 0x56, 0xcc, 0xcd, 0xd0, 0x55, 0x70, 0x96, 0x0b, 0x4e, 0x68, 0xa5,
	0x9f, 0x4d, 0x70, 0x20, 0x5e, 0x2c, 0x26, 0x81, 0xe7, 0x4e, 0x03, 0x02, 0xbf, 0x06, 0x1b, 0x03,
	0x9b, 0xda, 0x3c, 0xd3, 0xbb, 0xe5, 0x17, 0x53, 0x3b, 0x60, 0xe5, 0xd8, 0xb1, 0xa9, 0xdd, 0xbe,
	0x5e, 0x86, 0xd6, 0x01, 0x73, 0x97, 0x42, 0xac, 0x27, 0x21, 0x4a, 0x52, 0x2a, 0x3a, 0x59, 0xc5,
	0x1c, 0x0c, 0xbf, 0x33, 0xc0, 0xe6, 0xd8, 0x99, 0xde, 0x07, 0xc8, 0xe4, 0x9b, 0xf0, 0x65, 0x6d,
	0x08, 0x97, 0xd7, 0xcc, 0xed, 0xdd, 0x29, 0xf5, 0x17, 0x62, 0x4b, 0xf0, 0x61, 0x9a, 0x2d, 0xa1,
	0x68, 0x6a, 0x2c, 0xaa, 0x8c, 0x45, 0x0c, 0x67, 0x4d, 0x00, 0x12, 0x0a, 0x3c, 0x02, 0x85, 0x7b,
	0xb2, 0x10, 0xbb, 0x10, 0xb3, 0x9f, 0xf0, 0x14, 0x6c, 0xce, 0xed, 0xf1, 0x8c, 0x88, 0xfd, 0x82,
	0xc5, 0xc3, 0x23, 0xb3, 0x69, 0x94, 0x7e, 0x37, 0xc1, 0xc9, 0xe3, 0xd9, 0x98, 0x3a, 0xda, 0x04,
	0x17, 0xfe, 0x9d, 0x04, 0x7f, 0xaf, 0x24, 0xf8, 0x95, 0xb5, 0x10, 0x32, 0x02, 0xfe, 0x6f, 0x67,
	0xf9, 0x37, 0x13, 0xc0, 0x74, 0xde, 0xe0, 0xe7, 0x60, 0x83, 0x2e, 0x3c, 0x12, 0xd5, 0x0b, 0x9e,
	0x45, 0xf6, 0xac, 0xc9, 0xa2, 0x2c, 0xa9, 0x41, 0x2b, 0x2a, 0xe6, 0x33, 0x47, 0xa5, 0xd9, 0xfc,
	0xbb, 0x4b, 0x33, 0xfc, 0xc1, 0x00, 0xc0, 0xa6, 0xd4, 0x77, 0x7a, 0x33, 0x4a, 0x02, 0x5e, 0x0f,
	0x32, 0x0b, 0xf1, 0x67, 0xcb, 0xd0, 0x3a, 0x4d, 0x1c, 0x25, 0x76, 0x3b, 0x66, 0x67, 0x39, 0xa8,
	0x31, 0x64, 0xfa, 0xe0, 0xb5, 0x20, 0x4a, 0x3f, 0x01, 0xb0, 0x25, 0xb0, 0xff, 0x54, 0x47, 0x82,
	0x1e, 0xd8, 0x76, 0x9f, 0x4d, 0x89, 0x1f, 0xe7, 0xf4, 0x09, 0xab, 0xbc, 0xdc, 0x74, 0xab, 0x20,
	0xd6, 0x3a, 0x0d, 0x1f, 0x91, 0xc3, 0x49, 0x3b, 0xe0, 0x15, 0x06, 0x7e, 0x05, 0x8a, 0x7d, 0xd6,
	0x68, 0x5c, 0xc6, 0x2c, 0x70, 0xe6, 0xa7, 0x2c, 0xa7, 0x91, 0x51, 0xa5, 0x26, 0x75, 0x77, 0x35,
	0x2a, 0x87, 0x9b, 0xe5, 0x82, 0x13, 0x1c, 0xfc, 0xd6, 0x00, 0xbb, 0x9e, 0xed, 0x53, 0x87, 0x3a,
	0xee, 0xb4, 0xdb, 0x41, 0x1b, 0x1c, 0xdf, 0x5b, 0x86, 0xd6, 0xff, 0x63, 0xb3, 0x1a, 0x40, 0x67,
	0xad, 0xe1, 0xac, 0x46, 0xe6, 0xf6, 0x9c, 0x2c, 0x27, 0xbc, 0x8e, 0x85, 0x4f, 0xc5, 0x41, 0xdb,
	0xe4, 0xf4, 0xee, 0x32, 0xb4, 0xf6, 0xef, 0xc9, 0x42, 0x82, 0x56, 0x63, 0xa8, 0xa4, 0xa8, 0x30,
	0x59, 0x14, 0x67, 0x96, 0xac, 0xce, 0xec, 0x56, 0xd2, 0xb3, 0xb9, 0x41, 0x53, 0x3a, 0x14, 0x4d,
	0x45, 0xa8, 0x72, 0x54, 0x04, 0x58, 0xf7, 0x64, 0x67, 0x83, 0xf4, 0x29, 0x19, 0xa0, 0xed, 0xa4,
	0x7b, 0xc6, 0x46, 0xcd, 0x5b, 0xcc, 0xd0, 0x53, 0x29, 0xcc, 0x70, 0xc1, 0x09, 0x8d, 0xed, 0x59,
	0x9f, 0x0c, 0x6f, 0xa6, 0xe3, 0x05, 0xda, 0xe1, 0x60, 0xbe, 0x67, 0x7d, 0x32, 0xbc, 0x75, 0xa7,
	0xe3, 0x85, 0x66, 0xcf, 0x46, 0x23, 0xf4, 0xd4, 0xb4, 0x03, 0x5e, 0x61, 0xd8, 0x62, 0xa9, 0x33,
	0x21, 0x01, 0xb5, 0x27, 0x1e, 0x2a, 0x9e, 0x1b, 0x17, 0x05, 0xb1, 0xd8, 0xd8, 0xa8, 0x59, 0x6c,
	0x86, 0x9e, 0xaa, 0x72, 0x19, 0x2e, 0x38, 0xa1, 0xc1, 0x19, 0xd8, 0xf1, 0x7c, 0x32, 0x7f, 0xdf,
	0x0e, 0x46, 0x08, 0xf0, 0x37, 0xfa, 0x89, 0x48, 0x33, 0x99, 0xdf, 0x8e, 0xec, 0x60, 0xa4, 0xfb,
	0x38, 0x8a, 0xc6, 0xe4, 0x65, 0x39, 0xe5, 0x81, 0x63, 0x14, 0x2b, 0xe4, 0x6c, 0x76, 0xb4, 0x9b,
	0x14, 0xf2, 0x14, 0x2d, 0x29, 0xe4, 0xa3, 0x5c, 0x92, 0xa2, 0x62, 0x3e, 0x33, 0xfc, 0x12, 0x6c,
	0x60, 0x32, 0x0c, 0xd0, 0x1e, 0xaf, 0xb4, 0x87, 0x6b, 0x95, 0x96, 0x99, 0x05, 0xd2, 0x27, 0xc3,
	0x40, 0x83, 0x94, 0xa5, 0x8c, 0x97, 0x29, 0x55, 0x55, 0x8e, 0x2a, 0xfd, 0xba, 0x2d, 0x98, 0x6c,
	0x75, 0x3e, 0x19, 0x5e, 0x71, 0x76, 0x31, 0x46, 0x5d, 0xe9, 0x51, 0x57, 0xb9, 0xa8, 0x2b, 0x09,
	0xc5, 0x0c, 0x11, 0xa1, 0x8c, 0xf6, 0x25, 0x42, 0x59, 0x4f, 0x28, 0xe7, 0x12, 0xca, 0x2a, 0xa1,
	0x1c, 0x11, 0x2a, 0xe8, 0x40, 0x22, 0x54, 0xf4, 0x84, 0x4a, 0x2e, 0xa1, 0xa2, 0x12, 0x2a, 0x11,
	0xa1, 0x8a, 0x0e, 0x25, 0x42, 0x55, 0x4f, 0xa8, 0xe6, 0x12, 0xaa, 0x2a, 0xa1, 0x1a, 0x11, 0x6a,
	0xe8, 0x48, 0x22, 0xd4, 0xf4, 0x84, 0x5a, 0x2e, 0xa1, 0xa6, 0x12, 0x6a, 0x11, 0xa1, 0x8e, 0x8e,
	0x25, 0x42, 0x5d, 0x4f, 0xa8, 0xe7, 0x12, 0xea, 0x2a, 0xa1, 0x1e, 0x11, 0x1a, 0x08, 0x4a, 0x84,
	0x86, 0x9e, 0xd0, 0xc8, 0x25, 0x34, 0x54, 0x42, 0x23, 0x22, 0x34, 0xd1, 0x89, 0x44, 0x68, 0xea,
	0x09, 0xcd, 0x5c, 0x42, 0x53, 0x25, 0x34, 0x23, 0x42, 0x0b, 0x9d, 0x4a, 0x84, 0x96, 0x9e, 0xd0,
	0xca, 0x25, 0xb4, 0x54, 0x42, 0x8b, 0x75, 0x25, 0x76, 0x2e, 0x1e, 0xa2, 0xff, 0x25, 0x5d, 0x89,
	0x1b, 0x34, 0x5d, 0x49, 0xd1, 0xb2, 0x0e, 0xdd, 0x43, 0xa9, 0x2b, 0x71, 0x4b, 0xf9, 0x17, 0x03,
	0x14, 0xde, 0xfe, 0xa0, 0x0b, 0xbb, 0xe0, 0x40, 0xbe, 0xcd, 0xc2, 0x17, 0xd6, 0x0a, 0x4c, 0xea,
	0xa2, 0x7b, 0xf6, 0x9c, 0xf6, 0x3a, 0x00, 0xaf, 0xc1, 0xbe, 0x74, 0x65, 0x87, 0xcf, 0xa7, 0x66,
	0x4a, 0x2e, 0xf3, 0x67, 0x2f, 0xe5, 0x7f, 0xd4, 0xc3, 0x37, 0xc0, 0xa1, 0x18, 0x93, 0xfc, 0xb3,
	0xe1, 0x74, 0x6d, 0x48, 0x6c, 0x3d, 0xcb, 0xb4, 0xf6, 0xb6, 0xb8, 0xb1, 0xf2, 0x57, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xb7, 0x1f, 0xc4, 0xa9, 0xf6, 0x10, 0x00, 0x00,
}
