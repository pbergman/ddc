package ddc

import (
	"reflect"
)

const (
	EDID_ADDR uintptr = 0x50
)

type EDIDBlock interface {
	*EDIDHeaderInfo | *EDIDDescriptor | *EDID
	Unmarshal(data []byte) error
}

type EDID struct {
	EDIDHeaderInfo
	EDIDDescriptor
}

func (e *EDID) Unmarshal(data []byte) error {
	e.EDIDHeaderInfo.Unmarshal(data)
	e.EDIDDescriptor.Unmarshal(data)
	return nil
}

type EDIDDescriptor struct {
	DisplaySerialNumber string
	UnspecifiedText     string
	DisplayName         string
}

// https://en.wikipedia.org/wiki/Extended_Display_Identification_Data#EDID_1.4_data_format
func (e *EDIDDescriptor) Unmarshal(data []byte) error {

	var blocks = [4][]byte{
		data[54:71],
		data[72:89],
		data[90:107],
		data[108:128],
	}

	for i := 0; i < 4; i++ {
		var block = blocks[i]
		var ref *string

		if block[0] == 0x00 && block[1] == 0x00 && block[2] == 0x00 && block[4] == 0x00 {
			switch block[3] {
			case 0xff:
				ref = &e.DisplaySerialNumber
			case 0xfe:
				ref = &e.UnspecifiedText
			case 0xfc:
				ref = &e.DisplayName
			}
		}

		if ref != nil {
			var value = block[5:]

			for x := 0; x < 12; x++ {
				if value[x] == '\n' {
					value = value[:x]
					break
				}
			}

			(*ref) = string(value)
		}
	}

	return nil

}

type EDIDHeaderInfo struct {
	ManufacturerId         string
	ManufactureProductCode uint16
	SerialNumber           uint32
	WeekOfManufacture      byte
	YearOfManufacture      uint16
	Version                byte
	Revision               byte
}

func (e *EDIDHeaderInfo) GetManufacturer() string {

	if v, o := PNPRegistry[e.ManufacturerId]; o {
		return e.ManufacturerId + " - " + v
	}

	return e.ManufacturerId
}

func (e *EDIDHeaderInfo) Unmarshal(data []byte) error {
	e.ManufacturerId = string([]byte{
		64 + ((data[8] >> 0x02) & 0x1f),
		64 + (((data[8] & 0x03) << 0x03) | ((data[9] >> 0x05) & 0x07)),
		64 + (data[9] & 0x1f),
	})
	e.ManufactureProductCode = uint16(data[11])<<8 | uint16(data[10])
	e.SerialNumber = uint32(data[12]) | uint32(data[13])<<8 | uint32(data[14])<<16 | uint32(data[15])<<24
	e.WeekOfManufacture = data[16]
	e.YearOfManufacture = uint16(data[17]) + 1990
	e.Version = data[18]
	e.Revision = data[19]

	return nil
}

// IsActive will try to read the first 8 bits of the EDID
// response which should be the static header. On success,
// it will take the assumption that the screen is reachable (ON)
func IsActive(h *Wire) bool {

	if err := h.SetAddress(EDID_ADDR, false); err != nil {
		return false
	}

	if _, err := h.fd.Write([]byte{0x00}); err != nil {
		return false
	}

	buf := make([]byte, 8)

	if _, err := h.fd.Read(buf); err != nil {
		return false
	}

	return isValidEDIDFixedHeader(buf)
}

// GetEDID will try to read the EDID package at address 0x50,
// for now we only support the decoding of descriptor and
// header block.
//
// handler, err := ddc.NewDisplayHandler(10)
//
// if err != nil {
//
//	   log.Fatal(err)
//	}
//
// info, err := ddc.GetEDID[*ddc.EDID](handler)
//
// or for getting only the descriptor block
//
// info, err := ddc.GetEDID[*ddc.EDIDDescriptor](handler)
//
// see: https://en.wikipedia.org/wiki/Extended_Display_Identification_Data
func GetEDID[E EDIDBlock](h *Wire) (E, error) {

	h.WriteAt(EDID_ADDR, []byte{0x00})

	buf := make([]byte, 128)

	if _, err := h.fd.Read(buf); err != nil {
		return nil, err
	}

	if IsValidEDIDData(buf) {
		var gtype E
		var ginst = reflect.New(reflect.TypeOf(gtype).Elem()).Interface()

		if v, o := ginst.(E); o {
			if err := v.Unmarshal(buf); err != nil {
				return nil, err
			}
		}

		return ginst.(E), nil
	}

	return nil, &Error{Code: ERROR_EDID_INVALID_RESPONSE, Message: "Invalid EDID response"}
}

func IsValidEDIDData(buf []byte) bool {
	return (len(buf) >= 128) && isValidEDIDFixedHeader(buf) && IsValidEDIDChecksum(buf)
}

func IsValidEDIDChecksum(buf []byte) bool {
	var out byte

	for i := 0; i < 128; i++ {
		out += buf[i]
	}

	return out == 0
}

func isValidEDIDFixedHeader(buf []byte) bool {
	return buf[0] == 0x00 &&
		buf[1] == 0xFF &&
		buf[2] == 0xFF &&
		buf[3] == 0xFF &&
		buf[4] == 0xFF &&
		buf[5] == 0xFF &&
		buf[6] == 0xFF &&
		buf[7] == 0x00
}
