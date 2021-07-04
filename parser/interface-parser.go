package parser

import (
	"go/ast"
)

func (p *Parser) parseInterfaceDeclsByName(name string, genDecls []*ast.GenDecl) *ast.GenDecl {
	//loop over all general declarations in the file
	for _, genDeclaration := range genDecls {
		genDeclarationSpec := genDeclaration.Specs[0]

		switch genDeclarationSpec.(type) {
		case *ast.TypeSpec:
			genDeclarationType := genDeclarationSpec.(*ast.TypeSpec).Type
			switch genDeclarationType.(type) {
			case *ast.InterfaceType:
				if genDeclarationSpec.(*ast.TypeSpec).Name.Name == name {
					return genDeclaration
				}
			default:
				continue
			}
		default:
			continue
		}
	}
	return nil
}

// Takes in general declarations and returns only the general declarations with type interface
func (p *Parser) parseInterfaceDecls(genDeclarations []*ast.GenDecl) (interfaces []*ast.GenDecl) {
	genDeclsWithInterfaceType := []*ast.GenDecl{}
	//loop over all general declarations in the file
	for _, genDeclaration := range genDeclarations {
		genDeclarationSpec := genDeclaration.Specs[0]

		switch genDeclarationSpec.(type) {
		case *ast.TypeSpec:
			genDeclarationType := genDeclarationSpec.(*ast.TypeSpec).Type
			switch genDeclarationType.(type) {
			case *ast.InterfaceType:
				genDeclsWithInterfaceType = append(genDeclsWithInterfaceType, genDeclaration)
			default:
				continue
			}
		default:
			continue
		}

	}
	return genDeclsWithInterfaceType
}

func (p *Parser) parseInterfaceMethodName(astField *ast.Field) string {
	if astField.Names == nil {
		return ""
	}
	return astField.Names[0].Name
}

// takes in ast interfaces and returns this libraries representation of interface
func (p *Parser) convertInterfaceDeclsIntoInterface(genInterfaceDecls []*ast.GenDecl) (
	interfaces []*Interface, err error) {

	for _, genInterfaceDecl := range genInterfaceDecls {
		theInterface := &Interface{}
		genDeclSpec := genInterfaceDecl.Specs[0].(*ast.TypeSpec)
		theInterface.Name = genDeclSpec.Name.Name

		for _, method := range genDeclSpec.Type.(*ast.InterfaceType).Methods.List {
			interfaceMethod := &Method{}

			params, err := p.ParseParameters(method.Type.(*ast.FuncType))
			if err != nil {
				return nil, err
			}
			results, err := p.ParseResults(method.Type.(*ast.FuncType))
			if err != nil {
				return nil, err
			}

			interfaceMethod.Params = append(interfaceMethod.Params, params...)
			interfaceMethod.Results = append(interfaceMethod.Results, results...)
			interfaceMethod.Name = p.parseInterfaceMethodName(method)
			theInterface.Methods = append(theInterface.Methods, interfaceMethod)
		}

		commentGroup, err := p.ParseComments(genInterfaceDecl)
		if err != nil {
			return nil, err
		}

		theInterface.Doc = commentGroup
		interfaces = append(interfaces, theInterface)
	}
	return interfaces, nil
}
