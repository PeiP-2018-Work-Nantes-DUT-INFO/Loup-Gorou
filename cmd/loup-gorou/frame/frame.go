package frame

import (
	"encoding/binary"
	"fmt"

	"github.com/smallnest/goframe"
)

func WriteFrame(encoderConfig goframe.EncoderConfig, p []byte) (out []byte, err1 error) {
	length := len(p) + encoderConfig.LengthAdjustment
	if encoderConfig.LengthIncludesLengthFieldLength {
		length += encoderConfig.LengthFieldLength
	}

	if length < 0 {
		err1 = goframe.ErrTooLessLength
		return
	}

	switch encoderConfig.LengthFieldLength {
	case 1:
		if length >= 256 {
			err1 = fmt.Errorf("length does not fit into a byte: %d", length)
		}
		out = append(out, byte(length))
	case 2:
		if length >= 65536 {
			err1 = fmt.Errorf("length does not fit into a short integer: %d", length)
			return
		}
		encoderConfig.ByteOrder.PutUint16(out, uint16(length))
	case 3:
		if length >= 16777216 {
			err1 = fmt.Errorf("length does not fit into a medium integer: %d", length)
			return
		}
		out = writeUint24(encoderConfig.ByteOrder, length)
	case 4:
		encoderConfig.ByteOrder.PutUint32(out, uint32(length))
	case 8:
		encoderConfig.ByteOrder.PutUint64(out, uint64(length))
	default:
		err1 = goframe.ErrUnsupportedlength
		return
	}
	out = append(out, p...)
	return
}
func writeUint24(byteOrder binary.ByteOrder, v int) []byte {
	b := make([]byte, 3)
	if byteOrder == binary.LittleEndian {
		b[0] = byte(v)
		b[1] = byte(v >> 8)

		b[2] = byte(v >> 16)
	} else {
		b[2] = byte(v)
		b[1] = byte(v >> 8)
		b[0] = byte(v >> 16)
	}
	return b
}
