package parser

import (
	"fmt"
	"go/ast"
)

//FuncStringConversion converts a function expression into a string representation of the function
func funcStringConversion(expression ast.Expr) (theFunc string, err error) {
	funcType := expression.(*ast.FuncType)
	params, err := parseParameters(funcType)
	if err != nil {
		return "", err
	}
	results, err := parseResults(funcType)
	if err != nil {
		return "", err
	}

	//TODO improve this to be more readable
	ExprValues := ""
	ExprValues += "func("
	for _, param := range params {
		ExprValues += fmt.Sprintf("%v", param)
	}
	ExprValues += ") "
	ExprValues += "("
	for _, result := range results {
		ExprValues += fmt.Sprintf("%v", result)
	}
	ExprValues += ")"

	return ExprValues, nil
}

//MapStringConversion converts a map expression into a string representation of the map
func mapStringConversion(expression ast.Expr) string {
	key := expression.(*ast.MapType).Key.(*ast.Ident).Name
	value := expression.(*ast.MapType).Value.(*ast.Ident).Name
	return fmt.Sprintf("map[%s]%s", key, value)
}

//StarStringConversion converts a star(aka pointer) expression into a string representation of the star
func starStringConversion(expression ast.Expr) string {
	name := expression.(*ast.StarExpr).X.(*ast.Ident).Name
	return fmt.Sprintf("*%s", name)
}

//ArrayStringConversion converts an array expression into a string representation of the array
func arrayStringConversion(expression ast.Expr) string {
	var formatedName string
	arrayType := expression.(*ast.ArrayType).Elt
	switch arrayType := arrayType.(type) {
	case *ast.Ident:
		name := arrayType.Name
		formatedName = fmt.Sprintf("[]%s", name)
	case *ast.StarExpr:
		name := arrayType.X.(*ast.Ident).Name
		formatedName = fmt.Sprintf("[]*%s", name)
	}
	return formatedName
}

//IdentityStringConversion converts an identity expression into a string representation of the identity
func identityStringConversion(expression ast.Expr) string {
	return expression.(*ast.Ident).Name
}
