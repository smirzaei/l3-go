package frame

import (
	"encoding/binary"
	"fmt"
)

var (
	ErrInvalidVersion = fmt.Errorf("invalid frame version")
	ErrInvalidMsgLen  = fmt.Errorf("invalid message length")
)

type Frame struct {
	Version       uint8
	MessageLength uint32
}

func Build(data [8]byte) (Frame, error) {
	v := data[0]
	if v != 1 {
		return Frame{}, ErrInvalidVersion
	}

	msgLen := binary.LittleEndian.Uint32(data[3:])
	if msgLen == 0 {
		return Frame{}, ErrInvalidMsgLen
	}

	f := Frame{
		Version:       v,
		MessageLength: msgLen,
	}

	return f, nil
}
