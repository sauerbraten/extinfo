package extinfo

import (
	"errors"
	"regexp"
)

type extinfoResponse struct {
	buf      []byte
	posInBuf int
}

func (r *extinfoResponse) Len() int {
	return len(r.buf)
}

func (r *extinfoResponse) HasRemaining() bool {
	return r.posInBuf < len(r.buf)
}

// returns a byte and increases the position by one
func (r *extinfoResponse) ReadByte() (byte, error) {
	if r.posInBuf < len(r.buf) {
		r.posInBuf++
		return r.buf[r.posInBuf-1], nil
	} else {
		return 0, errors.New("extinfo: buf overread!")
	}
}

// decodes the bytes read from the connection into ints
// returns the decoded byte slice as int and the amount of bytes used up of the slice
func (r *extinfoResponse) ReadInt() (value int, err error) {
	// n is the size of the buffer
	n := len(r.buf)

	if n < 1 {
		err = errors.New("extinfo: getInt: buf too short!")
		return
	}

	var b1 byte
	b1, err = r.ReadByte()
	if err != nil {
		return
	}

	// 0x80 means: value is contained in the next 2 more bytes
	if b1 == 0x80 {
		var b2, b3 byte

		b2, err = r.ReadByte()
		if err != nil {
			return
		}

		b3, err = r.ReadByte()
		if err != nil {
			return
		}

		value = int(b2) + int(b3)<<8
		return
	}

	// 0x81 means: value is contained in the next 4 more bytes
	if b1 == 0x81 {
		var b2, b3, b4, b5 byte

		b2, err = r.ReadByte()
		if err != nil {
			return
		}

		b3, err = r.ReadByte()
		if err != nil {
			return
		}

		b4, err = r.ReadByte()
		if err != nil {
			return
		}

		b5, err = r.ReadByte()
		if err != nil {
			return
		}

		value = int(b2) + int(b3)<<8 + int(b4)<<16 + int(b5)<<24
		return
	}

	// neither 0x80 nor 0x81: value was already fully contained in the first byte
	if b1 > 0x7F {
		value = int(b1) - int(1<<8)
	} else {
		value = int(b1)
	}

	return
}

// returns a string of the next bytes up to 0x00 and sets the position to the next attribute's first byte
func (r *extinfoResponse) ReadString() (s string, err error) {
	var value int
	value, err = r.ReadInt()
	if err != nil {
		return
	}

	for value != 0 {
		codepoint := uint8(value)

		s += string(cubeCodeChars[codepoint])

		value, err = r.ReadInt()
		if err != nil {
			return
		}
	}

	return sanitizeString(s), err
}

// removes C 0x00 bytes and sauer color codes from strings (especially things like \f3 etc. from server description)
func sanitizeString(s string) string {
	re := regexp.MustCompile("\\f.|\\x00")
	return re.ReplaceAllLiteralString(s, "")
}
