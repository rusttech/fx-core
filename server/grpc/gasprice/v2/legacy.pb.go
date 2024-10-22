// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/base/v1/legacy.proto

package v2

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Deprecated: after upgrade v4
type GetGasPriceRequest struct {
}

func (m *GetGasPriceRequest) Reset()         { *m = GetGasPriceRequest{} }
func (m *GetGasPriceRequest) String() string { return proto.CompactTextString(m) }
func (*GetGasPriceRequest) ProtoMessage()    {}
func (*GetGasPriceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a72c8556a163ae4f, []int{0}
}
func (m *GetGasPriceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetGasPriceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetGasPriceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetGasPriceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetGasPriceRequest.Merge(m, src)
}
func (m *GetGasPriceRequest) XXX_Size() int {
	return m.Size()
}
func (m *GetGasPriceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetGasPriceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetGasPriceRequest proto.InternalMessageInfo

// Deprecated: after upgrade v4
type GetGasPriceResponse struct {
	GasPrices github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=gas_prices,json=gasPrices,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"gas_prices" yaml:"gas_prices"`
}

func (m *GetGasPriceResponse) Reset()         { *m = GetGasPriceResponse{} }
func (m *GetGasPriceResponse) String() string { return proto.CompactTextString(m) }
func (*GetGasPriceResponse) ProtoMessage()    {}
func (*GetGasPriceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a72c8556a163ae4f, []int{1}
}
func (m *GetGasPriceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetGasPriceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetGasPriceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetGasPriceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetGasPriceResponse.Merge(m, src)
}
func (m *GetGasPriceResponse) XXX_Size() int {
	return m.Size()
}
func (m *GetGasPriceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetGasPriceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetGasPriceResponse proto.InternalMessageInfo

func (m *GetGasPriceResponse) GetGasPrices() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.GasPrices
	}
	return nil
}

func init() {
	proto.RegisterType((*GetGasPriceRequest)(nil), "fx.base.v1.GetGasPriceRequest")
	proto.RegisterType((*GetGasPriceResponse)(nil), "fx.base.v1.GetGasPriceResponse")
}

func init() { proto.RegisterFile("fx/base/v1/legacy.proto", fileDescriptor_a72c8556a163ae4f) }

var fileDescriptor_a72c8556a163ae4f = []byte{
	// 356 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x86, 0x93, 0x7b, 0xb9, 0x17, 0x9c, 0xae, 0x8c, 0x95, 0x6a, 0xd1, 0x89, 0x64, 0xd5, 0x4d,
	0x67, 0x68, 0x75, 0xe5, 0xb2, 0x22, 0xdd, 0x89, 0x76, 0xe9, 0x46, 0x26, 0xe3, 0x64, 0x0c, 0x36,
	0x99, 0x98, 0x33, 0x09, 0xc9, 0x4a, 0xf0, 0x09, 0x04, 0xc1, 0x87, 0xf0, 0x49, 0xba, 0x2c, 0xb8,
	0x71, 0x55, 0xa5, 0xf5, 0x09, 0x7c, 0x02, 0x49, 0xd2, 0xda, 0x8a, 0xb8, 0xca, 0xe1, 0xfc, 0x27,
	0x3f, 0xdf, 0xff, 0x0f, 0x6a, 0x78, 0x19, 0x75, 0x19, 0x08, 0x9a, 0x76, 0xe8, 0x50, 0x48, 0xc6,
	0x73, 0x12, 0xc5, 0x4a, 0x2b, 0x0b, 0x79, 0x19, 0x29, 0x04, 0x92, 0x76, 0x9a, 0x98, 0x2b, 0x08,
	0x14, 0x2c, 0x0e, 0x5d, 0xa1, 0x59, 0x87, 0x72, 0xe5, 0x87, 0xd5, 0x6d, 0xb3, 0x2e, 0x95, 0x54,
	0xe5, 0x48, 0x8b, 0x69, 0xbe, 0xdd, 0x91, 0x4a, 0xc9, 0xa1, 0xa0, 0x2c, 0xf2, 0x29, 0x0b, 0x43,
	0xa5, 0x99, 0xf6, 0x55, 0x08, 0x95, 0xea, 0xd4, 0x91, 0xd5, 0x17, 0xba, 0xcf, 0xe0, 0x34, 0xf6,
	0xb9, 0x18, 0x88, 0x9b, 0x44, 0x80, 0x76, 0x1e, 0x4d, 0xb4, 0xf1, 0x6d, 0x0d, 0x91, 0x0a, 0x41,
	0x58, 0xb7, 0x08, 0x49, 0x06, 0x17, 0x51, 0xb1, 0x84, 0x2d, 0x73, 0xef, 0x6f, 0xab, 0xd6, 0xdd,
	0x26, 0x15, 0xd6, 0x02, 0xb3, 0xc4, 0x22, 0x47, 0xca, 0x0f, 0x7b, 0xc7, 0xa3, 0x89, 0x6d, 0x7c,
	0x4c, 0xec, 0xf5, 0x9c, 0x05, 0xc3, 0x43, 0x67, 0xf9, 0xab, 0xf3, 0xf4, 0x6a, 0xb7, 0xa4, 0xaf,
	0xaf, 0x12, 0x97, 0x70, 0x15, 0xd0, 0x79, 0xb0, 0xea, 0xd3, 0x86, 0xcb, 0x6b, 0xaa, 0xf3, 0x48,
	0x40, 0xe9, 0x02, 0x83, 0x35, 0x39, 0xe7, 0x80, 0x6e, 0x8a, 0xfe, 0x9d, 0x25, 0x22, 0xce, 0xad,
	0x00, 0xd5, 0x56, 0x00, 0x2d, 0x4c, 0x96, 0x3d, 0x91, 0x9f, 0x81, 0x9a, 0xf6, 0xaf, 0x7a, 0x95,
	0xcc, 0xd9, 0xbd, 0x7b, 0x7e, 0x7f, 0xf8, 0xd3, 0xb0, 0x36, 0xe9, 0xca, 0x4b, 0x7c, 0x01, 0xf7,
	0x4e, 0x46, 0x53, 0x6c, 0x8e, 0xa7, 0xd8, 0x7c, 0x9b, 0x62, 0xf3, 0x7e, 0x86, 0x8d, 0xf1, 0x0c,
	0x1b, 0x2f, 0x33, 0x6c, 0x9c, 0x1f, 0xac, 0xc4, 0xf0, 0x92, 0x90, 0x17, 0xe5, 0x66, 0xd4, 0xcb,
	0xda, 0x5c, 0xc5, 0x82, 0x82, 0x88, 0x53, 0x11, 0x53, 0x19, 0x47, 0xbc, 0x70, 0x2b, 0xcd, 0x68,
	0xda, 0x75, 0xff, 0x97, 0xed, 0xef, 0x7f, 0x06, 0x00, 0x00, 0xff, 0xff, 0x28, 0xbc, 0xb9, 0xe2,
	0xf8, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Deprecated: please use cosmos.base.node.v1beta1.Service.Config
	GetGasPrice(ctx context.Context, in *GetGasPriceRequest, opts ...grpc.CallOption) (*GetGasPriceResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GetGasPrice(ctx context.Context, in *GetGasPriceRequest, opts ...grpc.CallOption) (*GetGasPriceResponse, error) {
	out := new(GetGasPriceResponse)
	err := c.cc.Invoke(ctx, "/fx.base.v1.Query/GetGasPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Deprecated: please use cosmos.base.node.v1beta1.Service.Config
	GetGasPrice(context.Context, *GetGasPriceRequest) (*GetGasPriceResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) GetGasPrice(ctx context.Context, req *GetGasPriceRequest) (*GetGasPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGasPrice not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GetGasPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGasPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetGasPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.base.v1.Query/GetGasPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetGasPrice(ctx, req.(*GetGasPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fx.base.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGasPrice",
			Handler:    _Query_GetGasPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fx/base/v1/legacy.proto",
}

func (m *GetGasPriceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetGasPriceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetGasPriceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *GetGasPriceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetGasPriceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetGasPriceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.GasPrices) > 0 {
		for iNdEx := len(m.GasPrices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.GasPrices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLegacy(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintLegacy(dAtA []byte, offset int, v uint64) int {
	offset -= sovLegacy(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GetGasPriceRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *GetGasPriceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.GasPrices) > 0 {
		for _, e := range m.GasPrices {
			l = e.Size()
			n += 1 + l + sovLegacy(uint64(l))
		}
	}
	return n
}

func sovLegacy(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLegacy(x uint64) (n int) {
	return sovLegacy(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GetGasPriceRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLegacy
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GetGasPriceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetGasPriceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipLegacy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLegacy
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GetGasPriceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLegacy
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GetGasPriceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetGasPriceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLegacy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLegacy
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLegacy
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GasPrices = append(m.GasPrices, types.Coin{})
			if err := m.GasPrices[len(m.GasPrices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLegacy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLegacy
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipLegacy(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLegacy
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLegacy
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLegacy
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthLegacy
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLegacy
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLegacy
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLegacy        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLegacy          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLegacy = fmt.Errorf("proto: unexpected end of group")
)