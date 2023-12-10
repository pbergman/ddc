package ddc

import (
	"fmt"
	"time"
)

const (
	CAPABILITIE_ADDR uintptr = 0xF3
)

// https://glenwing.github.io/docs/VESA-DDCCI-1.1.pdf
// 4.6 Capabilities Request & Capabilities Reply

func capabilitiesRequestPacket(offset int) []byte {
	return ddcCreatePackage([]byte{
		byte(CAPABILITIE_ADDR),
		byte(offset >> 8),
		byte(offset),
	})
}

func readCapabilitiesResponse(wire *Wire) (string, int, error) {

	var offset = 6 //  addr + lenght + (2) offset + reply code + chks
	var buf = make([]byte, 32+offset)

	x, err := wire.Read(buf)

	if err != nil {
		return "", 0, err
	}

	buf = buf[:x]

	if buf[0] != byte(ADDR_DEST) {
		return "", 0, &Error{Code: ERROR_DCC_RESP_INVALID_ADDR, Message: fmt.Sprintf("invalid address byte in response, expected 0x6e, actual %#02x", buf[0])}
	}

	if buf[2] != 0xe3 {
		return "", 0, &Error{Code: ERROR_DCC_RESP_INVALID_FEATURE, Message: fmt.Sprintf("invalid capabilities reply op code, expected 0xe3, actual %#02x", buf[3])}
	}

	var length = int(buf[1]&0x7f) + 3
	var mem = make([]byte, length+1)

	mem[0] = 0x6f
	mem[1] = 0x6e

	copy(mem[2:], buf[1:length])

	var chks byte = 0x50

	xor(mem[1:length], &chks)

	if chks != mem[length] {
		return "", 0, fmt.Errorf("invalid checksum, expecting %#02x got %#02x", chks, mem[length])
	}

	if length == offset {
		return "", 0, nil
	}

	return string(buf[offset-1 : length-1]), length - offset, nil
}

func GetCapabilities(wire *Wire) (string, error) {

	var offset = 0
	var capabilities string

	for {

		_, err := wire.WriteAt(DDC_ADDR, capabilitiesRequestPacket(offset))

		if err != nil {
			return "", err
		}

		time.Sleep(wire.sleep)

		resp, size, err := readCapabilitiesResponse(wire)

		if err != nil {
			panic(err)
		}

		offset += size

		if resp == "" {
			break
		}

		capabilities += resp
	}

	return capabilities, nil
}
