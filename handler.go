package ddc

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

type (
	ERROR_CODE uint8

	Error struct {
		Code    ERROR_CODE
		Message string
	}

	DCCResponse struct {
		Code byte
		Max  uint16
		Curr uint16
	}

	DisplayHandler struct {
		bus   int
		fd    *os.File
		addr  uintptr
		retry int
	}
)

func (e *Error) Error() string {
	return e.Message
}

const (
	ERROR_DCC_RESP_EMPTY ERROR_CODE = iota + 1
	ERROR_DCC_RESP_NULL
	ERROR_DCC_RESP_INVALID_ADDR
	ERROR_DCC_RESP_INVALID_LENGTH
	ERROR_DCC_RESP_INVALID_FEATURE
	ERROR_DCC_RESP_CHECKSUM
	ERROR_DCC_RESP_UNSUPPORTED_VCP_CODE
	ERROR_DCC_RESP_UNEXPECTED_VCP_CODE
	ERROR_DCC_BUS_NOT_FOUND
	ERROR_DCC_BUS_NOT_OPEN
	ERROR_EDID_INVALID_RESPONSE
)

const (
	I2C_BUS_MAX     = 0x0020
	I2C_SLAVE       = 0x0703
	I2C_SLAVE_FORCE = 0x0706

	EDID_ADDR uintptr = 0x50
)

func NewDisplayHandler(bus int) (*DisplayHandler, error) {

	fd, err := os.OpenFile("/dev/i2c-"+strconv.Itoa(int(bus)), os.O_RDWR, 0600)

	if err != nil {

		if os.IsNotExist(err) {
			return nil, &Error{Code: ERROR_DCC_BUS_NOT_FOUND, Message: "could not find i2c bus " + strconv.Itoa(int(bus))}
		}

		return nil, &Error{Code: ERROR_DCC_BUS_NOT_OPEN, Message: "could not open i2c bus: " + err.Error()}
	}

	return &DisplayHandler{bus: bus, fd: fd, retry: 10}, nil
}

func (h *DisplayHandler) ioctl(cmd, arg uintptr) error {

	if h.addr == arg {
		return nil
	}

	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, h.fd.Fd(), cmd, arg, 0, 0, 0)

	if err != 0 {
		return err
	}

	h.addr = arg

	return nil
}

func (h *DisplayHandler) isZeroSlice(slice []byte) bool {
	var s byte

	for i, c := 0, len(slice); i < c; i++ {
		s |= slice[i]
	}

	return s == 0
}

func (h *DisplayHandler) validateDccResponse(data []byte) error {

	var length = data[2] & 0x7f

	if h.isZeroSlice(data) {
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

	if chcks := h.packageChecksum(data[:11]); chcks != data[11] {
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

func (h *DisplayHandler) Set(code uint8, value uint16) error {

	if err := h.ioctl(I2C_SLAVE, DDC_ADDR); err != nil {
		return err
	}

	_, err := h.fd.Write(h.newSetRequestPackage(code, value))

	if err != nil {
		return err
	}

	return nil
}

func (h *DisplayHandler) Get(code uint8) (*DCCResponse, error) {

	if err := h.ioctl(I2C_SLAVE, DDC_ADDR); err != nil {
		return nil, err
	}

	_, err := h.fd.Write(h.newGetRequestPackage(code))

	if err != nil {
		return nil, err
	}

	for i := 1; i <= h.retry; i++ {

		resp, err := func() (*DCCResponse, error) {
			var buf = make([]byte, 12)
			time.Sleep(20 * time.Millisecond)
			_, err = h.fd.Read(buf[1:])
			if err != nil {
				return nil, err
			}
			if err := h.validateDccResponse(buf); err != nil {
				return nil, err
			}
			return &DCCResponse{
				Code: buf[5],
				Max:  (uint16(buf[7]) << 8) + uint16(buf[8]),
				Curr: (uint16(buf[9]) << 8) + uint16(buf[10]),
			}, nil
		}()

		if err != nil {

			if i != h.retry {
				continue
			}

			return nil, err
		}

		return resp, nil
	}

	return nil, nil
}

func (h *DisplayHandler) IsCLosed() bool {
	return false == IsActive(h)
}

func (h *DisplayHandler) Close() error {
	return h.fd.Close()
}

func (h *DisplayHandler) packageChecksum(pkg []byte) byte {
	var chk byte = pkg[0]
	for i, c := 1, len(pkg); i < c; i++ {
		chk ^= pkg[i]
	}
	return chk
}

func (h *DisplayHandler) newBasePackage(payload []byte) []byte {
	size := byte(len(payload)) + 4
	buf := make([]byte, size)
	buf[0] = 0x6e              // x37<<1 + 0   destination address, write
	buf[1] = 0x51              // x28<<1 + 1   source address
	buf[2] = (size - 4) | 0x80 // size

	if size > 4 {
		copy(buf[3:], payload)
	}

	buf[size-1] = h.packageChecksum(buf[:size])

	return buf[1:]
}

func (h *DisplayHandler) newPingRequestPackage() []byte {
	return h.newBasePackage([]byte{
		0x07,
	})
}

func (h *DisplayHandler) newGetRequestPackage(code byte) []byte {
	return h.newBasePackage([]byte{
		0x01,
		code,
	})
}

func (h *DisplayHandler) newSetRequestPackage(code byte, value uint16) []byte {
	return h.newBasePackage([]byte{
		0x03,
		code,
		byte(value >> 8),
		byte(value),
	})
}
