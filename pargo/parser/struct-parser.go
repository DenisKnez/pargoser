package parser

import (
	"go/ast"
)

//parseStructDeclsByName returns the first gen decl that contains the provided name
func (p *Parser) parseStructDeclsByName(name string, genDecls []*ast.GenDecl) *ast.GenDecl {
	//loop over all general declarations in the file
	for _, genDeclaration := range genDecls {
		genDeclarationSpec := genDeclaration.Specs[0].(*ast.TypeSpec)
		switch genDeclarationSpec.Type.(type) {
		case *ast.StructType:
			if genDeclarationSpec.Name.Name == name {
				return genDeclaration
			}
		default:
			continue
		}
	}
	return nil
}

// returns the methods from the provided struct
func (p *Parser) parseStructMethods(receiverName string, funcDecls []*ast.FuncDecl) (funcs []*ast.FuncDecl) {
	funcs = p.parseMethodDeclsByReceiverName(receiverName, funcDecls)
	return funcs
}

// takes in the file and returns only the general declarations with type struct
func (p *Parser) parseStructDecls(genDeclarations []*ast.GenDecl) (structs []*ast.GenDecl) {
	genDeclsWithStructType := []*ast.GenDecl{}
	//loop over all general declarations in the file
	for _, genDeclaration := range genDeclarations {
		genDeclarationType := genDeclaration.Specs[0].(*ast.TypeSpec).Type
		switch genDeclarationType.(type) {
		case *ast.StructType:
			genDeclsWithStructType = append(genDeclsWithStructType, genDeclaration)
		default:
			continue
		}
	}
	return genDeclsWithStructType
}

func (p *Parser) convertStructDeclsIntoStruct(file *ast.File, genStructDecls []*ast.GenDecl) (structs []*Struct, err error) {
	for _, genStructDecl := range genStructDecls {
		theStruct := &Struct{}
		genDeclSpec := genStructDecl.Specs[0].(*ast.TypeSpec)
		theStruct.Name = genDeclSpec.Name.Name
		for _, field := range genDeclSpec.Type.(*ast.StructType).Fields.List {
			structField := &Field{}
			structField.Name = field.Names[0].Name
			fieldTypeString, err := p.ConvertFieldTypeToString(field)
			if err != nil {
				return nil, err
			}
			structField.Type = fieldTypeString
			theStruct.Fields = append(theStruct.Fields, structField)
		}
		commentGroup, err := p.ParseComments(genStructDecl)
		if err != nil {
			return nil, err
		}
		theStruct.CommentGroup = commentGroup
		// get struct methods
		funcDecls := p.ParseFuncDeclarations(file)
		astMethods := p.parseMethodDecls(funcDecls)
		methods, err := p.convertFunctionDeclsIntoMethod(astMethods)
		if err != nil {
			return nil, err
		}
		theStruct.Methods = methods
		structs = append(structs, theStruct)
	}
	return structs, nil
}
