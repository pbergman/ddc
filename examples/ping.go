package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pbergman/ddc"
)

func main() {

	handler, err := ddc.NewDisplayHandler(10)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 50; i++ {
		if handler.IsCLosed() {
			fmt.Println("Closed")
		} else {
			fmt.Println("Open")
		}

		time.Sleep(5 * time.Second)
	}
}
