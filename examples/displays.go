package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"github.com/pbergman/ddc"
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

		handler, err := ddc.NewWire(i, nil)

		if err != nil {

			if v, o := err.(*ddc.Error); o && v.Code == ddc.ERROR_DCC_BUS_NOT_FOUND {
				continue
			}

			log.Fatal(err)
		}

		handler.Debug(NopCloser{os.Stdout})

		defer handler.Close()

		if false == handler.IsActive() {
			continue
		}

		info, err := ddc.GetEDID[*ddc.EDID](handler)

		if err == nil {
			fmt.Printf("Found display at bus %d\n", i)
			fmt.Fprintf(writer, "Display Name\t%s\n", info.DisplayName)
			fmt.Fprintf(writer, "Model Serial Number\t%s\n", info.DisplaySerialNumber)
			fmt.Fprintf(writer, "Manufacture Year\t%d\n", info.YearOfManufacture)
			fmt.Fprintf(writer, "Manufacture Week\t%d\n", info.WeekOfManufacture)
			fmt.Fprintf(writer, "Version\t%d.%d\n", info.Version, info.Revision)
			fmt.Fprintf(writer, "Serial Number\t%d\n", info.SerialNumber)
			fmt.Fprintf(writer, "Product Code\t%d\n", info.ManufactureProductCode)
			fmt.Fprintf(writer, "Unspecified Text\t%s\n", info.UnspecifiedText)
			writer.Flush()
			fmt.Printf("\n")
		}
	}
}
