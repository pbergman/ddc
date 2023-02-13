package main

import (
	"fmt"
	"log"

	"github.com/pbergman/ddc"
)

func main() {

	handler, err := ddc.NewDisplayHandler(10)

	if err != nil {
		log.Fatal(err)
	}

	value, err := handler.Get(0x12)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current: %d\n", value.Curr)
	fmt.Printf("Max: 	 %d\n", value.Max)
	fmt.Printf("Index: 	 %#02.x\n", value.Code)
}
