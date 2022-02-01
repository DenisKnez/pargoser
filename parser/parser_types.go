package parser

import (
	"go/ast"
	"io/fs"
)

// parser package, directory, file are used just for the parser, this should not be used
// outside of the internals of the library
type parserPackage struct {
	GoFiles          []*parserGoFile
	InnerDirectories []parserDirectory
	DirectoryPath    string
}

type parserDirectory struct {
	FsDirectory   fs.DirEntry
	DirectoryPath string
}

type parserGoFile struct {
	Name    string
	Path    string
	AstFile *ast.File
}
