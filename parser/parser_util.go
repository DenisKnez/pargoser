package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"regexp"
)

var goFilesRegex = regexp.MustCompile(".*.go")

//isPointer checks if the methods has a pointer receiver
func isPointer(expression ast.Expr) bool {
	switch expression.(type) {
	case *ast.StarExpr:
		return true
	default:
		return false
	}
}

//ParseFuncDeclarations takes in the file and returns only the general declarations
func parseFuncDeclarations(file parserGoFile) []*ast.FuncDecl {
	funcDeclarations := []*ast.FuncDecl{}
	//loop over all declaration in the file
	for _, declaration := range file.AstFile.Decls {
		switch declType := declaration.(type) {
		case *ast.FuncDecl:
			funcDeclarations = append(funcDeclarations, declType)
		default:
			continue
		}
	}
	return funcDeclarations
}

//ParseGenDeclarations takes in the file and returns only the general declarations
func parseGenDeclarations(file parserGoFile) []*ast.GenDecl {
	genDeclarations := []*ast.GenDecl{}
	//loop over all declaration in the file
	for _, declaration := range file.AstFile.Decls {
		switch declType := declaration.(type) {
		case *ast.GenDecl:
			genDeclarations = append(genDeclarations, declType)
		default:
			continue
		}
	}
	return genDeclarations
}

func parsePackages(directoryName string) (packageDirectories []*parserPackage, err error) {

	directoryFiles, err := os.ReadDir(directoryName)
	if err != nil {
		return packageDirectories, err
	}

	// initialize the root package
	packageDirectory := &parserPackage{}
	packageDirectory.DirectoryPath = directoryName
	packageDirectory.InnerDirectories = make([]parserDirectory, 0)
	packageDirectory.GoFiles = make([]*parserGoFile, 0)

	for _, file := range directoryFiles {
		// if it's a directory append it to inner directorise of the package
		// if it's a go file append it to go files of the package
		if file.IsDir() {
			packageDirectory.InnerDirectories = append(packageDirectory.InnerDirectories, parserDirectory{
				DirectoryPath: fmt.Sprintf("%s/%s", directoryName, file.Name()),
				FsDirectory:   file,
			})
		} else {
			// if it's not a go file skip it, checks the extension .go to determine if it's a go file
			if !goFilesRegex.MatchString(file.Name()) {
				continue
			}

			goFile, err := parseFile(file, file.Name())
			if err != nil {
				fmt.Println("the parse go file error")
				fmt.Println(err)
				return packageDirectories, err
			}

			packageDirectory.GoFiles = append(packageDirectory.GoFiles, goFile)
		}
	}

	packageDirectories = append(packageDirectories, packageDirectory)

	err = parsePackage(&packageDirectories, packageDirectory)
	if err != nil {
		return packageDirectories, err
	}

	return packageDirectories, nil
}

func parsePackage(allPackageDirectories *[]*parserPackage, parserPackages ...*parserPackage) (err error) {
	// fmt.Println("print once: ", len(packageDirectories))
	packageDirectories := make([]*parserPackage, 0)

	if len(parserPackages) == 0 {
		return
	}

	for _, parserPkg := range parserPackages {

		// fmt.Println("inner: ", len(packageDirectory.InnerDirectories))
		for _, innerDirectory := range parserPkg.InnerDirectories {
			// fmt.Println("the inner directory: ", innerDirectory.DirectoryPath)
			// fmt.Println("reading: ", innerDirectory.DirectoryPath)
			files, err := os.ReadDir(innerDirectory.DirectoryPath)
			if err != nil {
				return err
			}

			pp := new(parserPackage)
			pp.DirectoryPath = innerDirectory.DirectoryPath

			for _, file := range files {
				// fmt.Println("the ff filej: ", file.Name())
				if file.IsDir() {
					pp.InnerDirectories = append(pp.InnerDirectories, parserDirectory{
						DirectoryPath: fmt.Sprintf("%s/%s", innerDirectory.DirectoryPath, file.Name()),
						FsDirectory:   file,
					})
					continue
				}

				// if it's not a go file skip it, checks the extension .go to determine if it's a go file
				if !goFilesRegex.MatchString(file.Name()) {
					continue
				}

				goFile, err := parseFile(file, fmt.Sprintf("%s/%s", innerDirectory.DirectoryPath, file.Name()))
				if err != nil {
					fmt.Println("parsing go stuff: ", err)
					return err
				}

				pp.GoFiles = append(pp.GoFiles, goFile)
			}

			// package is only a package if it has .go files
			if len(pp.GoFiles) != 0 {
				packageDirectories = append(packageDirectories, pp)
			}

		}
	}

	*allPackageDirectories = append(*allPackageDirectories, packageDirectories...)
	err = parsePackage(allPackageDirectories, packageDirectories...)
	if err != nil {
		return err
	}

	return nil
}

func parseFile(file fs.DirEntry, filePath string) (goFile *parserGoFile, err error) {
	astFile, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return goFile, err
	}

	return &parserGoFile{
		AstFile: astFile,
		Path:    filePath,
		Name:    file.Name(),
	}, nil
}

//ConvertFieldTypeToString convers the provided field into a string representation
func convertFieldTypeToString(field *ast.Field) (string, error) {
	expression := field.Type
	switch expression.(type) {
	case *ast.Ident:
		return identityStringConversion(expression), nil
	case *ast.ArrayType:
		return arrayStringConversion(expression), nil
	case *ast.StarExpr:
		return starStringConversion(expression), nil
	case *ast.MapType:
		return mapStringConversion(expression), nil
	case *ast.FuncType:
		theFunc, err := funcStringConversion(expression)
		if err != nil {
			return theFunc, err
		}
		return theFunc, nil
	default:
		return "", fmt.Errorf("field type not supported: %v", expression)
	}
}

//ParseParameters returns the parameters from the provided function
func parseParameters(astFunc *ast.FuncType) (parameters []*Parameter, err error) {
	funcParams := []*Parameter{}
	if astFunc.Params == nil {
		return funcParams, nil
	}
	astParams := astFunc.Params.List

	for _, astParam := range astParams {
		funcParam := &Parameter{}
		//parse parameter names if there are any names assigned
		if astParam.Names != nil {
			funcParam.Name = astParam.Names[0].Name
		}

		funcParam.Type, err = convertFieldTypeToString(astParam)
		if err != nil {
			return funcParams, err
		}
		funcParam.IsTypePointer = isPointer(astParam.Type)
		funcParams = append(funcParams, funcParam)
	}
	return funcParams, nil
}

//ParseResults returns the results from the provided function
func parseResults(astFunc *ast.FuncType) (results []*Result, err error) {
	funcResults := []*Result{}
	if astFunc.Results == nil {
		return funcResults, nil
	}
	astResults := astFunc.Results.List

	for _, astResult := range astResults {
		funcResult := &Result{}
		//parse result names if there are any names assigned
		if astResult.Names != nil {
			funcResult.Name = astResult.Names[0].Name
		}

		funcResult.Type, err = convertFieldTypeToString(astResult)
		if err != nil {
			return funcResults, err
		}
		funcResult.IsTypePointer = isPointer(astResult.Type)

		funcResults = append(funcResults, funcResult)
	}
	return funcResults, nil
}

//ParseComments returns the comments from the provided declaration
func parseComments(astDecl ast.Decl) (*CommentGroup, error) {
	var astCommentGroup *ast.CommentGroup
	commentGroup := &CommentGroup{}

	switch astDecl := astDecl.(type) {
	case *ast.FuncDecl:
		astCommentGroup = astDecl.Doc
	case *ast.GenDecl:
		astCommentGroup = astDecl.Doc
	default:
		return commentGroup, errors.New("declaration type not supported")
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
