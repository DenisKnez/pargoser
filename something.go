package pargoser

type DifferentCat struct {
}

// one gen decl
var ranVariable string

// one gen decl
const sdfKapow string = "something"

// this is one gen decl
const (
	//sometihng
	oneone string = "skdfj"
	//something
	twotwo string = "slkdjfs"
)

// this is one gen decl
var (
	three int          = 3
	four  *string      = &ranVariable
	five  DifferentCat = DifferentCat{}
)

// type ss func(s int) (err error)

// // some comment
// func something(ss int, s ss) (r DifferentCat) {

// }

// //makaw like a doc comment
// func makaw(meow string) *int {
