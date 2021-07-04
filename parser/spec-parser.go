package parser

import (
	"fmt"
	"go/ast"
)

//ParseSpecs get const and var variables
func (p *Parser) parseSpecs(genDecls []*ast.GenDecl) (variables []*ast.GenDecl) {
	specGenDecls := []*ast.GenDecl{}

	for _, genDecl := range genDecls {
		genDeclType := genDecl.Specs[0].(*ast.TypeSpec).Type
		switch genDeclType.(type) {
		case ast.Spec:
			specGenDecls = append(specGenDecls, genDecl)
		default:
			continue
		}
	}
	return specGenDecls
}

//parseValues get const and var values
func (p *Parser) parseValues(specs []ast.Spec) (variables []*ast.ValueSpec) {
	valueSpecs := []*ast.ValueSpec{}

	for _, spec := range specs {
		switch specType := spec.(type) {
		case *ast.ValueSpec:
			valueSpecs = append(valueSpecs, specType)
		default:
			continue
		}
	}
	return valueSpecs
}

//parseStructDeclsByName returns the first gen decl that contains the provided name
// func (p *Parser) parseValueSpecsByNames(genDeclarations []*ast.GenDecl, name ...string) *ast.ValueSpec {
// 	//loop over all general declarations in the file
// 	for _, genDeclaration := range genDeclarations {
// 		genDeclarationSpec := genDeclaration.Specs[0].(*ast.TypeSpec)
// 		switch genDeclarationSpec.Type.(type) {
// 		case *ast.StructType:
// 			if genDeclarationSpec.Name.Name == name {
// 				return genDeclaration
// 			}
// 		default:
// 			continue
// 		}
// 	}
// 	return nil
// }

//ParseImports get 3rd party package imports
func (p *Parser) parseImports(specs []ast.Spec) (variables []*ast.ImportSpec) {
	importSpecs := []*ast.ImportSpec{}

	for _, spec := range specs {
		switch spec.(type) {
		case *ast.ImportSpec:
			importSpec, ok := spec.(*ast.ImportSpec)
			if ok {
				fmt.Println("it's ok")
				importSpecs = append(importSpecs, importSpec)
			}
		default:
			continue
		}
	}
	return importSpecs
}

//ParseTypeSpecs get types inside the type declaration
func (p *Parser) parseTypeSpecs(specs []ast.Spec) (variables []*ast.TypeSpec) {
	typeSpecs := []*ast.TypeSpec{}

	for _, spec := range specs {
		switch spec.(type) {
		case *ast.TypeSpec:
			typeSpec, ok := spec.(*ast.TypeSpec)
			if ok {
				fmt.Println("it's ok")
				typeSpecs = append(typeSpecs, typeSpec)
			}
		default:
			continue
		}
	}
	return typeSpecs
}
