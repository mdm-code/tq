package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/lexer"
)

func main() {
	s, err := scanner.New(os.Stdin)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	l, err := lexer.New(s)
	fmt.Println(l.Buffer)
}
