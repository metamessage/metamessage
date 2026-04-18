package mm

import "errors"

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
	if required > maxCap {
		err = ErrMaxSizeExceeded
		return
	}

	if int(e.offset)+l > cap(e.buf) {
		newCap := cap(e.buf) * 2
		if newCap > maxCap || newCap < required {
			newCap = required
		}
		newBuf := make([]byte, len(e.buf), newCap)
		copy(newBuf, e.buf)
		e.buf = newBuf
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
	if required > maxCap {
		err = ErrMaxSizeExceeded
		return
	}

	if int(e.offset)+l > cap(e.buf) {
		newCap := cap(e.buf) * 2
		if newCap > maxCap || newCap < required {
			newCap = required
		}
		newBuf := make([]byte, len(e.buf), newCap)
		copy(newBuf, e.buf)
		e.buf = newBuf
	}

	copy(e.buf[e.offset:], bs)
	n = uint32(l)
	e.offset += n
	return
}

func (e *encoder) writeString(s string) (n uint32, err error) {
	l := len(s)
	required := int(e.offset) + l
	if required > maxCap {
		err = ErrMaxSizeExceeded
		return
	}

	if int(e.offset)+l > cap(e.buf) {
		newCap := cap(e.buf) * 2
		if newCap > maxCap || newCap < required {
			newCap = required
		}
		newBuf := make([]byte, len(e.buf), newCap)
		copy(newBuf, e.buf)
		e.buf = newBuf
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
	if required > maxCap {
		err = ErrMaxSizeExceeded
		return
	}

	if int(e.offset)+l > cap(e.buf) {
		newCap := cap(e.buf) * 2
		if newCap > maxCap || newCap < required {
			newCap = required
		}
		newBuf := make([]byte, len(e.buf), newCap)
		copy(newBuf, e.buf)
		e.buf = newBuf
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
	if required > maxCap {
		err = ErrMaxSizeExceeded
		return
	}

	if int(e.offset)+l > cap(e.buf) {
		newCap := cap(e.buf) * 2
		if newCap > maxCap || newCap < required {
			newCap = required
		}
		newBuf := make([]byte, len(e.buf), newCap)
		copy(newBuf, e.buf)
		e.buf = newBuf
	}

	copy(e.buf[e.offset:], prefix)
	copy(e.buf[e.offset+uint32(lp):], s)
	n = uint32(l)
	e.offset += n
	return
}
