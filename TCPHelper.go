/*
   a package to help to handle data for a conn(tcp or unix)
*/
package TCPHelper

import (
	"bytes"
	"errors"
	"io"
	"net"
)

type Helper struct {
	Conn net.Conn
	PacketProtocol
}

func NewHelper(conn net.Conn, protocol *PacketProtocol) *Helper {
	return &Helper{
		Conn:           conn,
		PacketProtocol: *protocol,
	}
}

func (h *Helper) ReadLoop(handle func([]byte)) {
	for {
		headBuf := make([]byte, h.PacketProtocol.HeadSize)
		_, err := io.ReadFull(h.Conn, headBuf)
		if err != nil && err.Error() != "EOF" {
			panic(err)
		}

		bodyLen, err := h.DecodeHead(headBuf)
		if err != nil {
			panic(err)
		}

		if bodyLen <= 0 {
			continue
		}

		body := make([]byte, bodyLen)
		_, err = io.ReadFull(h.Conn, body)
		if err != nil && err.Error() != "EOF" {
			panic(err)
		}
		go handle(body)
	}
}

func (h *Helper) Write(buf []byte) (int, error) {
	bodyLen := len(buf)
	if bodyLen == 0 {
		return 0, errors.New("buf is empty")
	}
	head := make([]byte, h.HeadSize)
	buffer := bytes.NewBuffer(head)
	buffer.Write(buf)
	buf = buffer.Bytes()

	err := h.EncodeHead(buf, bodyLen)
	if err != nil {
		panic(err)
	}

	return h.Conn.Write(buf)
}
