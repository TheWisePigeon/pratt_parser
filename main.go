package main

import (
	"os"
	"pratt_parser/lexer"
)

func main() {
	bytes, err := os.ReadFile("./examples/00.pigeon")
	if err != nil {
		panic(err)
	}
	tokens := lexer.Tokenize(string(bytes))
	for _, token := range tokens {
		token.Debug()
	}
}
