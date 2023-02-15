package main

import (
	"fmt"
	"log"

	"github.com/pbergman/ddc"
)

func main() {

	handler, err := ddc.NewWire(10)

	if err != nil {
		log.Fatal(err)
	}

	defer handler.Close()

	value, err := handler.GetVCP(0x12)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current: %d\n", value.Curr)
	fmt.Printf("Max: 	 %d\n", value.Max)
	fmt.Printf("Index: 	 %#02.x\n", value.Code)
}
