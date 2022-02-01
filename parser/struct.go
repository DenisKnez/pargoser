package parser

import (
	"go/ast"
	"regexp"
)

var tagSplitRegex = regexp.MustCompile(":")
var tagSyntaxRegex = regexp.MustCompile("[\"\x60]")

func getStructs(file parserGoFile) (structs []*Struct, err error) {
	genDecls := parseGenDeclarations(file)
	structDecls := parseStructDecls(genDecls)
	structStructs, err := convertStructDeclsIntoStruct(file, structDecls)
	if err != nil {
		return structs, err
	}
	structs = append(structs, structStructs...)

	return structs, nil
}

//parseStructDeclsByName returns the first gen decl that contains the provided name
func parseStructDeclsByName(name string, genDeclarations []*ast.GenDecl) *ast.GenDecl {
	//loop over all general declarations in the file
	for _, genDeclaration := range genDeclarations {
		genDeclarationTypeSpec := genDeclaration.Specs[0]
		switch typeSpec := genDeclarationTypeSpec.(type) {
		case *ast.TypeSpec:
			switch typeSpec.Type.(type) {
			case *ast.StructType:
				return genDeclaration
			default:
				continue
			}
		}
	}

	return nil
}

// takes in the file and returns only the general declarations with type struct
func parseStructDecls(genDeclarations []*ast.GenDecl) (structs []*ast.GenDecl) {
	genDeclsWithStructType := []*ast.GenDecl{}
	//loop over all general declarations in the file
	for _, genDeclaration := range genDeclarations {
		genDeclarationTypeSpec := genDeclaration.Specs[0]
		switch typeSpec := genDeclarationTypeSpec.(type) {
		case *ast.TypeSpec:
			switch typeSpec.Type.(type) {
			case *ast.StructType:
				genDeclsWithStructType = append(genDeclsWithStructType, genDeclaration)
			default:
				continue
			}
		}

	}
	return genDeclsWithStructType
}

func parseStructsByPackage(file parserGoFile) (structs []*Struct, err error) {
	genDecls := parseGenDeclarations(file)
	structDecls := parseStructDecls(genDecls)
	structStructs, err := convertStructDeclsIntoStruct(file, structDecls)
	if err != nil {
		return nil, err
	}
	structs = append(structs, structStructs...)

	return structs, nil
}

func convertStructDeclsIntoStruct(file parserGoFile, genStructDecls []*ast.GenDecl) (structs []*Struct, err error) {

	for _, genStructDecl := range genStructDecls {
		theStruct := &Struct{}
		genDeclSpec := genStructDecl.Specs[0].(*ast.TypeSpec)
		theStruct.Name = genDeclSpec.Name.Name
		for _, field := range genDeclSpec.Type.(*ast.StructType).Fields.List {
			structField := &Field{}
			if field.Names == nil {
				structField.Name = ""
			} else {
				structField.Name = field.Names[0].Name
			}

			fieldTypeString, err := convertFieldTypeToString(field)
			if err != nil {
				return nil, err
			}
			structField.Type = fieldTypeString

			if field.Tag != nil {
				tagValue := tagSplitRegex.Split(string(tagSyntaxRegex.ReplaceAll([]byte(field.Tag.Value), []byte(""))), -1)

				structField.Tag = &Tag{
					Type: field.Tag.Kind,
					Value: TagValue{
						Type:  tagValue[0],
						Value: tagValue[1],
					},
				}
			}

			theStruct.Fields = append(theStruct.Fields, structField)
		}
		commentGroup, err := parseComments(genStructDecl)
		if err != nil {
			return nil, err
		}
		theStruct.Doc = commentGroup
		// get struct methods
		funcDecls := parseFuncDeclarations(file)
		astMethods := parseMethodDecls(funcDecls)
		methods, err := convertFunctionDeclsIntoMethod(astMethods)
		if err != nil {
			return nil, err
		}
		theStruct.Methods = methods
		structs = append(structs, theStruct)
	}
	return structs, nil
}
