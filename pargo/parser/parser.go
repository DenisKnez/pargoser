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
}

//Parser used to parse go files
type Parser struct {
	fileToParse string
	*ast.File
}

func NewParser(fileToParse string) IParser {
	return &Parser{fileToParse: fileToParse}
}

//GetInterfaces returns all interfaces in the file,
//returns nil if there is not interfaces in the file
func (p *Parser) GetInterfaces() (interfaces []*Interface, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	genDeclarations := p.ParseGenDeclarations(file)
	genInterfaceDeclarations := p.parseInterfaceDecls(genDeclarations)
	interfaces, err = p.convertInterfaceDeclsIntoInterface(genInterfaceDeclarations)
	if err != nil {
		return nil, err
	}
	return interfaces, nil
}

//GetInterface gets the first occurence of interface that has the provided name,
//returns nil if the interface does not exist
func (p *Parser) GetInterface(name string) (theInterface *Interface, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	genDeclarations := p.ParseGenDeclarations(file)
	interfaceDecl := p.parseInterfaceDeclsByName(name, genDeclarations)
	if interfaceDecl == nil {
		return nil, nil
	}
	ids := []*ast.GenDecl{interfaceDecl}
	interfaces := []*Interface{}
	interfaces, err = p.convertInterfaceDeclsIntoInterface(ids)
	if err != nil {
		return nil, err
	}
	return interfaces[0], nil
}

//GetStruct gets the first struct that has the provided name, if no struct is found returns nil
func (p *Parser) GetStruct(structName string) (theStruct *Struct, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	genDecls := p.ParseGenDeclarations(file)
	structDecls := p.parseStructDeclsByName(structName, genDecls)
	if structDecls == nil {
		return nil, nil
	}
	ids := []*ast.GenDecl{structDecls}
	structs, err := p.convertStructDeclsIntoStruct(file, ids)
	if err != nil {
		return nil, err
	}
	return structs[0], nil
}

//GetStructs get all the structs
func (p *Parser) GetStructs() (structs []*Struct, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	genDecls := p.ParseGenDeclarations(file)
	structDecls := p.parseStructDecls(genDecls)
	structStructs, err := p.convertStructDeclsIntoStruct(file, structDecls)
	if err != nil {
		return nil, err
	}
	return structStructs, nil
}

//GetFunctions get all the functions
func (p *Parser) GetFunctions() (functions []*Function, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	funcDecls := p.ParseFuncDeclarations(file)
	astFunctions := p.parseFunctionDecls(funcDecls)
	functions, err = p.convertFunctionDeclsIntoFunction(astFunctions)
	if err != nil {
		return nil, err
	}
	return functions, nil
}

//GetFunction get the function with the provided name
func (p *Parser) GetFunction(funcName string) (theFunc *Function, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	funcDecls := p.ParseFuncDeclarations(file)
	astFunction := p.parseFunctionDeclsByName(funcName, funcDecls)
	ids := []*ast.FuncDecl{astFunction}
	functions, err := p.convertFunctionDeclsIntoFunction(ids)
	if err != nil {
		return nil, err
	}
	return functions[0], nil
}

//GetVariables gets all the variables
func (p *Parser) GetVariables() (theVar []Variable, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	genDecls := p.ParseGenDeclarations(file)
	variableGenDecls := p.parseVariableDecls(genDecls)
	theVar, err = p.convertGenDeclsIntoVariable(1, variableGenDecls)
	if err != nil {
		return nil, err
	}
	return theVar, nil
}

//GetConstantVariables gets all the constant variables
func (p *Parser) GetConstantVariables() (consts []Variable, err error) {
	file, err := p.GetFile(p.fileToParse)
	if err != nil {
		return nil, err
	}
	genDecls := p.ParseGenDeclarations(file)
	constVariableGenDecls := p.parseConstDecls(genDecls)
	consts, err = p.convertGenDeclsIntoVariable(0, constVariableGenDecls)
	if err != nil {
		return nil, err
	}
	return consts, nil
}
