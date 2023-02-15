package ddc

import (
	"os"
	"strconv"
	"syscall"
)

func NewWire(bus int) (*Wire, error) {

	fd, err := os.OpenFile("/dev/i2c-"+strconv.Itoa(int(bus)), os.O_RDWR, 0600)

	if err != nil {

		if os.IsNotExist(err) {
			return nil, &Error{Code: ERROR_DCC_BUS_NOT_FOUND, Message: "could not find i2c bus " + strconv.Itoa(int(bus))}
		}

		return nil, &Error{Code: ERROR_DCC_BUS_NOT_OPEN, Message: "could not open i2c bus: " + err.Error()}
	}

	return &Wire{fd: fd}, nil
}

type Wire struct {
	fd   *os.File
	addr uintptr
}

func (w *Wire) SetAddress(addr uintptr, force bool) error {

	if w.addr == addr && force == false {
		return nil
	}

	var cmd uintptr

	if force {
		cmd = 0x0706 // I2C_SLAVE_FORCE
	} else {
		cmd = 0x0703 // I2C_SLAVE
	}

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, w.fd.Fd(), cmd, addr)

	if err != 0 {
		return err
	}

	w.addr = addr

	return nil
}

func (w *Wire) WriteAt(addr uintptr, d []byte) (int, error) {

	if err := w.SetAddress(addr, false); err != nil {
		return 0, err
	}

	return w.Write(d)
}

func (w *Wire) Write(d []byte) (int, error) {
	return w.fd.Write(d)
}

func (w *Wire) Read(d []byte) (int, error) {
	return w.fd.Read(d)
}

func (w *Wire) Close() error {
	return w.fd.Close()
}

func (w *Wire) GetVCP(index byte) (*VCPResponse, error) {
	return GetVCP(w, index)
}

func (w *Wire) SetVCP(index byte, value uint16) error {
	return SetVCP(w, index, value)
}

func (w *Wire) GetEDID() (*EDID, error) {
	return GetEDID[*EDID](w)
}

func (w *Wire) IsActive() bool {
	return IsActive(w)
}
