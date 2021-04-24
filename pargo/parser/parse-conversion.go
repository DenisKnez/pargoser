package parser

import (
	"fmt"
	"go/ast"
)

//FuncStringConversion converts a function expression into a string representation of the function
func (p *Parser) FuncStringConversion(expression ast.Expr) (theFunc string, err error) {
	funcType := expression.(*ast.FuncType)
	params, err := p.ParseParameters(funcType)
	if err != nil {
		return "", err
	}
	results, err := p.ParseResults(funcType)
	if err != nil {
		return "", err
	}
	//TODO improve this to be more readable
	ExprValues := ""
	ExprValues += "func("
	for _, param := range params {
		ExprValues += fmt.Sprintf("%s", param)
	}
	ExprValues += fmt.Sprint(") ")
	ExprValues += fmt.Sprint("(")
	for _, result := range results {
		ExprValues += fmt.Sprintf("%s", result)
	}
	ExprValues += fmt.Sprint(")")
	return ExprValues, nil
}

//MapStringConversion converts a map expression into a string representation of the map
func (p *Parser) MapStringConversion(expression ast.Expr) string {
	key := expression.(*ast.MapType).Key.(*ast.Ident).Name
	value := expression.(*ast.MapType).Value.(*ast.Ident).Name
	return fmt.Sprintf("map[%s]%s", key, value)
}

//StarStringConversion converts a star(aka pointer) expression into a string representation of the star
func (p *Parser) StarStringConversion(expression ast.Expr) string {
	name := expression.(*ast.StarExpr).X.(*ast.Ident).Name
	return fmt.Sprintf("*%s", name)
}

//ArrayStringConversion converts an array expression into a string representation of the array
func (p *Parser) ArrayStringConversion(expression ast.Expr) string {
	name := expression.(*ast.ArrayType).Elt.(*ast.Ident).Name
	return fmt.Sprintf("[]%s", name)
}

//IdentityStringConversion converts an identity expression into a string representation of the identity
func (p *Parser) IdentityStringConversion(expression ast.Expr) string {
	return expression.(*ast.Ident).Name
}
