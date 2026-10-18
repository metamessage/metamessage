package core

import (
	"errors"
)

// ensureCap ensures the internal buffer has enough capacity for required bytes
// starting from the current offset. Returns ErrMaxSizeExceeded if the required
// size exceeds maxCap.
func (e *encoder) ensureCap(required int) error {
	if required > maxCap {
		return ErrMaxSizeExceeded
	}

	if required > cap(e.buf) {
		newCap := cap(e.buf) * 2
		if newCap > maxCap || newCap < required {
			newCap = required
		}
		newBuf := make([]byte, newCap)
		copy(newBuf, e.buf[:e.offset])
		e.buf = newBuf
	}

	return nil
}

// for io.ByteWriter
func (e *encoder) WriteByte(c byte) error {
	_, err := e.writeByte(c)
	return err
}

// for io.Writer
func (e *encoder) Write(bs []byte) (n int, err error) {
	var u32 uint32
	u32, err = e.writeBytes(bs)
	return int(u32), err
}

func (e *encoder) writeBytes(bs []byte) (n uint32, err error) {
	l := len(bs)
	if l == 0 {
		return
	}

	required := int(e.offset) + l
	if err = e.ensureCap(required); err != nil {
		return
	}

	copy(e.buf[e.offset:], bs)
	n = uint32(l)
	e.offset += n
	return
}

var ErrMaxSizeExceeded = errors.New("maximum size exceeded")

func (e *encoder) writeByte(bs ...byte) (n uint32, err error) {
	l := len(bs)
	required := int(e.offset) + l
	if err = e.ensureCap(required); err != nil {
		return
	}

	copy(e.buf[e.offset:], bs)
	n = uint32(l)
	e.offset += n
	return
}

func (e *encoder) writeString(s string) (n uint32, err error) {
	l := len(s)
	required := int(e.offset) + l
	if err = e.ensureCap(required); err != nil {
		return
	}

	copy(e.buf[e.offset:], s)
	n = uint32(l)
	e.offset += n
	return
}

func (e *encoder) writeBytesWithPrefix(bs []byte, prefix ...byte) (n uint32, err error) {
	lp := len(prefix)
	l := lp + len(bs)
	required := int(e.offset) + l
	if err = e.ensureCap(required); err != nil {
		return
	}

	copy(e.buf[e.offset:], prefix)
	copy(e.buf[e.offset+uint32(lp):], bs)
	n = uint32(l)
	e.offset += n
	return
}

func (e *encoder) writeStringWithPrefix(s string, prefix ...byte) (n uint32, err error) {
	lp := len(prefix)
	l := lp + len(s)
	required := int(e.offset) + l
	if err = e.ensureCap(required); err != nil {
		return
	}

	copy(e.buf[e.offset:], prefix)
	copy(e.buf[e.offset+uint32(lp):], s)
	n = uint32(l)

e.offset += n
	return
}

// v0.5.0 尾定界写出原语：body 先出，定界字节（档位长度大端 + sign）追加在 body 之后。
// 编码器一趟正向写完 body 即可写长度与 sign，无需回写、无需事后反转。
func (e *encoder) writeBytesWithSuffix(bs []byte, suffix ...byte) (n uint32, err error) {
	lb := len(bs)
	l := lb + len(suffix)
	required := int(e.offset) + l
	if err = e.ensureCap(required); err != nil {
		return
	}

	copy(e.buf[e.offset:], bs)
	copy(e.buf[e.offset+uint32(lb):], suffix)
	n = uint32(l)

e.offset += n
	return
}

func (e *encoder) writeStringWithSuffix(s string, suffix ...byte) (n uint32, err error) {
	ls := len(s)
	l := ls + len(suffix)
	required := int(e.offset) + l
	if err = e.ensureCap(required); err != nil {
		return
	}

	copy(e.buf[e.offset:], s)
	copy(e.buf[e.offset+uint32(ls):], suffix)
	n = uint32(l)
	e.offset += n
	return
}
