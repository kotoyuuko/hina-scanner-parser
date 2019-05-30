package main

import (
	"fmt"
	"os"

	"github.com/kotoyuuko/hina-scanner-parser/parser"
	"github.com/kotoyuuko/hina-scanner-parser/scanner"
)

func main() {
	file, err := os.Open("test.hina")
	if err != nil {
		panic(err)
	}
	s := scanner.NewScanner(file)
	tokens := s.Scan()

	fmt.Println("Tokens:")
	for key, item := range tokens {
		fmt.Printf("  [%3d]( %9s , `%s` )\n", key, item.Token.String(), item.Code)
	}
	fmt.Println()

	fmt.Print("Parser: ")
	parser.Init(tokens)
	err = parser.Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}
}
