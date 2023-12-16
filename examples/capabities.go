package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pbergman/ddc"
)

type NopCloser struct {
	io.ReadWriteCloser
}

func (n *NopCloser) Close() error {
	return nil
}

func main() {

	for i := 0; i < 32; i++ {

		handler, err := ddc.NewWire(i, nil)

		if err != nil {

			if v, o := err.(*ddc.Error); o && v.Code == ddc.ERROR_DCC_BUS_NOT_FOUND {
				continue
			}

			log.Fatal(err)
		}

		handler.SetDefaultSleep(50 * time.Millisecond)
		handler.Debug(NopCloser{os.Stdout})

		defer handler.Close()

		if false == handler.IsActive() {
			continue
		}

		fmt.Println(handler.GetCapabilities())

	}
}
