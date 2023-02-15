package ddc

type ERROR_CODE byte

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

type Error struct {
	Code    ERROR_CODE
	Message string
}

func (e *Error) Error() string {
	return e.Message
}