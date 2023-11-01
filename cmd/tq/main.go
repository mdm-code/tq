package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mdm-code/scanner"
)

func main() {
	s, err := scanner.New(os.Stdin)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	for s.Scan() {
		t := s.Token()
		fmt.Printf("%v\n", t)
	}
}
