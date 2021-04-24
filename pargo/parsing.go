package main

import (
	"fmt"
	"go/ast"
	"go/token"

	//"go/types"
	"reflect"
	//util "golang.org/x/tools/go/types/typeutil"
)

func StartParsing(fileSet *token.FileSet, file *ast.File) {
	parseDeclarations(file)
}

func parseDeclarations(file *ast.File) {

	for _, declaration := range file.Decls {
		if reflect.TypeOf(declaration) != reflect.TypeOf(&ast.GenDecl{}) {
			continue
		}

		genDeclaration, _ := declaration.(*ast.GenDecl)
		parseSpecificDeclaration(genDeclaration)
	}

}

type MethodStructure struct {
	Name       string
	Parameters []string
	Results    []string
}

func parseInterface(theInterface *ast.InterfaceType) {
	interfaceMethods := theInterface.Methods.List

	methodStructures := []MethodStructure{}

	//method
	for _, method := range interfaceMethods {
		methodStructure := MethodStructure{}
		methodStructure.Name = method.Names[0].Name
		methodStructure.Parameters = parseMethodParameters(method)
		methodStructure.Results = parseMethodResults(method)

		methodStructures = append(methodStructures, methodStructure)
	}

	fmt.Print("method structures")
	fmt.Print(methodStructures)
	//createTemplate(methodStructures)

}

//parserMethodParameters parse the parameters from the provided method
func parseMethodParameters(method *ast.Field) (parameters []string) {
	parameterFields := method.Type.(*ast.FuncType).Params.List

	for _, parameterField := range parameterFields {
		switch parameterField.Type.(type) {
		case *ast.ArrayType:
			parameters = append(parameters, parameterField.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name)
		case *ast.Ident:
			parameters = append(parameters, parameterField.Type.(*ast.Ident).Name)
		case *ast.StarExpr:
			parameters = append(parameters, parameterField.Type.(*ast.StarExpr).X.(*ast.Ident).Name)
		case *ast.StructType:
			parameters = append(parameters, parameterField.Type.(*ast.Ident).Name)
		default:
			parameters = append(parameters, parameterField.Type.(*ast.Ident).Name)
		}
	}
	return parameters
}

func parseMethodResults(method *ast.Field) (results []string) {
	resultsFields := method.Type.(*ast.FuncType).Results.List

	for _, resultField := range resultsFields {
		results = append(results, resultField.Type.(*ast.Ident).Name)
	}
	return results
}

func parseSpecificDeclaration(declaration *ast.GenDecl) {
	spec := declaration.Specs[0]

	typeSpec := spec.(*ast.TypeSpec)
	expression := typeSpec.Type

	switch expression.(type) {
	case *ast.InterfaceType:
		interfaceType := expression.(*ast.InterfaceType)
		parseInterface(interfaceType)
	default:
		break
	}

}
