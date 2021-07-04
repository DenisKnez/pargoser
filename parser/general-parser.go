package parser

import (
	"go/ast"
	"go/token"
)

//parseFunctionDecls returns only variable declarations from the provided declarations
func (p *Parser) parseVariableDecls(genDecls []*ast.GenDecl) (valueSpecs []*ast.GenDecl) {
	for _, genDecl := range genDecls {
		switch genDecl.Tok {
		case token.VAR:
			valueSpecs = append(valueSpecs, genDecl)
		}
	}
	return valueSpecs
}

//parseConstDecls returns only const variable declarations from the provided declarations
func (p *Parser) parseConstDecls(genDecls []*ast.GenDecl) (valueSpecs []*ast.GenDecl) {
	for _, genDecl := range genDecls {
		switch genDecl.Tok {
		case token.CONST:
			valueSpecs = append(valueSpecs, genDecl)
		}
	}
	return valueSpecs
}

//singleSpecVariableDeclarationConversion this is when there is just a single
// variable per var keyword
func (p *Parser) singleSpecVariableDeclarationConversion(kind VariableKind, genDecl *ast.GenDecl) (variables []Variable, err error) {
	variable := Variable{}
	commentGroup, err := p.ParseComments(genDecl)
	if err != nil {
		return variables, err
	}
	variable.Doc = commentGroup
	variable.Kind = kind
	variable.Name = genDecl.Specs[0].(*ast.ValueSpec).Names[0].Name
	if len(genDecl.Specs[0].(*ast.ValueSpec).Values) == 0 {
		variable.Value = nil
	} else {
		theType := genDecl.Specs[0].(*ast.ValueSpec).Values[0]
		switch theType.(type) {
		case (*ast.UnaryExpr):
			value := genDecl.Specs[0].(*ast.ValueSpec).Values[0].(*ast.UnaryExpr).Op.String() + genDecl.Specs[0].(*ast.ValueSpec).Values[0].(*ast.UnaryExpr).X.(*ast.Ident).Name
			variable.Value = &value
		case (*ast.BasicLit):
			variable.Value = &genDecl.Specs[0].(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value
		case (*ast.CompositeLit):
			variable.Value = &genDecl.Specs[0].(*ast.ValueSpec).Values[0].(*ast.CompositeLit).Type.(*ast.Ident).Name
		default:
			panic("unsupported type: convertGenDeclsIntoVariable")
		}
	}
	variables = append(variables, variable)
	return variables, nil
}

//multiSpecVariableDeclarationConversion this takes care of the scenario
// where there is multiple variables inside a single var keyword
func (p *Parser) multiSpecVariableDeclarationConversion(kind VariableKind, genDecl *ast.GenDecl) (variables []Variable, err error) {
	for _, spec := range genDecl.Specs {
		variable := Variable{}
		commentGroup, err := parseSpecComments(spec)
		if err != nil {
			return variables, err
		}
		variable.Doc = commentGroup
		variable.Kind = kind
		variable.Name = spec.(*ast.ValueSpec).Names[0].Name
		if len(genDecl.Specs[0].(*ast.ValueSpec).Values) == 0 {
			variable.Value = nil
		} else {
			theType := spec.(*ast.ValueSpec).Values[0]
			switch theType.(type) {
			case (*ast.UnaryExpr):
				theValue := spec.(*ast.ValueSpec).Values[0].(*ast.UnaryExpr).Op.String() + spec.(*ast.ValueSpec).Values[0].(*ast.UnaryExpr).X.(*ast.Ident).Name
				variable.Value = &theValue
			case (*ast.BasicLit):
				variable.Value = &spec.(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value
			case (*ast.CompositeLit):
				variable.Value = &spec.(*ast.ValueSpec).Values[0].(*ast.CompositeLit).Type.(*ast.Ident).Name
			default:
				panic("unsupported type: convertGenDeclsIntoVariable")
			}
		}
		variables = append(variables, variable)
	}
	return variables, nil
}

func (p *Parser) convertGenDeclsIntoVariable(kind VariableKind, genDecls []*ast.GenDecl) (variables []Variable, err error) {
	for _, genDecl := range genDecls {
		if len(genDecl.Specs) == 1 {
			vars, err := p.singleSpecVariableDeclarationConversion(kind, genDecl)
			if err != nil {
				return variables, nil
			}
			variables = append(variables, vars...)
		} else {
			vars, err := p.multiSpecVariableDeclarationConversion(kind, genDecl)
			if err != nil {
				return variables, err
			}
			variables = append(variables, vars...)
		}
	}
	return variables, nil
}

func parseSpecComments(spec ast.Spec) (*CommentGroup, error) {
	astCommentGroup := &ast.CommentGroup{}
	commentGroup := &CommentGroup{}

	switch spec := spec.(type) {
	case *ast.ValueSpec:
		astCommentGroup = spec.Doc
	case *ast.TypeSpec:
		astCommentGroup = spec.Doc
	case *ast.ImportSpec:
		astCommentGroup = spec.Doc
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
