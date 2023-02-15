package main

import (
	"fmt"
	"log"

	"github.com/pbergman/ddc"
)

func main() {

	for i := 0; i < 32; i++ {

		handler, err := ddc.NewWire(i)

		if err != nil {

			if v, o := err.(*ddc.Error); o && v.Code == ddc.ERROR_DCC_BUS_NOT_FOUND {
				continue
			}

			log.Fatal(err)
		}

		defer handler.Close()

		if handler.IsActive() {
			fmt.Printf("Found display with quick check on bus %d\n", i)
		}
	}
}
