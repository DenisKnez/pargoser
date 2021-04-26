package main

type DifferentCat struct {
}

// one gen decl
var ranVariable string

// one gen decl
const sdfKapow string = "something"

// this is one gen decl
const (
	oneone string = "skdfj"
	twotwo string = "slkdjfs"
)

// this is one gen decl
var (
	three int          = 3
	four  *string      = &ranVariable
	five  DifferentCat = DifferentCat{}
)
