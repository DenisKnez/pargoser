package parser

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
)

//isPointer checks if the methods has a pointer receiver
func (p *Parser) isPointer(expression ast.Expr) bool {
	switch expression.(type) {
	case *ast.StarExpr:
		return true
	default:
		return false
	}
}

//ParseFuncDeclarations takes in the file and returns only the general declarations
func (p *Parser) ParseFuncDeclarations(file *ast.File) []*ast.FuncDecl {
	funcDeclarations := []*ast.FuncDecl{}
	//loop over all declaration in the file
	for _, declaration := range file.Decls {
		switch declaration.(type) {
		case *ast.FuncDecl:
			funcDeclarations = append(funcDeclarations, declaration.(*ast.FuncDecl))
		default:
			continue
		}
	}
	return funcDeclarations
}

//ParseGenDeclarations takes in the file and returns only the general declarations
func (p *Parser) ParseGenDeclarations(file *ast.File) []*ast.GenDecl {
	genDeclarations := []*ast.GenDecl{}
	//loop over all declaration in the file
	for _, declaration := range file.Decls {
		switch declaration.(type) {
		case *ast.GenDecl:
			genDeclarations = append(genDeclarations, declaration.(*ast.GenDecl))
		default:
			continue
		}
	}
	return genDeclarations
}

//GetFile get the file to parse
func (p *Parser) GetFile(path string) (*ast.File, error) {
	if p.File != nil {
		return p.File, nil
	}
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}
	p.File = file
	return file, nil
}

//ConvertFieldTypeToString convers the provided field into a string representation
func (p *Parser) ConvertFieldTypeToString(field *ast.Field) (string, error) {
	expression := field.Type
	switch expression.(type) {
	case *ast.Ident:
		return p.IdentityStringConversion(expression), nil
	case *ast.ArrayType:
		return p.ArrayStringConversion(expression), nil
	case *ast.StarExpr:
		return p.StarStringConversion(expression), nil
	case *ast.MapType:
		return p.MapStringConversion(expression), nil
	case *ast.FuncType:
		theFunc, err := p.FuncStringConversion(expression)
		if err != nil {
			return theFunc, err
		}
		return theFunc, nil
	default:
		return "", errors.New("field type not supported")
	}
}

//ParseParameters returns the parameters from the provided function
func (p *Parser) ParseParameters(astFunc *ast.FuncType) (parameters []*Parameter, err error) {
	funcParams := []*Parameter{}
	if astFunc.Params == nil {
		return funcParams, nil
	}
	astParams := astFunc.Params.List

	for _, astParam := range astParams {
		funcParam := &Parameter{}
		//parse parameter names if there are any names assigned
		if astParam.Names != nil {
			funcParam.Name = astParam.Names[0].Name
		}

		funcParam.Type, err = p.ConvertFieldTypeToString(astParam)
		if err != nil {
			return funcParams, err
		}
		funcParam.IsTypePointer = p.isPointer(astParam.Type)
		funcParams = append(funcParams, funcParam)
	}
	return funcParams, nil
}

//ParseResults returns the results from the provided function
func (p *Parser) ParseResults(astFunc *ast.FuncType) (results []*Result, err error) {
	funcResults := []*Result{}
	if astFunc.Results == nil {
		return funcResults, nil
	}
	astResults := astFunc.Results.List

	for _, astResult := range astResults {
		funcResult := &Result{}
		//parse result names if there are any names assigned
		if astResult.Names != nil {
			funcResult.Name = astResult.Names[0].Name
		}

		funcResult.Type, err = p.ConvertFieldTypeToString(astResult)
		if err != nil {
			return funcResults, err
		}
		funcResult.IsTypePointer = p.isPointer(astResult.Type)

		funcResults = append(funcResults, funcResult)
	}
	return funcResults, nil
}

//ParseComments returns the comments from the provided declaration
func (p *Parser) ParseComments(astDecl ast.Decl) (*CommentGroup, error) {
	astCommentGroup := &ast.CommentGroup{}
	commentGroup := &CommentGroup{}

	switch astDecl.(type) {
	case *ast.FuncDecl:
		astCommentGroup = astDecl.(*ast.FuncDecl).Doc
	case *ast.GenDecl:
		astCommentGroup = astDecl.(*ast.GenDecl).Doc

	default:
		return commentGroup, errors.New("declaration type not supported")
	}

	if astCommentGroup == nil {
		commentGroup = &CommentGroup{}
		return commentGroup, nil
	}
	for _, astComment := range astCommentGroup.List {
		comment := Comment{}
		comment.Text = astComment.Text
		commentGroup.Comments = append(commentGroup.Comments, &comment)
	}

	return commentGroup, nil
}
