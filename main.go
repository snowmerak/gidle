package main

import (
	"os"

	"github.com/alecthomas/participle/v2"
)

func main() {
	data, err := os.ReadFile("./test.gidle")
	if err != nil {
		panic(err)
	}

	parser := participle.MustBuild[Grammar]()

	values, err := parser.ParseBytes("", data)
	if err != nil {
		panic(err)
	}

	if err := NewGoGenerator().Generate("test.go", values); err != nil {
		panic(err)
	}
}
