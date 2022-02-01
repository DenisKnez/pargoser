package parser

import (
	"go/ast"
)

func getFunctions(file parserGoFile) (functions []*Function, err error) {
	funcDecls := parseFuncDeclarations(file)
	astFunctions := parseFunctionDecls(funcDecls)
	fileFunctions, err := convertFunctionDeclsIntoFunction(astFunctions)
	if err != nil {
		return nil, err
	}
	functions = append(functions, fileFunctions...)
	return functions, nil
}

//parseFunctionDeclsByName returns the first gen decl that contains the provided name
//if none is found returns nil
func parseFunctionDeclsByName(name string, funcDecls []*ast.FuncDecl) *ast.FuncDecl {
	for _, funcDecl := range funcDecls {
		if funcDecl.Recv == nil && funcDecl.Name.Name == name {
			return funcDecl
		}
	}
	return nil
}

//parseFunctionDecls returns only function declarations from the provided declarations
func parseFunctionDecls(funcDecls []*ast.FuncDecl) (funcs []*ast.FuncDecl) {
	for _, funcDecl := range funcDecls {
		if funcDecl.Recv == nil {
			funcs = append(funcs, funcDecl)
		}
	}
	return funcs
}

func convertFunctionDeclsIntoFunction(funcDecls []*ast.FuncDecl) (functions []*Function, err error) {
	for _, funcDecl := range funcDecls {
		theFunc := &Function{}
		theFunc.Name = funcDecl.Name.Name

		results, err := parseResults(funcDecl.Type)
		if err != nil {
			return nil, err
		}
		theFunc.Results = results
		params, err := parseParameters(funcDecl.Type)
		if err != nil {
			return nil, err
		}
		theFunc.Params = params

		// get comments
		commentGroup, err := parseComments(funcDecl)
		if err != nil {
			return nil, err
		}
		theFunc.Doc = commentGroup
		functions = append(functions, theFunc)
	}
	return functions, nil
}
