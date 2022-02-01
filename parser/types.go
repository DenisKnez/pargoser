package parser

import (
	"bytes"
	"go/token"
	"text/template"
)

var golangBasicDataTypes = []string{
	"bool",
	"string",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",
	"byte",
	"rune",
	"float32",
	"float64",
	"complex64",
	"complex128",
}

type GolangType int

const (
	INTERFACE GolangType = iota
	STRUCT
	FUNCTION
	METHOD
)

func (gt GolangType) String() string {
	return [...]string{"interface", "struct", "function", "method"}[gt]
}

type VariableKind int

const (
	Const VariableKind = iota
	Var
)

func (vk VariableKind) String() string {
	return [...]string{"const", "var"}[vk]
}

type Interface struct {
	PackageName string        `json:"packageName"`
	Doc         *CommentGroup `json:"doc"`
	Name        string        `json:"name"`
	Methods     []*Method     `json:"methods"`
}

type Variable struct {
	PackageName string `json:"packageName"`
	Doc         *CommentGroup
	Kind        VariableKind `json:"kind"`
	Name        string       `json:"name"`
	Value       *string      `json:"value"`
}

type Import struct {
	PackageName string        `json:"packageName"`
	Doc         *CommentGroup `json:"doc"`
	Name        *string       `json:"name"`
	Path        string        `json:"path"`
	Comment     *CommentGroup `json:"comment"`
}

type Struct struct {
	PackageName string        `json:"packageName"`
	Doc         *CommentGroup `json:"doc"`
	Name        string        `json:"name"`
	Fields      []*Field      `json:"fields"`
	Methods     []*Method     `json:"methods"`
}

const (
	StructTemplate    string = "StructTemplate"
	InterfaceTemplate string = "InterfaceTemplate"
	ImportTemplate    string = "ImportTemplate"
)

func (s Struct) String() (string, error) {
	declarationString, err := generateTypeDeclaration(StructTemplate, StructTemplateString, s)
	if err != nil {
		return declarationString, err
	}

	return declarationString, nil
}

func (i Interface) String() (string, error) {
	declarationString, err := generateTypeDeclaration(InterfaceTemplate, InterfaceTemplateString, i)
	if err != nil {
		return declarationString, err
	}

	return declarationString, nil
}

func (i Import) String() (string, error) {
	declarationString, err := generateTypeDeclaration(ImportTemplate, ImportTemplateString, i)
	if err != nil {
		return declarationString, err
	}
	return declarationString, nil
}

func generateTypeDeclaration(templateName string, templateString string, object interface{}) (declarationString string, err error) {
	tmpl, err := template.New(templateName).Parse(templateString)
	if err != nil {
		return declarationString, err
	}

	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, object)
	if err != nil {
		return declarationString, err
	}

	return buf.String(), nil
}

//Field represents different things depending on the type
// Struct it represents the struct field
// Interface it represents the method list
// Declaration signature it represents parameters/results
type Field struct {
	Doc           *CommentGroup `json:"doc"`
	Name          string        `json:"name"`
	IsTypePointer bool          `json:"isTypePointer"`
	Type          string        `json:"type"`
	Tag           *Tag          `json:"tag"`
}

type TagValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Tag struct {
	Type  token.Token `json:"type"`
	Value TagValue    `json:"value"`
}

type Receiver struct {
	PointerReceiver bool   `json:"pointerReceiver"`
	Type            string `json:"type"`
	Name            string `json:"name"`
}

//Method represents a struct or interface method
type Method struct {
	Doc      *CommentGroup `json:"doc"`
	Receiver *Receiver     `json:"receiver"`
	Name     string        `json:"name"`
	Params   []*Parameter  `json:"params"`
	Results  []*Result     `json:"results"`
}

type Package struct {
	Name          string   `json:"name"`
	Files         []GoFile `json:"files"`
	DirectoryPath string   `json:"directoryPath"`
}

type GoFile struct {
	Structs    []*Struct    `json:"structs"`
	Interfaces []*Interface `json:"interfaces"`
	Imports    []*Import    `json:"imports"`
	Variables  []Variable   `json:"variables"`
	Functions  []*Function  `json:"functions"`
}

// TODO: anonimous functions

//Function represents a function
type Function struct {
	PackageName string        `json:"packageName"`
	Doc         *CommentGroup `json:"doc"`
	Name        string        `json:"name"`
	Params      []*Parameter  `json:"params"`
	Results     []*Result     `json:"results"`
}

//Parameter a variable that is passed into a method/function
type Parameter struct {
	Name          string `json:"name"`
	IsTypePointer bool   `json:"isTypePointer"`
	Type          string `json:"type"`
}

//Result result is a variable returned from a function or method
type Result struct {
	Name          string `json:"name"`
	IsTypePointer bool   `json:"isTypePointer"`
	Type          string `json:"type"`
}

//Comment represents a single line of a comment
type Comment struct {
	Text string `json:"text"`
}

//CommentGroup comments with 1 line or more grouped together
type CommentGroup struct {
	PackageName string     `json:"packageName"`
	Comments    []*Comment `json:"comments"`
}
