// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services/consensus/hotstuff.proto

package consensus

import (
	context "context"
	fmt "fmt"
	math "math"

	shared "github.com/cumulus-sidechain/cumulus/hotstuff/gen/grpc/shared"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ProcessViewChangeRequest struct {
	HighestQc            *shared.QuorumCertificate `protobuf:"bytes,1,opt,name=highestQc,proto3" json:"highestQc,omitempty"`
	Signature            *shared.Signature         `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *ProcessViewChangeRequest) Reset()         { *m = ProcessViewChangeRequest{} }
func (m *ProcessViewChangeRequest) String() string { return proto.CompactTextString(m) }
func (*ProcessViewChangeRequest) ProtoMessage()    {}
func (*ProcessViewChangeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4bef6b6c621508dd, []int{0}
}

func (m *ProcessViewChangeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessViewChangeRequest.Unmarshal(m, b)
}
func (m *ProcessViewChangeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessViewChangeRequest.Marshal(b, m, deterministic)
}
func (m *ProcessViewChangeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessViewChangeRequest.Merge(m, src)
}
func (m *ProcessViewChangeRequest) XXX_Size() int {
	return xxx_messageInfo_ProcessViewChangeRequest.Size(m)
}
func (m *ProcessViewChangeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessViewChangeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessViewChangeRequest proto.InternalMessageInfo

func (m *ProcessViewChangeRequest) GetHighestQc() *shared.QuorumCertificate {
	if m != nil {
		return m.HighestQc
	}
	return nil
}

func (m *ProcessViewChangeRequest) GetSignature() *shared.Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

type ProcessBlockProposalRequest struct {
	// BlockProposal blockProposal = 1;
	Block                *shared.Block     `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	Signature            *shared.Signature `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ProcessBlockProposalRequest) Reset()         { *m = ProcessBlockProposalRequest{} }
func (m *ProcessBlockProposalRequest) String() string { return proto.CompactTextString(m) }
func (*ProcessBlockProposalRequest) ProtoMessage()    {}
func (*ProcessBlockProposalRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4bef6b6c621508dd, []int{1}
}

func (m *ProcessBlockProposalRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessBlockProposalRequest.Unmarshal(m, b)
}
func (m *ProcessBlockProposalRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessBlockProposalRequest.Marshal(b, m, deterministic)
}
func (m *ProcessBlockProposalRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessBlockProposalRequest.Merge(m, src)
}
func (m *ProcessBlockProposalRequest) XXX_Size() int {
	return xxx_messageInfo_ProcessBlockProposalRequest.Size(m)
}
func (m *ProcessBlockProposalRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessBlockProposalRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessBlockProposalRequest proto.InternalMessageInfo

func (m *ProcessBlockProposalRequest) GetBlock() *shared.Block {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *ProcessBlockProposalRequest) GetSignature() *shared.Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

type ProcessVoteRequest struct {
	BlockHash            []byte            `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Signature            *shared.Signature `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ProcessVoteRequest) Reset()         { *m = ProcessVoteRequest{} }
func (m *ProcessVoteRequest) String() string { return proto.CompactTextString(m) }
func (*ProcessVoteRequest) ProtoMessage()    {}
func (*ProcessVoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4bef6b6c621508dd, []int{2}
}

func (m *ProcessVoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessVoteRequest.Unmarshal(m, b)
}
func (m *ProcessVoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessVoteRequest.Marshal(b, m, deterministic)
}
func (m *ProcessVoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessVoteRequest.Merge(m, src)
}
func (m *ProcessVoteRequest) XXX_Size() int {
	return xxx_messageInfo_ProcessVoteRequest.Size(m)
}
func (m *ProcessVoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessVoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessVoteRequest proto.InternalMessageInfo

func (m *ProcessVoteRequest) GetBlockHash() []byte {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

func (m *ProcessVoteRequest) GetSignature() *shared.Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

type QueryBlockRequest struct {
	BlockHash            []byte   `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryBlockRequest) Reset()         { *m = QueryBlockRequest{} }
func (m *QueryBlockRequest) String() string { return proto.CompactTextString(m) }
func (*QueryBlockRequest) ProtoMessage()    {}
func (*QueryBlockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4bef6b6c621508dd, []int{3}
}

func (m *QueryBlockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryBlockRequest.Unmarshal(m, b)
}
func (m *QueryBlockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryBlockRequest.Marshal(b, m, deterministic)
}
func (m *QueryBlockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBlockRequest.Merge(m, src)
}
func (m *QueryBlockRequest) XXX_Size() int {
	return xxx_messageInfo_QueryBlockRequest.Size(m)
}
func (m *QueryBlockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBlockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBlockRequest proto.InternalMessageInfo

func (m *QueryBlockRequest) GetBlockHash() []byte {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

type QueryBlockReply struct {
	Block                *shared.Block `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *QueryBlockReply) Reset()         { *m = QueryBlockReply{} }
func (m *QueryBlockReply) String() string { return proto.CompactTextString(m) }
func (*QueryBlockReply) ProtoMessage()    {}
func (*QueryBlockReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_4bef6b6c621508dd, []int{4}
}

func (m *QueryBlockReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryBlockReply.Unmarshal(m, b)
}
func (m *QueryBlockReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryBlockReply.Marshal(b, m, deterministic)
}
func (m *QueryBlockReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBlockReply.Merge(m, src)
}
func (m *QueryBlockReply) XXX_Size() int {
	return xxx_messageInfo_QueryBlockReply.Size(m)
}
func (m *QueryBlockReply) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBlockReply.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBlockReply proto.InternalMessageInfo

func (m *QueryBlockReply) GetBlock() *shared.Block {
	if m != nil {
		return m.Block
	}
	return nil
}

func init() {
	proto.RegisterType((*ProcessViewChangeRequest)(nil), "flow.services.consensus.ProcessViewChangeRequest")
	proto.RegisterType((*ProcessBlockProposalRequest)(nil), "flow.services.consensus.ProcessBlockProposalRequest")
	proto.RegisterType((*ProcessVoteRequest)(nil), "flow.services.consensus.ProcessVoteRequest")
	proto.RegisterType((*QueryBlockRequest)(nil), "flow.services.consensus.QueryBlockRequest")
	proto.RegisterType((*QueryBlockReply)(nil), "flow.services.consensus.QueryBlockReply")
}

func init() { proto.RegisterFile("services/consensus/hotstuff.proto", fileDescriptor_4bef6b6c621508dd) }

var fileDescriptor_4bef6b6c621508dd = []byte{
	// 392 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0xcd, 0xca, 0xd3, 0x40,
	0x14, 0x25, 0x8a, 0x42, 0x6e, 0x05, 0xed, 0xe0, 0x4f, 0x49, 0x05, 0xb5, 0x1b, 0x8b, 0xc2, 0x84,
	0x56, 0xb7, 0x6e, 0x5a, 0x0b, 0x2e, 0xdb, 0x14, 0x5c, 0xb8, 0x10, 0x92, 0xf4, 0xe6, 0x07, 0xd3,
	0x4c, 0x3a, 0x77, 0xc6, 0x92, 0xa7, 0xd0, 0x47, 0xf8, 0x1e, 0xf5, 0xa3, 0x49, 0x93, 0x14, 0x9a,
	0xd0, 0xf2, 0x75, 0x99, 0xb9, 0xe7, 0x9e, 0x73, 0x72, 0xcf, 0x81, 0x0f, 0x84, 0xf2, 0x6f, 0xec,
	0x23, 0xd9, 0xbe, 0x48, 0x09, 0x53, 0xd2, 0x64, 0x47, 0x42, 0x91, 0xd2, 0x41, 0xc0, 0x33, 0x29,
	0x94, 0x60, 0x6f, 0x82, 0x44, 0xec, 0x79, 0x85, 0xe3, 0x35, 0xce, 0x7a, 0x45, 0x91, 0x2b, 0x71,
	0x63, 0x6f, 0x91, 0xc8, 0x0d, 0x91, 0x4a, 0xbc, 0x35, 0x0c, 0x85, 0x08, 0x13, 0xb4, 0x8b, 0x2f,
	0x4f, 0x07, 0x36, 0x6e, 0x33, 0x95, 0x97, 0xc3, 0xd1, 0x9d, 0x01, 0x83, 0xa5, 0x14, 0x3e, 0x12,
	0xfd, 0x8c, 0x71, 0x3f, 0x8f, 0xdc, 0x34, 0x44, 0x07, 0x77, 0x1a, 0x49, 0xb1, 0x05, 0x98, 0x51,
	0x1c, 0x46, 0x48, 0x6a, 0xe5, 0x0f, 0x8c, 0xf7, 0xc6, 0xb8, 0x37, 0xfd, 0xc8, 0x4b, 0xf5, 0x42,
	0x89, 0xd7, 0x4a, 0x2b, 0x2d, 0xa4, 0xde, 0xce, 0x51, 0xaa, 0x38, 0x88, 0x7d, 0x57, 0xa1, 0xd3,
	0x6c, 0xb2, 0x6f, 0x60, 0x52, 0x1c, 0xa6, 0xae, 0xd2, 0x12, 0x07, 0x8f, 0x0a, 0x9a, 0x77, 0xed,
	0x34, 0xeb, 0x0a, 0xe6, 0x34, 0x1b, 0xa3, 0x7f, 0x06, 0x0c, 0x8f, 0x16, 0x67, 0x89, 0xf0, 0xff,
	0x2c, 0xa5, 0xc8, 0x04, 0xb9, 0x49, 0xe5, 0x72, 0x02, 0x4f, 0xbc, 0xc3, 0xfb, 0xd1, 0xe1, 0xb0,
	0x9d, 0xba, 0x58, 0x75, 0x4a, 0xe4, 0xad, 0x8e, 0x76, 0xc0, 0xaa, 0x9b, 0x09, 0x55, 0x5f, 0xeb,
	0x2d, 0x98, 0x05, 0xfb, 0x0f, 0x97, 0xa2, 0xc2, 0xcb, 0x33, 0xa7, 0x79, 0xb8, 0x55, 0x72, 0x02,
	0xfd, 0x95, 0x46, 0x99, 0x97, 0xbf, 0x71, 0x8d, 0xe2, 0xe8, 0x3b, 0x3c, 0x3f, 0x5d, 0xc9, 0x92,
	0xfc, 0x01, 0xa7, 0x9a, 0xfe, 0x7f, 0x0c, 0x2f, 0xe6, 0x55, 0xc5, 0xd6, 0x65, 0xe9, 0xd8, 0x6f,
	0xe8, 0x9f, 0x95, 0x86, 0x4d, 0x78, 0x47, 0x31, 0x79, 0x57, 0xc1, 0xac, 0xd7, 0xbc, 0xec, 0x26,
	0xaf, 0xba, 0xc9, 0x17, 0x87, 0x6e, 0xb2, 0x0d, 0xbc, 0x6c, 0x4b, 0x9c, 0x7d, 0xbd, 0x24, 0xd1,
	0x56, 0x90, 0x4e, 0x15, 0x07, 0x7a, 0x27, 0x31, 0xb2, 0xcf, 0x17, 0xfd, 0x37, 0x61, 0x77, 0x72,
	0x7a, 0x00, 0xcd, 0xd1, 0xd9, 0xa7, 0x4e, 0xca, 0xb3, 0x30, 0xad, 0xf1, 0x55, 0xd8, 0x2c, 0xc9,
	0x67, 0xbd, 0x5f, 0x66, 0x3d, 0xf4, 0x9e, 0x16, 0x06, 0xbe, 0xdc, 0x07, 0x00, 0x00, 0xff, 0xff,
	0x87, 0x07, 0xdb, 0x3e, 0x39, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConsensusServiceClient is the client API for ConsensusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConsensusServiceClient interface {
	// Process view change messages from replicas
	ProcessViewChange(ctx context.Context, in *ProcessViewChangeRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Process block proposal from replicas
	ProcessBlockProposal(ctx context.Context, in *ProcessBlockProposalRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Process vote from replicas
	ProcessVote(ctx context.Context, in *ProcessVoteRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	QueryBlock(ctx context.Context, in *QueryBlockRequest, opts ...grpc.CallOption) (*QueryBlockReply, error)
}

type consensusServiceClient struct {
	cc *grpc.ClientConn
}

func NewConsensusServiceClient(cc *grpc.ClientConn) ConsensusServiceClient {
	return &consensusServiceClient{cc}
}

func (c *consensusServiceClient) ProcessViewChange(ctx context.Context, in *ProcessViewChangeRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/flow.services.consensus.ConsensusService/ProcessViewChange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusServiceClient) ProcessBlockProposal(ctx context.Context, in *ProcessBlockProposalRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/flow.services.consensus.ConsensusService/ProcessBlockProposal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusServiceClient) ProcessVote(ctx context.Context, in *ProcessVoteRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/flow.services.consensus.ConsensusService/ProcessVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusServiceClient) QueryBlock(ctx context.Context, in *QueryBlockRequest, opts ...grpc.CallOption) (*QueryBlockReply, error) {
	out := new(QueryBlockReply)
	err := c.cc.Invoke(ctx, "/flow.services.consensus.ConsensusService/QueryBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsensusServiceServer is the server API for ConsensusService service.
type ConsensusServiceServer interface {
	// Process view change messages from replicas
	ProcessViewChange(context.Context, *ProcessViewChangeRequest) (*empty.Empty, error)
	// Process block proposal from replicas
	ProcessBlockProposal(context.Context, *ProcessBlockProposalRequest) (*empty.Empty, error)
	// Process vote from replicas
	ProcessVote(context.Context, *ProcessVoteRequest) (*empty.Empty, error)
	QueryBlock(context.Context, *QueryBlockRequest) (*QueryBlockReply, error)
}

// UnimplementedConsensusServiceServer can be embedded to have forward compatible implementations.
type UnimplementedConsensusServiceServer struct {
}

func (*UnimplementedConsensusServiceServer) ProcessViewChange(ctx context.Context, req *ProcessViewChangeRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessViewChange not implemented")
}
func (*UnimplementedConsensusServiceServer) ProcessBlockProposal(ctx context.Context, req *ProcessBlockProposalRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessBlockProposal not implemented")
}
func (*UnimplementedConsensusServiceServer) ProcessVote(ctx context.Context, req *ProcessVoteRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessVote not implemented")
}
func (*UnimplementedConsensusServiceServer) QueryBlock(ctx context.Context, req *QueryBlockRequest) (*QueryBlockReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryBlock not implemented")
}

func RegisterConsensusServiceServer(s *grpc.Server, srv ConsensusServiceServer) {
	s.RegisterService(&_ConsensusService_serviceDesc, srv)
}

func _ConsensusService_ProcessViewChange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessViewChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServiceServer).ProcessViewChange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flow.services.consensus.ConsensusService/ProcessViewChange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServiceServer).ProcessViewChange(ctx, req.(*ProcessViewChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsensusService_ProcessBlockProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessBlockProposalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServiceServer).ProcessBlockProposal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flow.services.consensus.ConsensusService/ProcessBlockProposal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServiceServer).ProcessBlockProposal(ctx, req.(*ProcessBlockProposalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsensusService_ProcessVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServiceServer).ProcessVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flow.services.consensus.ConsensusService/ProcessVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServiceServer).ProcessVote(ctx, req.(*ProcessVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsensusService_QueryBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServiceServer).QueryBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/flow.services.consensus.ConsensusService/QueryBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServiceServer).QueryBlock(ctx, req.(*QueryBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConsensusService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "flow.services.consensus.ConsensusService",
	HandlerType: (*ConsensusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessViewChange",
			Handler:    _ConsensusService_ProcessViewChange_Handler,
		},
		{
			MethodName: "ProcessBlockProposal",
			Handler:    _ConsensusService_ProcessBlockProposal_Handler,
		},
		{
			MethodName: "ProcessVote",
			Handler:    _ConsensusService_ProcessVote_Handler,
		},
		{
			MethodName: "QueryBlock",
			Handler:    _ConsensusService_QueryBlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/consensus/hotstuff.proto",
}
