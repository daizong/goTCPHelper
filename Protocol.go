package TCPHelper

import (
	"encoding/binary"
	"errors"
)

type PacketProtocol struct {
	HeadSize  int              //头部字节数
	ByteOrder binary.ByteOrder //字节序
}

func DefaultPacketProtocol() *PacketProtocol {
	return &PacketProtocol{
		HeadSize:  4,
		ByteOrder: binary.BigEndian,
	}
}

func NewPacketProtocol(headSize int, byteOrder binary.ByteOrder) (*PacketProtocol, error) {
	switch headSize {
	case 1, 2, 4:
	default:
		return nil, errors.New("this head size is not supported")
	}

	return &PacketProtocol{
		HeadSize:  headSize,
		ByteOrder: byteOrder,
	}, nil
}

func (p *PacketProtocol) DecodeHead(buf []byte) (int, error) {
	switch p.HeadSize {
	case 1:
		return int(buf[0]), nil
	case 2:
		return int(p.ByteOrder.Uint16(buf)), nil
	case 4:
		return int(p.ByteOrder.Uint32(buf)), nil
	default:
		return 0, errors.New("this head size is not supported")
	}
}

func (p *PacketProtocol) EncodeHead(buf []byte, length int) error {
	switch p.HeadSize {
	case 1:
		buf[0] = byte(length)
	case 2:
		p.ByteOrder.PutUint16(buf, uint16(length))
	case 4:
		p.ByteOrder.PutUint32(buf, uint32(length))
	default:
		return errors.New("this head size is not supported")
	}

	return nil
}
