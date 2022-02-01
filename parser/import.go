package parser

import (
	"go/ast"
)

func getImports(file parserGoFile) (imports []*Import) {
	genDeclarations := parseGenDeclarations(file)
	genImportDeclarations := parseImportDecls(genDeclarations)
	fileImports := convertImportDeclsIntoImport(genImportDeclarations)
	imports = append(imports, fileImports...)
	return imports
}

func parseImportDecls(genDeclarations []*ast.GenDecl) (imports []*ast.GenDecl) {
	genDeclsWithImportSpec := []*ast.GenDecl{}
	//loop over all general declarations in the file
	for _, genDeclaration := range genDeclarations {
		genDeclarationSpec := genDeclaration.Specs[0]
		switch genDeclarationSpec.(type) {
		case *ast.ImportSpec:
			genDeclsWithImportSpec = append(genDeclsWithImportSpec, genDeclaration)
		default:
			continue
		}

	}

	return genDeclsWithImportSpec
}

func parseImportsByPackage(file parserGoFile) (imports []*Import) {
	genDeclarations := parseGenDeclarations(file)
	genImportDeclarations := parseImportDecls(genDeclarations)
	fileImports := convertImportDeclsIntoImport(genImportDeclarations)
	imports = append(imports, fileImports...)

	return imports
}

func convertImportDeclsIntoImport(genImportDecls []*ast.GenDecl) (imports []*Import) {
	for _, genImportDecl := range genImportDecls {
		for _, spec := range genImportDecl.Specs {
			theImport := &Import{}
			importSpec := spec.(*ast.ImportSpec)

			// PATH & NAME
			path := importSpec.Path
			name := importSpec.Name

			// DOC
			if importSpec.Doc != nil {
				theImport.Doc = &CommentGroup{}
				theImport.Doc.Comments = []*Comment{}
				for _, docComment := range importSpec.Doc.List {
					theImport.Doc.Comments = append(theImport.Doc.Comments, &Comment{
						Text: docComment.Text,
					})
				}
			}

			// COMMENT
			if importSpec.Comment != nil {
				theImport.Comment = &CommentGroup{}
				theImport.Comment.Comments = []*Comment{}
				for _, comment := range importSpec.Comment.List {
					theImport.Comment.Comments = append(theImport.Comment.Comments, &Comment{
						Text: comment.Text,
					})
				}
			}

			// if the name is provided set it if not
			// set null (naming import paths is optional)
			if name == nil {
				theImport.Name = nil
			} else {
				theImport.Name = &name.Name
			}

			theImport.Path = path.Value

			imports = append(imports, theImport)
		}
	}
	return imports
}
