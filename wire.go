package ddc

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/pbergman/logger"
)

func NewWire(bus int, logger *logger.Logger) (*Wire, error) {

	file, err := os.OpenFile("/dev/i2c-"+strconv.Itoa(int(bus)), os.O_RDWR, 0600)

	if err != nil {

		if os.IsNotExist(err) {
			return nil, &Error{Code: ERROR_DCC_BUS_NOT_FOUND, Message: "could not find i2c bus " + strconv.Itoa(int(bus))}
		}

		return nil, &Error{Code: ERROR_DCC_BUS_NOT_OPEN, Message: "could not open i2c bus: " + err.Error()}
	}

	return &Wire{file: file, logger: logger, sleep: time.Millisecond * 10, fd: file.Fd()}, nil
}

type Wire struct {
	file   io.ReadWriteCloser
	fd     uintptr
	sleep  time.Duration
	logger *logger.Logger
	addr   uintptr
}

func (w *Wire) Debug(out io.Writer) {
	w.file = &debug{inner: w.file, outer: out}
}

func (w *Wire) SetDefaultSleep(sleep time.Duration) {
	w.sleep = sleep
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

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, w.fd, cmd, addr)

	if err != 0 {
		return err
	}

	if w.logger != nil {
		w.logger.Debug(fmt.Sprintf("set i2c address to 0x%02x", addr))
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

	n, err := w.file.Write(d)

	if err != nil && nil != w.logger {
		w.logger.Error(fmt.Sprintf("failed to write to i2c, %s", err.Error()))
	}

	if err == nil && nil != w.logger {
		w.logger.Debug(fmt.Sprintf("written %v to i2c", d[:n]))
	}

	return n, err
}

func (w *Wire) Read(d []byte) (int, error) {
	n, err := w.file.Read(d)

	if err != nil && nil != w.logger {
		w.logger.Error(fmt.Sprintf("failed to read from i2c, %s", err.Error()))
	}

	if err == nil && nil != w.logger {
		w.logger.Debug(fmt.Sprintf("read %v from i2c", d[:n]))
	}

	return n, err
}

func (w *Wire) Close() error {
	return w.file.Close()
}

func (w *Wire) GetCapabilities() (string, error) {
	return GetCapabilities(w)
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
