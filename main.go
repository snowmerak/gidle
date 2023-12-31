package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/alecthomas/participle/v2"
)

const (
	LanguageGo         = "go"
	LanguageDart       = "dart"
	LanguageTypeScript = "ts"
	LanguageRust       = "rust"
)

func main() {
	inputFile := flag.String("i", "", "input file")
	outputFile := flag.String("o", "", "output file")
	lang := flag.String("l", "", "output language")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" || *lang == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		panic(err)
	}

	parser := participle.MustBuild[Grammar]()

	values, err := parser.ParseBytes("", data)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Dir(*outputFile), 0755); err != nil {
		panic(err)
	}

	var generator Generator
	switch *lang {
	case LanguageGo:
		generator = NewGoGenerator()
	case LanguageDart:
		generator = NewDartGenerator()
	// case LanguageTypeScript:
	// 	generator = NewTypeScriptGenerator()
	// case LanguageRust:
	// 	generator = NewRustGenerator()
	default:
		panic("unknown language")
	}

	if err := generator.Generate(*outputFile, values); err != nil {
		panic(err)
	}
}
