package parser

import (
	"go/ast"
)

//parseFunctionDecls returns only method declarations from the provided declarations
func (p *Parser) parseMethodDecls(funcDecls []*ast.FuncDecl) (funcs []*ast.FuncDecl) {
	for _, funcDecl := range funcDecls {
		if funcDecl.Recv != nil {
			funcs = append(funcs, funcDecl)
		}
	}
	return funcs
}

func (p *Parser) convertFunctionDeclsIntoMethod(funcDecls []*ast.FuncDecl) (methods []*Method, err error) {
	for _, funcDecl := range funcDecls {
		theMethod := &Method{}
		receiver := Receiver{}
		theMethod.Name = funcDecl.Name.Name

		// get parameters
		for _, astParam := range funcDecl.Type.Params.List {
			param := &Parameter{}
			param.Name = astParam.Names[0].Name
			fieldType, err := p.ConvertFieldTypeToString(astParam)
			if err != nil {
				return nil, err
			}
			param.Type = fieldType

			theMethod.Params = append(theMethod.Params, param)
		}

		receiver.Name = funcDecl.Recv.List[0].Names[0].Name

		receiver.PointerReceiver = p.isPointer(funcDecl.Recv.List[0].Type)
		recvType, err := p.ConvertFieldTypeToString(funcDecl.Recv.List[0])
		if err != nil {
			return nil, err
		}
		receiver.Type = recvType
		//get results
		for _, astResult := range funcDecl.Type.Results.List {
			result := &Result{}
			result.Name = astResult.Names[0].Name
			fieldType, err := p.ConvertFieldTypeToString(astResult)
			if err != nil {
				return nil, err
			}
			result.Type = fieldType

			theMethod.Results = append(theMethod.Results, result)
		}
		// get comments
		commentGroup, err := p.ParseComments(funcDecl)
		if err != nil {
			return nil, err
		}
		theMethod.Doc = commentGroup
		methods = append(methods, theMethod)
	}
	return methods, nil
}
