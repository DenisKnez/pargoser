package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Person struct {
	Name     string
	Lastname string
	Age      int
}

//go:generate ../generator/gen
func main() {
	jason := Person{
		Name:     "Jason",
		Lastname: "Pierce",
		Age:      21,
	}

	personTemplate := Template{
		OutputFilePath:   "generated-files",
		TemplateFilePath: "templates",
	}

	err := personTemplate.Generate("Jason.go", "person.tmpl", jason)
	if err != nil {
		fmt.Println(err)
	}

	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, "something.go", nil, parser.DeclarationErrors)
	if err != nil {
		panic(err)
	}

	StartParsing(fileSet, file)

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
	//ast.Print(fileSet, file)
	//fmt.Println("================================================")
	return
}
