package parser

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
	Doc     *CommentGroup `json:"doc"`
	Name    string        `json:"name"`
	Methods []*Method     `json:"methods"`
}

type Variable struct {
	Doc   *CommentGroup
	Kind  VariableKind `json:"kind"`
	Name  string       `json:"name"`
	Value *string      `json:"value"`
}

type Import struct {
	Doc     *CommentGroup `json:"doc"`
	Name    *string       `json:"name"`
	Path    string        `json:"path"`
	Comment *CommentGroup `json:"comment"`
}

type TypeSpec struct {
	Doc     *CommentGroup `json:"doc"`
	Name    *string       `json:"name"`
	Type    string        `json:"type"`
	Comment *CommentGroup `json:"comment"`
}

type Struct struct {
	Doc     *CommentGroup `json:"doc"`
	Name    string        `json:"name"`
	Fields  []*Field      `json:"fields"`
	Methods []*Method     `json:"methods"`
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

//TODO anonimous functions

//Function represents a function
type Function struct {
	Doc     *CommentGroup `json:"doc"`
	Name    string        `json:"name"`
	Params  []*Parameter  `json:"params"`
	Results []*Result     `json:"results"`
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

//CommentGroup comments with 1 line line or more grouped together
type CommentGroup struct {
	Comments []*Comment `json:"comments"`
}
