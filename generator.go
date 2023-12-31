package main

type Generator interface {
	Generate(outPath string, values *Grammar) error
}
