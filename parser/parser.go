package parser

import (
	"go/ast"
)

type IParser interface {
	//Interface
	GetInterfaces() (interfaces []*Interface, err error)
	GetInterface(name string) (theInterface *Interface, err error)
	//Struct
	GetStruct(structName string) (theStruct *Struct, err error)
	GetStructs() (structs []*Struct, err error)
	//Function
	GetFunction(funcName string) (theFunc *Function, err error)
	GetFunctions() (funcs []*Function, err error)
	//Variables
	//GetVariable(variableName string) (theVariable *Variable, err error) //TODO
	//GetVariablesByNames(...string) (variables []*Variable, err error)   //TODO
	GetVariables() (vars []Variable, err error)
	GetConstantVariables() (consts []Variable, err error)
	//Import
	//GetImports() (types []*TypeSpec, err error)                 //TODO
	//GetTypes   //TODO
	ParseFiles(directoryName string) ([]string, error)
}

//Parser used to parse go files
type Parser struct {
	astFiles []*ast.File
}

func NewParser(fileToParse string) (parsedFiles []string, parser IParser, err error) {
	par := Parser{}
	parsedFiles, err = par.ParseFiles(fileToParse)
	if err != nil {
		return []string{}, nil, err
	}
	return parsedFiles, &par, nil
}

//GetInterfaces returns all interfaces in the file,
//returns nil if there is not interfaces in the file
func (p *Parser) GetInterfaces() (interfaces []*Interface, err error) {
	for _, astFile := range p.astFiles {
		genDeclarations := p.ParseGenDeclarations(astFile)
		genInterfaceDeclarations := p.parseInterfaceDecls(genDeclarations)
		fileInterfaces, err := p.convertInterfaceDeclsIntoInterface(genInterfaceDeclarations)
		if err != nil {
			return nil, err
		}
		interfaces = append(interfaces, fileInterfaces...)
	}

	return interfaces, nil
}

//GetInterface gets the first occurence of interface that has the provided name,
//returns nil if the interface does not exist
func (p *Parser) GetInterface(name string) (theInterface *Interface, err error) {
	interfaces := []*Interface{}
	for _, astFile := range p.astFiles {
		genDeclarations := p.ParseGenDeclarations(astFile)
		interfaceDecl := p.parseInterfaceDeclsByName(name, genDeclarations)
		if interfaceDecl == nil {
			return nil, nil
		}
		ids := []*ast.GenDecl{interfaceDecl}
		fileInterfaces, err := p.convertInterfaceDeclsIntoInterface(ids)
		if err != nil {
			return nil, err
		}
		interfaces = append(interfaces, fileInterfaces...)
	}

	if len(interfaces) == 0 {
		return nil, nil
	}
	return interfaces[0], nil
}

//GetStruct gets the first struct that has the provided name, if no struct is found returns nil
func (p *Parser) GetStruct(structName string) (theStruct *Struct, err error) {
	structs := []*Struct{}
	for _, astFile := range p.astFiles {
		genDecls := p.ParseGenDeclarations(astFile)
		structDecls := p.parseStructDeclsByName(structName, genDecls)
		if structDecls == nil {
			return nil, nil
		}
		ids := []*ast.GenDecl{structDecls}
		fileStructs, err := p.convertStructDeclsIntoStruct(astFile, ids)
		if err != nil {
			return nil, err
		}
		structs = append(structs, fileStructs...)
	}

	return structs[0], nil
}

//GetStructs get all the structs
func (p *Parser) GetStructs() (structs []*Struct, err error) {
	for _, astFile := range p.astFiles {
		genDecls := p.ParseGenDeclarations(astFile)
		structDecls := p.parseStructDecls(genDecls)
		structStructs, err := p.convertStructDeclsIntoStruct(astFile, structDecls)
		if err != nil {
			return nil, err
		}
		structs = append(structs, structStructs...)
	}
	return structs, nil
}

//GetFunctions get all the functions
func (p *Parser) GetFunctions() (functions []*Function, err error) {
	for _, astFile := range p.astFiles {
		funcDecls := p.ParseFuncDeclarations(astFile)
		astFunctions := p.parseFunctionDecls(funcDecls)
		fileFunctions, err := p.convertFunctionDeclsIntoFunction(astFunctions)
		if err != nil {
			return nil, err
		}
		functions = append(functions, fileFunctions...)
	}
	return functions, nil
}

//GetFunction get the function with the provided name
func (p *Parser) GetFunction(funcName string) (theFunc *Function, err error) {
	funcs := []*Function{}
	for _, astFile := range p.astFiles {
		funcDecls := p.ParseFuncDeclarations(astFile)
		astFunction := p.parseFunctionDeclsByName(funcName, funcDecls)
		ids := []*ast.FuncDecl{astFunction}
		fileFunctions, err := p.convertFunctionDeclsIntoFunction(ids)
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, fileFunctions...)
	}
	return funcs[0], nil
}

//GetVariables gets all the variables
func (p *Parser) GetVariables() (variables []Variable, err error) {
	for _, astFile := range p.astFiles {
		genDecls := p.ParseGenDeclarations(astFile)
		variableGenDecls := p.parseVariableDecls(genDecls)
		theVars, err := p.convertGenDeclsIntoVariable(1, variableGenDecls)
		if err != nil {
			return nil, err
		}
		variables = append(variables, theVars...)
	}
	return variables, nil
}

//GetConstantVariables gets all the constant variables
func (p *Parser) GetConstantVariables() (consts []Variable, err error) {
	for _, astFile := range p.astFiles {
		genDecls := p.ParseGenDeclarations(astFile)
		constVariableGenDecls := p.parseConstDecls(genDecls)
		fileConsts, err := p.convertGenDeclsIntoVariable(0, constVariableGenDecls)
		if err != nil {
			return nil, err
		}
		consts = append(consts, fileConsts...)
	}
	return consts, nil
}
