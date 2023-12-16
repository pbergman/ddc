package ddc

import "time"

const (
	DDC_ADDR uintptr = 0x37
)

type VCPResponse struct {
	Code byte
	MH   byte // Maximum value High byte
	ML   byte // Maximum value Low byte
	SH   byte // Present value High byte
	SL   byte // Present value Low byte
}

func (v VCPResponse) GetMax() uint16 {
	return (uint16(v.MH) << 8) + uint16(v.ML)
}

func (v VCPResponse) GetCurr() uint16 {
	return (uint16(v.SH) << 8) + uint16(v.SL)
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

	time.Sleep(w.sleep)

	var buf = make([]byte, 12)

	_, err = w.Read(buf[1:])

	if err != nil {
		return nil, err
	}

	if err := dccValidateResponse(buf); err != nil {
		return nil, err
	}

	return &VCPResponse{
		Code: buf[0x05],
		MH:   buf[0x07],
		ML:   buf[0x08],
		SH:   buf[0x09],
		SL:   buf[0x0A],
	}, nil
}
