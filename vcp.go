package ddc

import "time"

const (
	DDC_ADDR uintptr = 0x37
)

type VCPResponse struct {
	Code byte
	Max  uint16
	Curr uint16
}

func SetVCP(w *Wire, index byte, value uint16) error {

	_, err := w.WriteAt(DDC_ADDR, ddcCreatePackage([]byte{
		0x03,
		index,
		byte(value >> 8),
		byte(value),
	}))

	if err != nil {
		return err
	}

	return nil
}

func GetVCP(w *Wire, index byte) (*VCPResponse, error) {

	_, err := w.WriteAt(DDC_ADDR, ddcCreatePackage([]byte{
		0x01,
		index,
	}))

	if err != nil {
		return nil, err
	}

	time.Sleep(20 * time.Millisecond)

	var buf = make([]byte, 12)

	_, err = w.Read(buf[1:])

	if err != nil {
		return nil, err
	}

	if err := dccValidateResponse(buf); err != nil {
		return nil, err
	}

	return &VCPResponse{
		Code: buf[5],
		Max:  (uint16(buf[7]) << 8) + uint16(buf[8]),
		Curr: (uint16(buf[9]) << 8) + uint16(buf[10]),
	}, nil
}
