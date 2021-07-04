package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	dp "github.com/DenisKnez/pargoser/parser"
)

type Person struct {
	Name     string
	Lastname string
	Age      int
}

//go:generate ../generator/gen
func main() {

	par, err := dp.NewParser(".")
	if err != nil {
		fmt.Println("the error: ", err)
	}

	fun, err := par.GetInterface("Login")
	if err != nil {
		panic(err)
	}

	fmt.Println(fun)
	// par := dp.NewParser("something.go")
	// variables, err := par.GetVariables()
	// if err != nil {
	// 	panic(err)
	// }

	// for _, vvv := range variables {

	// 	buffer := &bytes.Buffer{}
	// 	encoder := json.NewEncoder(buffer)
	// 	encoder.SetEscapeHTML(false)
	// 	encoder.SetIndent("", "  ")
	// 	encoder.Encode(vvv)
	// 	fmt.Println(buffer.String())
	// }

	//ParsingFileStuff()

	fmt.Scan()
}

//ParsingFileStuff parsing stuff
func ParsingFileStuff() (fileSet *token.FileSet, file *ast.File) {
	fileSet = token.NewFileSet()

	file, err := parser.ParseFile(fileSet, "something.go", nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		panic(err)
	}
	//show the ast tree in the terminal
	ast.Print(fileSet, file)
	//fmt.Println("================================================")
	return
}
