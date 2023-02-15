package ddc

import "fmt"

func ddcPackageChecksum(pkg []byte) byte {
	var chk byte = pkg[0]
	for i, c := 1, len(pkg); i < c; i++ {
		chk ^= pkg[i]
	}
	return chk
}

func ddcCreatePackage(payload []byte) []byte {
	size := byte(len(payload)) + 4
	buf := make([]byte, size)
	buf[0] = 0x6e              // x37<<1 + 0   destination address, write
	buf[1] = 0x51              // x28<<1 + 1   source address
	buf[2] = (size - 4) | 0x80 // size
	if size > 4 {
		copy(buf[3:], payload)
	}
	buf[size-1] = ddcPackageChecksum(buf[:size])
	return buf[1:]
}

func isZeroSlice(slice []byte) bool {
	var s byte
	for i, c := 0, len(slice); i < c; i++ {
		s |= slice[i]
	}
	return s == 0
}

func dccValidateResponse(data []byte) error {

	var length = data[2] & 0x7f

	if isZeroSlice(data) {
		return &Error{Code: ERROR_DCC_RESP_EMPTY, Message: "received empty response"}
	}

	if data[1] == 0x6e && length == 0 && data[3] == 0xbe /** checksum */ {
		return &Error{Code: ERROR_DCC_RESP_NULL, Message: "received DDC null response"}
	}

	if data[1] != 0x6e {
		return &Error{Code: ERROR_DCC_RESP_INVALID_ADDR, Message: fmt.Sprintf("invalid address byte in response, expected 0x6e, actual %#02x", data[1])}
	}

	if length != 8 {
		return &Error{Code: ERROR_DCC_RESP_INVALID_LENGTH, Message: fmt.Sprintf("invalid query VCP response length: %d", length)}
	}

	if data[3] != 0x02 {
		return &Error{Code: ERROR_DCC_RESP_INVALID_FEATURE, Message: fmt.Sprintf("expected 0x02 in feature response field, actual value %#02x", data[3])}
	}

	data[0] = 0x50 // for calculating DDC checksum

	if chcks := ddcPackageChecksum(data[:11]); chcks != data[11] {
		return &Error{Code: ERROR_DCC_RESP_CHECKSUM, Message: fmt.Sprintf("unexpected checksum.  actual=%#02x, calculated=%#02x", data[11], chcks)}
	}

	switch data[4] {
	case 0x00:
		return nil // OK
	case 0x01:
		return &Error{Code: ERROR_DCC_RESP_UNSUPPORTED_VCP_CODE, Message: "unsupported VCP code"}
	default:
		return &Error{Code: ERROR_DCC_RESP_UNEXPECTED_VCP_CODE, Message: fmt.Sprintf("unexpected value in supported VCP code field: %#02x", data[4])}
	}
}
