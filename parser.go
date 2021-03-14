package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "gopher.go", nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, f)
}
