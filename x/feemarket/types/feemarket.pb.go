// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ethermint/feemarket/v1/feemarket.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// Params defines the EVM module parameters
type Params struct {
	// no base fee forces the EIP-1559 base fee to 0 (needed for 0 price calls)
	NoBaseFee bool `protobuf:"varint,1,opt,name=no_base_fee,json=noBaseFee,proto3" json:"no_base_fee,omitempty" yaml:"no_base_fee"`
	// base fee change denominator bounds the amount the base fee can change
	// between blocks.
	BaseFeeChangeDenominator uint32 `protobuf:"varint,2,opt,name=base_fee_change_denominator,json=baseFeeChangeDenominator,proto3" json:"base_fee_change_denominator,omitempty" yaml:"base_fee_change_denominator"`
	// elasticity multiplier bounds the maximum gas limit an EIP-1559 block may
	// have.
	ElasticityMultiplier uint32 `protobuf:"varint,3,opt,name=elasticity_multiplier,json=elasticityMultiplier,proto3" json:"elasticity_multiplier,omitempty" yaml:"elasticity_multiplier"`
	// height at which the base fee calculation is enabled.
	EnableHeight int64 `protobuf:"varint,4,opt,name=enable_height,json=enableHeight,proto3" json:"enable_height,omitempty" yaml:"enable_height"`
	// base fee for EIP-1559 blocks.
	BaseFee github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=base_fee,json=baseFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"base_fee" yaml:"base_fee"`
	// min base fee for EIP-1559 blocks.
	MinBaseFee github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=min_base_fee,json=minBaseFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_base_fee" yaml:"min_base_fee"`
	// max base fee for EIP-1559 blocks.
	MaxBaseFee github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,7,opt,name=max_base_fee,json=maxBaseFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"max_base_fee" yaml:"max_base_fee"`
	// replace block max gas, if > 0
	MaxGas github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,8,opt,name=max_gas,json=maxGas,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"max_gas" yaml:"max_gas"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_4feb8b20cf98e6e1, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetNoBaseFee() bool {
	if m != nil {
		return m.NoBaseFee
	}
	return false
}

func (m *Params) GetBaseFeeChangeDenominator() uint32 {
	if m != nil {
		return m.BaseFeeChangeDenominator
	}
	return 0
}

func (m *Params) GetElasticityMultiplier() uint32 {
	if m != nil {
		return m.ElasticityMultiplier
	}
	return 0
}

func (m *Params) GetEnableHeight() int64 {
	if m != nil {
		return m.EnableHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "fx.ethermint.feemarket.v1.Params")
}

func init() {
	proto.RegisterFile("ethermint/feemarket/v1/feemarket.proto", fileDescriptor_4feb8b20cf98e6e1)
}

var fileDescriptor_4feb8b20cf98e6e1 = []byte{
	// 461 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0x4f, 0x6b, 0xd4, 0x40,
	0x18, 0xc6, 0x37, 0xd6, 0xee, 0x6e, 0xc7, 0x56, 0x21, 0xae, 0x12, 0xff, 0x90, 0x84, 0x39, 0x2c,
	0xb9, 0x34, 0xa1, 0x08, 0x1e, 0x04, 0x41, 0xe3, 0xbf, 0x7a, 0x10, 0x24, 0xe0, 0x41, 0x11, 0xc2,
	0x24, 0x7d, 0x93, 0x0c, 0xcd, 0xcc, 0x2c, 0x99, 0xd9, 0x92, 0xfd, 0x16, 0x7e, 0xac, 0x3d, 0xf6,
	0x28, 0x1e, 0x82, 0xec, 0x7e, 0x83, 0x7c, 0x02, 0xd9, 0x64, 0xbb, 0x89, 0x20, 0x42, 0xe9, 0x69,
	0xe6, 0x7d, 0x9f, 0x67, 0x9e, 0x1f, 0xbc, 0xc3, 0x8b, 0xa6, 0xa0, 0x32, 0x28, 0x18, 0xe5, 0xca,
	0x4b, 0x00, 0x18, 0x29, 0xce, 0x41, 0x79, 0x17, 0x27, 0x5d, 0xe1, 0xce, 0x0a, 0xa1, 0x84, 0xfe,
	0x28, 0x29, 0xdd, 0x9d, 0xd5, 0xed, 0xd4, 0x8b, 0x93, 0xc7, 0x93, 0x54, 0xa4, 0xa2, 0x71, 0x79,
	0x9b, 0x5b, 0xfb, 0x00, 0x2f, 0xf7, 0xd1, 0xf0, 0x33, 0x29, 0x08, 0x93, 0xfa, 0x73, 0x74, 0x87,
	0x8b, 0x30, 0x22, 0x12, 0xc2, 0x04, 0xc0, 0xd0, 0x6c, 0xcd, 0x19, 0xfb, 0x0f, 0xeb, 0xca, 0xd2,
	0x17, 0x84, 0xe5, 0x2f, 0x70, 0x4f, 0xc4, 0xc1, 0x01, 0x17, 0x3e, 0x91, 0xf0, 0x1e, 0x40, 0x07,
	0xf4, 0xe4, 0xaa, 0x1f, 0xc6, 0x19, 0xe1, 0x29, 0x84, 0x67, 0xc0, 0x05, 0xa3, 0x9c, 0x28, 0x51,
	0x18, 0xb7, 0x6c, 0xcd, 0x39, 0xf2, 0xa7, 0x75, 0x65, 0xe1, 0x36, 0xe7, 0x3f, 0x66, 0x1c, 0x18,
	0x51, 0x9b, 0xfa, 0xa6, 0xd1, 0xde, 0x76, 0x92, 0xfe, 0x05, 0x3d, 0x80, 0x9c, 0x48, 0x45, 0x63,
	0xaa, 0x16, 0x21, 0x9b, 0xe7, 0x8a, 0xce, 0x72, 0x0a, 0x85, 0xb1, 0xd7, 0x00, 0xec, 0xba, 0xb2,
	0x9e, 0xb6, 0x80, 0x7f, 0xda, 0x70, 0x30, 0xe9, 0xfa, 0x9f, 0x76, 0x6d, 0xfd, 0x25, 0x3a, 0x02,
	0x4e, 0xa2, 0x1c, 0xc2, 0x0c, 0x68, 0x9a, 0x29, 0xe3, 0xb6, 0xad, 0x39, 0x7b, 0xbe, 0x51, 0x57,
	0xd6, 0x64, 0x1b, 0xd7, 0x97, 0x71, 0x70, 0xd8, 0xd6, 0xa7, 0x4d, 0xa9, 0x7f, 0x47, 0xe3, 0xdd,
	0xc4, 0xf6, 0x6d, 0xcd, 0x39, 0xf0, 0x5f, 0x2f, 0x2b, 0x6b, 0xf0, 0xab, 0xb2, 0xa6, 0x29, 0x55,
	0xd9, 0x3c, 0x72, 0x63, 0xc1, 0xbc, 0x58, 0x48, 0x26, 0xe4, 0xf6, 0x38, 0x96, 0x67, 0xe7, 0x9e,
	0x5a, 0xcc, 0x40, 0xba, 0x1f, 0xb9, 0xaa, 0x2b, 0xeb, 0xde, 0xdf, 0x73, 0xc1, 0xc1, 0x68, 0x3b,
	0x04, 0x3d, 0x45, 0x87, 0x8c, 0xf2, 0xee, 0x4f, 0x86, 0x0d, 0xe1, 0xdd, 0xb5, 0x09, 0xf7, 0x5b,
	0x42, 0x3f, 0x0b, 0x07, 0x88, 0x51, 0xee, 0xf7, 0x40, 0xa4, 0xec, 0x40, 0xa3, 0x1b, 0x82, 0x7a,
	0x59, 0x1b, 0x10, 0x29, 0xaf, 0x40, 0x5f, 0xd1, 0x68, 0x23, 0xa6, 0x44, 0x1a, 0xe3, 0x86, 0xf1,
	0xea, 0xda, 0x8c, 0xbb, 0x1d, 0x23, 0x25, 0x12, 0x07, 0x43, 0x46, 0xca, 0x0f, 0x44, 0xfa, 0xa7,
	0xcb, 0x95, 0xa9, 0x5d, 0xae, 0x4c, 0xed, 0xf7, 0xca, 0xd4, 0x7e, 0xac, 0xcd, 0xc1, 0xe5, 0xda,
	0x1c, 0xfc, 0x5c, 0x9b, 0x83, 0x6f, 0x6e, 0x2f, 0x3b, 0x99, 0xf3, 0x58, 0x51, 0xc1, 0x4b, 0x2f,
	0x29, 0x8f, 0x63, 0x51, 0x80, 0x57, 0xf6, 0x56, 0xaa, 0xe1, 0x44, 0xc3, 0x66, 0x37, 0x9e, 0xfd,
	0x09, 0x00, 0x00, 0xff, 0xff, 0x84, 0x9e, 0xc9, 0xfd, 0x76, 0x03, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.MaxGas.Size()
		i -= size
		if _, err := m.MaxGas.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFeemarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.MaxBaseFee.Size()
		i -= size
		if _, err := m.MaxBaseFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFeemarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.MinBaseFee.Size()
		i -= size
		if _, err := m.MinBaseFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFeemarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.BaseFee.Size()
		i -= size
		if _, err := m.BaseFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintFeemarket(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.EnableHeight != 0 {
		i = encodeVarintFeemarket(dAtA, i, uint64(m.EnableHeight))
		i--
		dAtA[i] = 0x20
	}
	if m.ElasticityMultiplier != 0 {
		i = encodeVarintFeemarket(dAtA, i, uint64(m.ElasticityMultiplier))
		i--
		dAtA[i] = 0x18
	}
	if m.BaseFeeChangeDenominator != 0 {
		i = encodeVarintFeemarket(dAtA, i, uint64(m.BaseFeeChangeDenominator))
		i--
		dAtA[i] = 0x10
	}
	if m.NoBaseFee {
		i--
		if m.NoBaseFee {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintFeemarket(dAtA []byte, offset int, v uint64) int {
	offset -= sovFeemarket(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.NoBaseFee {
		n += 2
	}
	if m.BaseFeeChangeDenominator != 0 {
		n += 1 + sovFeemarket(uint64(m.BaseFeeChangeDenominator))
	}
	if m.ElasticityMultiplier != 0 {
		n += 1 + sovFeemarket(uint64(m.ElasticityMultiplier))
	}
	if m.EnableHeight != 0 {
		n += 1 + sovFeemarket(uint64(m.EnableHeight))
	}
	l = m.BaseFee.Size()
	n += 1 + l + sovFeemarket(uint64(l))
	l = m.MinBaseFee.Size()
	n += 1 + l + sovFeemarket(uint64(l))
	l = m.MaxBaseFee.Size()
	n += 1 + l + sovFeemarket(uint64(l))
	l = m.MaxGas.Size()
	n += 1 + l + sovFeemarket(uint64(l))
	return n
}

func sovFeemarket(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFeemarket(x uint64) (n int) {
	return sovFeemarket(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeemarket
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NoBaseFee", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.NoBaseFee = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFeeChangeDenominator", wireType)
			}
			m.BaseFeeChangeDenominator = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BaseFeeChangeDenominator |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ElasticityMultiplier", wireType)
			}
			m.ElasticityMultiplier = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ElasticityMultiplier |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableHeight", wireType)
			}
			m.EnableHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EnableHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeemarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeemarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BaseFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinBaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeemarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeemarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinBaseFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxBaseFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeemarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeemarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxBaseFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxGas", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeemarket
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeemarket
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeemarket
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxGas.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFeemarket(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeemarket
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
func skipFeemarket(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFeemarket
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
					return 0, ErrIntOverflowFeemarket
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
					return 0, ErrIntOverflowFeemarket
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
				return 0, ErrInvalidLengthFeemarket
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFeemarket
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFeemarket
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFeemarket        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFeemarket          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFeemarket = fmt.Errorf("proto: unexpected end of group")
)
