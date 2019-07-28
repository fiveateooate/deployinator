// Code generated by protoc-gen-go. DO NOT EDIT.
// source: deployproto.proto

package deployproto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type DeployMessage struct {
	Slug                 string   `protobuf:"bytes,1,opt,name=slug,proto3" json:"slug,omitempty"`
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Version              string   `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Domain               string   `protobuf:"bytes,4,opt,name=domain,proto3" json:"domain,omitempty"`
	Cenv                 string   `protobuf:"bytes,5,opt,name=cenv,proto3" json:"cenv,omitempty"`
	Cid                  string   `protobuf:"bytes,6,opt,name=cid,proto3" json:"cid,omitempty"`
	Deployertype         string   `protobuf:"bytes,7,opt,name=deployertype,proto3" json:"deployertype,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeployMessage) Reset()         { *m = DeployMessage{} }
func (m *DeployMessage) String() string { return proto.CompactTextString(m) }
func (*DeployMessage) ProtoMessage()    {}
func (*DeployMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e818a39b64e3e2b, []int{0}
}

func (m *DeployMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployMessage.Unmarshal(m, b)
}
func (m *DeployMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployMessage.Marshal(b, m, deterministic)
}
func (m *DeployMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployMessage.Merge(m, src)
}
func (m *DeployMessage) XXX_Size() int {
	return xxx_messageInfo_DeployMessage.Size(m)
}
func (m *DeployMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployMessage.DiscardUnknown(m)
}

var xxx_messageInfo_DeployMessage proto.InternalMessageInfo

func (m *DeployMessage) GetSlug() string {
	if m != nil {
		return m.Slug
	}
	return ""
}

func (m *DeployMessage) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *DeployMessage) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *DeployMessage) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *DeployMessage) GetCenv() string {
	if m != nil {
		return m.Cenv
	}
	return ""
}

func (m *DeployMessage) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *DeployMessage) GetDeployertype() string {
	if m != nil {
		return m.Deployertype
	}
	return ""
}

type DeployStatusMessage struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Runningversion       string   `protobuf:"bytes,2,opt,name=runningversion,proto3" json:"runningversion,omitempty"`
	Requestedversion     string   `protobuf:"bytes,3,opt,name=requestedversion,proto3" json:"requestedversion,omitempty"`
	MsgID                string   `protobuf:"bytes,4,opt,name=msgID,proto3" json:"msgID,omitempty"`
	Success              bool     `protobuf:"varint,5,opt,name=success,proto3" json:"success,omitempty"`
	Other                string   `protobuf:"bytes,6,opt,name=other,proto3" json:"other,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeployStatusMessage) Reset()         { *m = DeployStatusMessage{} }
func (m *DeployStatusMessage) String() string { return proto.CompactTextString(m) }
func (*DeployStatusMessage) ProtoMessage()    {}
func (*DeployStatusMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e818a39b64e3e2b, []int{1}
}

func (m *DeployStatusMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployStatusMessage.Unmarshal(m, b)
}
func (m *DeployStatusMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployStatusMessage.Marshal(b, m, deterministic)
}
func (m *DeployStatusMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployStatusMessage.Merge(m, src)
}
func (m *DeployStatusMessage) XXX_Size() int {
	return xxx_messageInfo_DeployStatusMessage.Size(m)
}
func (m *DeployStatusMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployStatusMessage.DiscardUnknown(m)
}

var xxx_messageInfo_DeployStatusMessage proto.InternalMessageInfo

func (m *DeployStatusMessage) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *DeployStatusMessage) GetRunningversion() string {
	if m != nil {
		return m.Runningversion
	}
	return ""
}

func (m *DeployStatusMessage) GetRequestedversion() string {
	if m != nil {
		return m.Requestedversion
	}
	return ""
}

func (m *DeployStatusMessage) GetMsgID() string {
	if m != nil {
		return m.MsgID
	}
	return ""
}

func (m *DeployStatusMessage) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *DeployStatusMessage) GetOther() string {
	if m != nil {
		return m.Other
	}
	return ""
}

type DeployResponse struct {
	Success              string   `protobuf:"bytes,1,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeployResponse) Reset()         { *m = DeployResponse{} }
func (m *DeployResponse) String() string { return proto.CompactTextString(m) }
func (*DeployResponse) ProtoMessage()    {}
func (*DeployResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e818a39b64e3e2b, []int{2}
}

func (m *DeployResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployResponse.Unmarshal(m, b)
}
func (m *DeployResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployResponse.Marshal(b, m, deterministic)
}
func (m *DeployResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployResponse.Merge(m, src)
}
func (m *DeployResponse) XXX_Size() int {
	return xxx_messageInfo_DeployResponse.Size(m)
}
func (m *DeployResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeployResponse proto.InternalMessageInfo

func (m *DeployResponse) GetSuccess() string {
	if m != nil {
		return m.Success
	}
	return ""
}

func init() {
	proto.RegisterType((*DeployMessage)(nil), "deployproto.DeployMessage")
	proto.RegisterType((*DeployStatusMessage)(nil), "deployproto.DeployStatusMessage")
	proto.RegisterType((*DeployResponse)(nil), "deployproto.DeployResponse")
}

func init() { proto.RegisterFile("deployproto.proto", fileDescriptor_7e818a39b64e3e2b) }

var fileDescriptor_7e818a39b64e3e2b = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0xc1, 0x4e, 0x02, 0x31,
	0x10, 0x86, 0x5d, 0x81, 0x45, 0x46, 0x20, 0x38, 0x1a, 0xd2, 0x10, 0x0f, 0x64, 0x0f, 0xc6, 0x70,
	0xe0, 0xa0, 0xaf, 0xc0, 0xc5, 0x83, 0x9a, 0xa0, 0x2f, 0xb0, 0xee, 0x4e, 0xd6, 0x26, 0xd0, 0xae,
	0x9d, 0x2e, 0x09, 0x4f, 0xe5, 0x13, 0x78, 0xf6, 0xb5, 0xcc, 0xb6, 0xdb, 0xc8, 0xaa, 0x37, 0x2f,
	0xcd, 0xfc, 0xff, 0xdf, 0x74, 0xe6, 0xeb, 0xc0, 0x59, 0x4e, 0xe5, 0x46, 0xef, 0x4b, 0xa3, 0xad,
	0x5e, 0xba, 0x13, 0x4f, 0x0f, 0xac, 0xe4, 0x23, 0x82, 0xd1, 0xca, 0xe9, 0x7b, 0x62, 0x4e, 0x0b,
	0x42, 0x84, 0x2e, 0x6f, 0xaa, 0x42, 0x44, 0xf3, 0xe8, 0x7a, 0xb0, 0x76, 0x35, 0x5e, 0xc2, 0x40,
	0xa5, 0x5b, 0xe2, 0x32, 0xcd, 0x48, 0x1c, 0xbb, 0xe0, 0xdb, 0x40, 0x01, 0xfd, 0x1d, 0x19, 0x96,
	0x5a, 0x89, 0x8e, 0xcb, 0x82, 0xc4, 0x29, 0xc4, 0xb9, 0xde, 0xa6, 0x52, 0x89, 0xae, 0x0b, 0x1a,
	0x55, 0xf7, 0xc8, 0x48, 0xed, 0x44, 0xcf, 0xf7, 0xa8, 0x6b, 0x9c, 0x40, 0x27, 0x93, 0xb9, 0x88,
	0x9d, 0x55, 0x97, 0x98, 0xc0, 0xd0, 0x8f, 0x4a, 0xc6, 0xee, 0x4b, 0x12, 0x7d, 0x17, 0xb5, 0xbc,
	0xe4, 0x33, 0x82, 0x73, 0x3f, 0xff, 0x93, 0x4d, 0x6d, 0xc5, 0x81, 0x62, 0x0a, 0x31, 0x3b, 0xa3,
	0xe1, 0x68, 0x14, 0x5e, 0xc1, 0xd8, 0x54, 0x4a, 0x49, 0x55, 0x84, 0x91, 0x3d, 0xce, 0x0f, 0x17,
	0x17, 0x30, 0x31, 0xf4, 0x56, 0x11, 0x5b, 0xca, 0xdb, 0x70, 0xbf, 0x7c, 0xbc, 0x80, 0xde, 0x96,
	0x8b, 0xbb, 0x55, 0x03, 0xe9, 0x45, 0xfd, 0x2b, 0x5c, 0x65, 0x19, 0x31, 0x3b, 0xcc, 0x93, 0x75,
	0x90, 0xf5, 0x7d, 0x6d, 0x5f, 0xc9, 0x34, 0xac, 0x5e, 0x24, 0x0b, 0x18, 0x7b, 0x90, 0x35, 0x71,
	0xa9, 0x15, 0xd3, 0xe1, 0x0b, 0x1e, 0x22, 0xc8, 0x9b, 0xf7, 0x08, 0x86, 0xfe, 0xb2, 0x54, 0xa9,
	0xd5, 0x06, 0x1f, 0x61, 0xf4, 0x6c, 0x64, 0x51, 0x90, 0xf1, 0x36, 0xce, 0x96, 0x87, 0x8b, 0x6f,
	0x6d, 0x78, 0x36, 0xff, 0x23, 0x6b, 0xfd, 0x5e, 0x72, 0x84, 0x0f, 0xa1, 0x81, 0x0f, 0xfe, 0xfb,
	0xde, 0x4b, 0xec, 0xc2, 0xdb, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x92, 0x32, 0x3e, 0x0e, 0x90,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DeployinatorClient is the client API for Deployinator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DeployinatorClient interface {
	TriggerDeploy(ctx context.Context, in *DeployMessage, opts ...grpc.CallOption) (*DeployStatusMessage, error)
	DeployStatus(ctx context.Context, in *DeployMessage, opts ...grpc.CallOption) (*DeployStatusMessage, error)
}

type deployinatorClient struct {
	cc *grpc.ClientConn
}

func NewDeployinatorClient(cc *grpc.ClientConn) DeployinatorClient {
	return &deployinatorClient{cc}
}

func (c *deployinatorClient) TriggerDeploy(ctx context.Context, in *DeployMessage, opts ...grpc.CallOption) (*DeployStatusMessage, error) {
	out := new(DeployStatusMessage)
	err := c.cc.Invoke(ctx, "/deployproto.Deployinator/TriggerDeploy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deployinatorClient) DeployStatus(ctx context.Context, in *DeployMessage, opts ...grpc.CallOption) (*DeployStatusMessage, error) {
	out := new(DeployStatusMessage)
	err := c.cc.Invoke(ctx, "/deployproto.Deployinator/DeployStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeployinatorServer is the server API for Deployinator service.
type DeployinatorServer interface {
	TriggerDeploy(context.Context, *DeployMessage) (*DeployStatusMessage, error)
	DeployStatus(context.Context, *DeployMessage) (*DeployStatusMessage, error)
}

// UnimplementedDeployinatorServer can be embedded to have forward compatible implementations.
type UnimplementedDeployinatorServer struct {
}

func (*UnimplementedDeployinatorServer) TriggerDeploy(ctx context.Context, req *DeployMessage) (*DeployStatusMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerDeploy not implemented")
}
func (*UnimplementedDeployinatorServer) DeployStatus(ctx context.Context, req *DeployMessage) (*DeployStatusMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployStatus not implemented")
}

func RegisterDeployinatorServer(s *grpc.Server, srv DeployinatorServer) {
	s.RegisterService(&_Deployinator_serviceDesc, srv)
}

func _Deployinator_TriggerDeploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeployinatorServer).TriggerDeploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deployproto.Deployinator/TriggerDeploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeployinatorServer).TriggerDeploy(ctx, req.(*DeployMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deployinator_DeployStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeployinatorServer).DeployStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deployproto.Deployinator/DeployStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeployinatorServer).DeployStatus(ctx, req.(*DeployMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _Deployinator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "deployproto.Deployinator",
	HandlerType: (*DeployinatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TriggerDeploy",
			Handler:    _Deployinator_TriggerDeploy_Handler,
		},
		{
			MethodName: "DeployStatus",
			Handler:    _Deployinator_DeployStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "deployproto.proto",
}
