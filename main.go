package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

const grammar = `
object Person {
    string name
    int32 age
    list of string friends
    map string for string properties
}

enum CASE for uint8 {
    UPPER = 0
    LOWER = 1
}

const BOUNDARY for int32 {
    MAX = 100
    MIN = 0
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
