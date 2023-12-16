package main

import (
	"fmt"
	"github.com/pbergman/ddc"
	"io"
	"log"
	"os"
	"text/tabwriter"
)

type NopCloser struct {
	io.ReadWriteCloser
}

func (n *NopCloser) Close() error {
	return nil
}
func main() {

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for i := 0; i < 32; i++ {

		if err := getValue(i, writer); err != nil {
			log.Fatalln(err)
		}
	}

}

func getValue(bus int, writer *tabwriter.Writer) error {
	handler, err := ddc.NewWire(bus, nil)

	if err != nil {
		if v, o := err.(*ddc.Error); o && v.Code == ddc.ERROR_DCC_BUS_NOT_FOUND {
			return nil
		}

		return err
	}

	defer handler.Close()

	handler.Debug(NopCloser{os.Stdout})

	var active = false

	for i := 0; i < 5; i++ {
		if handler.IsActive() {
			active = true
			break
		}
	}

	if !active {
		return nil
	}

	fmt.Printf("Found display at bus %d\n", bus)

	if value, err := handler.GetVCP(0xDF); err == nil {
		fmt.Fprintf(writer, "VCP Version: %d.%d\n", value.SH, value.SL)
	}

	if value, err := handler.GetVCP(0x10); err == nil {
		fmt.Printf("Found display at bus %d\n", bus)
		fmt.Fprintf(writer, "Current: %d\n", value.GetCurr())
		fmt.Fprintf(writer, "Max: 	 %d\n", value.GetMax())
		fmt.Fprintf(writer, "Index: 	 %#02.x\n", value.Code)
	}

	writer.Flush()

	return nil
}
