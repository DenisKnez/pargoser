package main

//NotPerson this is a person interface
type NotPerson interface {
	Hello(string) string
	Add(int, int) int
	How(string, string) float32
}

func Meow(string, string) int {
	return 2
}
