package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/pbergman/ddc"
)

func main() {

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for i := 0; i < 32; i++ {

		handler, err := ddc.NewWire(i)

		if err != nil {

			if v, o := err.(*ddc.Error); o && v.Code == ddc.ERROR_DCC_BUS_NOT_FOUND {
				continue
			}

			log.Fatal(err)
		}

		defer handler.Close()

		info, err := ddc.GetEDID[*ddc.EDID](handler)

		if err == nil {
			fmt.Printf("Found display at bus %d\n", i)
			fmt.Fprintf(writer, "Display Name\t%s\n", info.DisplayName)
			fmt.Fprintf(writer, "Model Serial Number\t%s\n", info.DisplaySerialNumber)
			fmt.Fprintf(writer, "Unspecified Text\t%s\n", info.UnspecifiedText)
			writer.Flush()
			fmt.Printf("\n")

		}
	}
}
