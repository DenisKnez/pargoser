package parser

import (
	"go/ast"
)

//parseFunctionDeclsByName returns the first gen decl that contains the provided name
//if none is found returns nil
func (p *Parser) parseFunctionDeclsByName(name string, funcDecls []*ast.FuncDecl) *ast.FuncDecl {
	for _, funcDecl := range funcDecls {
		if funcDecl.Recv == nil && funcDecl.Name.Name == name {
			return funcDecl
		}
	}
	return nil
}

//parseFunctionDecls returns only function declarations from the provided declarations
func (p *Parser) parseFunctionDecls(funcDecls []*ast.FuncDecl) (funcs []*ast.FuncDecl) {
	for _, funcDecl := range funcDecls {
		if funcDecl.Recv == nil {
			funcs = append(funcs, funcDecl)
		}
	}
	return funcs
}

func (p *Parser) convertFunctionDeclsIntoFunction(funcDecls []*ast.FuncDecl) (functions []*Function, err error) {
	for _, funcDecl := range funcDecls {
		theFunc := &Function{}
		theFunc.Name = funcDecl.Name.Name

		results, err := p.ParseResults(funcDecl.Type)
		if err != nil {
			return nil, err
		}
		theFunc.Results = results
		params, err := p.ParseParameters(funcDecl.Type)
		if err != nil {
			return nil, err
		}
		theFunc.Params = params

		// get comments
		commentGroup, err := p.ParseComments(funcDecl)
		if err != nil {
			return nil, err
		}
		theFunc.Doc = commentGroup
		functions = append(functions, theFunc)
	}
	return functions, nil
}
