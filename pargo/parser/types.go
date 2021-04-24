package parser

type GolangType string

const (
	INTERFACE GolangType = "interface"
	STRUCT               = "struct"
	FUNCTION             = "function"
	METHOD               = "method"
)

type Interface struct {
	CommentGroup *CommentGroup
	Name         string
	Methods      []*Method
}

type Struct struct {
	CommentGroup *CommentGroup
	Name         string
	Fields       []*Field
	Methods      []*Method
}

//Field represents different things depending on the type
// Struct it represents the struct field
// Interface it represents the method list
// Declaration signature it represents parameters/results
type Field struct {
	CommentGroup  *CommentGroup
	Name          string
	IsTypePointer bool
	Type          string
}

type Receiver struct {
	PointerReceiver bool
	Type            string
	Name            string
}

//Method represents a struct or interface method
type Method struct {
	CommentGroup *CommentGroup
	Receiver     *Receiver
	Name         string
	Params       []*Parameter
	Results      []*Result
}

//TODO anonimous functions

//Function represents a function
type Function struct {
	CommentGroup *CommentGroup
	Name         string
	Params       []*Parameter
	Results      []*Result
}

//Parameter a variable that is passed into a method/function
type Parameter struct {
	Name          string
	IsTypePointer bool
	Type          string
}

//Result result is a variable returned from a function or method
type Result struct {
	Name          string
	IsTypePointer bool
	Type          string
}

//Comment represents a single line of a comment
type Comment struct {
	Text string
}

//CommentGroup comments with 1 line line or more grouped together
type CommentGroup struct {
	Comments []*Comment
}
