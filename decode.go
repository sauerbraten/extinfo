package extinfo

import (
	"errors"
	"regexp"
)

// the current position in a response ([]byte)
// needed, since values are encoded in variable amount of bytes
// global to not have to pass around an int on every dump
var positionInResponse int

// decodes the bytes read from the connection into ints
// returns the decoded byte slice as int and the amount of bytes used up of the slice
func getInt(buf []byte) (value int, bytesRead int, err error) {
	// n is the size of the buffer
	n := len(buf)

	if n < 1 {
		err = errors.New("extinfo: getInt: buf too short!")
		return
	}

	// b is the first byte in buf
	b := buf[0]

	bytesRead = 1

	// 0x80 means: value is contained in the next 2 more bytes
	if b == 0x80 {
		if n < 3 {
			err = errors.New("extinfo: getInt: buf too short!")
			return
		}

		// return the decoded int (contained in the next 2 bytes) and the amount of bytes used
		value = int(buf[1]) + int(buf[2])<<8
		bytesRead += 2
		return
	}

	// 0x81 means: value is contained in the next 4 more bytes
	if b == 0x81 {
		if n < 5 {
			err = errors.New("extinfo: getInt: buf too short!")
			return
		}

		// return the decoded int (contained in the next 4 bytes) and the amount of bytes used
		value = int(buf[1]) + int(buf[2])<<8 + int(buf[3])<<16 + int(buf[4])<<24
		bytesRead += 4
		return
	}

	// value was already fully contained in the first byte
	if b > 0x7F {
		value = int(b) - int(1<<8)
	} else {
		value = int(b)
	}

	return
}

// returns a byte and increases the position by one
func dumpByte(buf []byte) byte {
	positionInResponse++
	return buf[positionInResponse-1]
}

// returns a decoded int and sets the position to the next attribute's first byte
func dumpInt(buf []byte) (int, error) {
	decodedInt, bytesRead, err := getInt(buf[positionInResponse:])
	positionInResponse = positionInResponse + bytesRead
	return decodedInt, err
}

// returns a string of the next bytes up to 0x00 and sets the position to the next attribute's first byte
func dumpString(buf []byte) (s string, err error) {
	value := -1

	for value != 0 {
		value, err = dumpInt(buf)
		if err != nil {
			return
		}

		// convert to 8-bit uint for lookup in cubecode table
		codepoint := uint8(value)

		s += string(cubeCodeChars[codepoint])
	}

	return
}

// removes C 0x00 bytes and sauer color codes from strings (especially things like \f3 etc. from server description)
func sanitizeString(s string) string {
	re := regexp.MustCompile("\\f.|\\x00")
	return re.ReplaceAllLiteralString(s, "")
}
