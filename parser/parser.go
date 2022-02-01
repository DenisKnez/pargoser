package parser

import (
	"fmt"
	"go/ast"
)

type IParser interface {
	// Get all packages, and the package contains all the other declarations inside it
	GetPackages() ([]Package, error)
	// Get everything inside the package that matches the provided name
	// GetPackage(packageName string)
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
	GetImports() (types []*Import, err error) //TODO

	//GetTypes   //TODO
	// example:  type Something string
	// parseFiles(directoryName string) ([]*GoFile, error)
}

//Parser used to parse go files
type Parser struct {
	packages []*parserPackage
	// files              []*parserGoFile
}

// parser needs to read one package at a time

func NewParser(fileToParse string) (parsedFiles []string, parser IParser, err error) {
	par := Parser{}
	par.packages, err = parsePackages(".")
	if err != nil {
		return parsedFiles, parser, err
	}

	return parsedFiles, &par, nil
}

func (p *Parser) GetPackages() (packages []Package, err error) {
	for _, parserPkg := range p.packages {
		pkg := Package{}
		pkg.Name = getPackageName(*parserPkg.GoFiles[0])
		pkg.DirectoryPath = parserPkg.DirectoryPath

		for _, parserGoFile := range parserPkg.GoFiles {
			// STRUCT
			structs, err := getStructs(*parserGoFile)
			if err != nil {
				return nil, err
			}

			// VARIABLES
			variables, err := getVariables(*parserGoFile)
			if err != nil {
				return nil, err
			}

			// FUNCTIONS
			functions, err := getFunctions(*parserGoFile)
			if err != nil {
				return nil, err
			}

			// INTERFACES
			interfaces, err := getInterfaces(*parserGoFile)
			if err != nil {
				return nil, err
			}

			// IMPORTS
			imports := getImports(*parserGoFile)

			pkg.Files = append(pkg.Files, GoFile{
				Imports:    imports,
				Structs:    structs,
				Variables:  variables,
				Functions:  functions,
				Interfaces: interfaces,
			})
		}
		packages = append(packages, pkg)

	}

	return packages, nil
}

func (p *Parser) GetImports() (imports []*Import, err error) {
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		imports = append(imports, getImports(goFile)...)
	}

	return imports, nil
}

//GetInterfaces returns all interfaces in the file,
//returns nil if there is not interfaces in the file
func (p *Parser) GetInterfaces() (interfaces []*Interface, err error) {
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		i, err := getInterfaces(goFile)
		if err != nil {
			return nil, err
		}

		interfaces = append(interfaces, i...)
	}

	return interfaces, nil
}

//GetInterface gets the first occurence of interface that has the provided name,
//returns nil if the interface does not exist
func (p *Parser) GetInterface(name string) (theInterface *Interface, err error) {
	interfaces := []*Interface{}
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		genDeclarations := parseGenDeclarations(goFile)
		interfaceDecl := parseInterfaceDeclsByName(name, genDeclarations)
		if interfaceDecl == nil {
			continue
		}
		ids := []*ast.GenDecl{interfaceDecl}
		fileInterfaces, err := convertInterfaceDeclsIntoInterface(ids)
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
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		genDecls := parseGenDeclarations(goFile)
		fmt.Println("amount: ", len(genDecls))
		structDecls := parseStructDeclsByName(structName, genDecls)
		if structDecls == nil {
			continue
		}

		ids := []*ast.GenDecl{structDecls}
		fileStructs, err := convertStructDeclsIntoStruct(goFile, ids)
		if err != nil {
			return nil, err
		}
		structs = append(structs, fileStructs...)
	}

	if len(structs) == 0 {
		return nil, nil
	}

	return structs[0], nil
}

//GetStructs get all the structs
func (p *Parser) GetStructs() (structs []*Struct, err error) {
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		s, err := getStructs(goFile)
		if err != nil {
			return nil, err
		}

		structs = append(structs, s...)
	}
	return structs, nil
}

//GetFunctions get all the functions
func (p *Parser) GetFunctions() (functions []*Function, err error) {
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		funcs, err := getFunctions(goFile)
		if err != nil {
			return nil, err
		}
		functions = append(functions, funcs...)
	}
	return functions, nil
}

// GetFunction gets the first occurance of the function with the provided name
func (p *Parser) GetFunction(funcName string) (theFunc *Function, err error) {
	funcs := []*Function{}
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		funcDecls := parseFuncDeclarations(goFile)
		astFunction := parseFunctionDeclsByName(funcName, funcDecls)
		ids := []*ast.FuncDecl{astFunction}
		fileFunctions, err := convertFunctionDeclsIntoFunction(ids)
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, fileFunctions...)
	}
	return funcs[0], nil
}

//GetVariables gets all the variables
func (p *Parser) GetVariables() (variables []Variable, err error) {
	goFiles := getAllGoFilesFromAllPackages(p.packages)

	for _, goFile := range goFiles {
		vars, err := getVariables(goFile)
		if err != nil {
			return nil, err
		}

		variables = append(variables, vars...)
	}
	return variables, nil
}

//GetConstantVariables gets all the constant variables
func (p *Parser) GetConstantVariables() (consts []Variable, err error) {
	goFiles := getAllGoFilesFromAllPackages(p.packages)
	for _, goFile := range goFiles {
		genDecls := parseGenDeclarations(goFile)
		constVariableGenDecls := parseConstDecls(genDecls)
		fileConsts, err := convertGenDeclsIntoVariable(0, constVariableGenDecls)
		if err != nil {
			return nil, err
		}
		consts = append(consts, fileConsts...)
	}
	return consts, nil
}
