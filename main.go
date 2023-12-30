package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

const grammar = `
object Person {
	int32 id
	string name
	list of string phoneNumbers
	map int8 for string phoneBook
}

object Test {
	int32 id
	string name
	list of float64 phoneNumbers
	map uint32 for string phoneBook
}

const int8 {
	A = 1
	B = 2
	C = 3
}
`

func main() {
	parser := participle.MustBuild[Grammar]()

	v, err := parser.ParseString("", grammar)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", v)
}
