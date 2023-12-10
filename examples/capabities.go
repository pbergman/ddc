package main

import (
	"github.com/pbergman/ddc"
	"log"
	"os"
)

func main() {

	//for i := 0; i < 32; i++ {
	//for i := 6; i < 7; i++ {

	handler, err := ddc.NewWire(6, nil)
	handler.Debug(os.Stdout)

	if err != nil {

		if v, o := err.(*ddc.Error); o && v.Code == ddc.ERROR_DCC_BUS_NOT_FOUND {
			//continue
			panic(err)
		}

		log.Fatal(err)
	}

	defer handler.Close()

	if false == handler.IsActive() {
		//continue
		panic(err)
	}

	ddc.GetCapabilities(handler)
	//}
}
