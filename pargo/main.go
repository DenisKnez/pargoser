package main

import (
	"bytes"
	"encoding/json"
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

	par := dp.NewParser("something.go")
	variables, err := par.GetConstantVariables()
	if err != nil {
		panic(err)
	}

	for _, vvv := range variables {
		// bits, err := json.MarshalIndent(vvv, "", "    ")
		// if err != nil {
		// 	panic(err)
		// }

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		encoder.Encode(vvv)
		fmt.Println(buffer.String())
	}

	//ParsingFileStuff()

	fmt.Scan()
}

// func getStructs(theStructs []*pargoser.Struct) {
// 	for _, theStruct := range theStructs {
// 		fmt.Printf("Struct name: %s\n", theStruct.Name)
// 		for _, field := range theStruct.Fields {
// 			fmt.Printf("\tField name: %s\n", field.Name)
// 			fmt.Printf("\tField type: %s\n", field.Type)
// 			fmt.Println()
// 		}
// 	}
// }

// func getInterfaces(theInterfaces []*pargoser.Interface) {
// 	for _, theInt := range theInterfaces {
// 		fmt.Printf("Interface name: %s\n", theInt.Name)
// 		for _, method := range theInt.Methods {
// 			fmt.Printf("\tMethod name: %s\n", method.Name)
// 			fmt.Printf("\t\tParams: ")
// 			fmt.Printf(" ")
// 			for _, param := range method.Params {
// 				fmt.Printf("%s, ", param)
// 			}
// 			fmt.Println()
// 			fmt.Printf("\t\tResults: ")
// 			fmt.Printf(" ")
// 			for _, result := range method.Results {
// 				fmt.Printf("%s, ", result)
// 			}
// 			fmt.Println()
// 			fmt.Println()
// 			fmt.Println()
// 		}
// 	}
// }

// func getSingleInterface(theInt *pargoser.Interface) {
// 	fmt.Printf("Interface name: %s\n", theInt.Name)
// 	for _, method := range theInt.Methods {
// 		fmt.Printf("\tMethod name: %s\n", method.Name)
// 		fmt.Printf("\t\tParams: ")
// 		fmt.Printf(" ")
// 		for _, param := range method.Params {
// 			fmt.Printf("%s, ", param)
// 		}
// 		fmt.Println()
// 		fmt.Printf("\t\tResults: ")
// 		fmt.Printf(" ")
// 		for _, result := range method.Results {
// 			fmt.Printf("%s, ", result)
// 		}
// 		fmt.Println()
// 		fmt.Println()
// 		fmt.Println()
// 	}
// }

// func Setup() {
// 	jason := Person{
// 		Name:     "Jason",
// 		Lastname: "Pierce",
// 		Age:      21,
// 	}

// 	personTemplate := Template{
// 		OutputFilePath:   "generated-files",
// 		TemplateFilePath: "templates",
// 	}

// 	err := personTemplate.Generate("Jason.go", "person.tmpl", jason)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fileSet := token.NewFileSet()

// 	file, err := parser.ParseFile(fileSet, "something.go", nil, parser.DeclarationErrors)
// 	if err != nil {
// 		panic(err)
// 	}

// 	StartParsing(fileSet, file)
// }

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
