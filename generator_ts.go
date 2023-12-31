package main

import (
	"bytes"
	"os"
	"strconv"
)

type TypeScriptGenerator struct {
	buffer *bytes.Buffer
}

func NewTypeScriptGenerator() *TypeScriptGenerator {
	return &TypeScriptGenerator{
		buffer: bytes.NewBuffer(nil),
	}
}

func (t *TypeScriptGenerator) Generate(outPath string, values *Grammar) error {
	t.buffer.Reset()

	for _, name := range values.Package.Names {
		t.buffer.WriteString("export namespace ")
		t.buffer.WriteString(name)
		t.buffer.WriteString(" {\n")
	}
	t.buffer.WriteString("\n")

	for _, entry := range values.Entries {
		if entry.Const != nil {
			t.generateConst(entry.Const)
		} else if entry.Enum != nil {
			t.generateEnum(entry.Enum)
		} else if entry.Object != nil {
			t.generateObject(entry.Object)
		}
	}

	t.buffer.WriteString("\n")
	for range values.Package.Names {
		t.buffer.WriteString("}\n")
	}

	if err := os.WriteFile(outPath, t.buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (t *TypeScriptGenerator) generatePrimitiveValue(value *PrimitiveValue) {
	if value.StringValue != nil {
		t.buffer.WriteString("\"")
		t.buffer.WriteString(*value.StringValue)
		t.buffer.WriteString("\"")
	} else if value.IntValue != nil {
		t.buffer.WriteString(strconv.FormatInt(*value.IntValue, 10))
	} else if value.FloatValue != nil {
		t.buffer.WriteString(strconv.FormatFloat(*value.FloatValue, 'f', -1, 64))
	} else if value.BoolValue != nil {
		t.buffer.WriteString(strconv.FormatBool(*value.BoolValue))
	} else {
		t.buffer.WriteString("null")
	}
}

func (t *TypeScriptGenerator) generateConst(constant *Const) {
	for _, f := range constant.Fields {
		t.buffer.WriteString("export const ")
		t.buffer.WriteString(constant.Name)
		t.buffer.WriteString("_")
		t.buffer.WriteString(f.Name)
		t.buffer.WriteString(" = ")
		t.generatePrimitiveValue(&f.Value)
		t.buffer.WriteString(";\n")
	}
}

func (t *TypeScriptGenerator) generateEnum(enum *Enum) {
	t.buffer.WriteString("export enum ")
	t.buffer.WriteString(enum.Name)
	t.buffer.WriteString(" {\n")
	for _, v := range enum.Body {
		t.buffer.WriteString("\t")
		t.buffer.WriteString(v.Name)
		t.buffer.WriteString(" = ")
		t.generatePrimitiveValue(&v.Value)
		t.buffer.WriteString(",\n")
	}
	t.buffer.WriteString("}\n\n")

	t.buffer.WriteString("export function indexOf")
	t.buffer.WriteString(enum.Name)
	t.buffer.WriteString("(value: ")
	t.buffer.WriteString(enum.Name)
	t.buffer.WriteString("): number {\n")
	t.buffer.WriteString("\t switch (value) {\n")
	for i, v := range enum.Body {
		t.buffer.WriteString("\t\t case ")
		t.buffer.WriteString(enum.Name)
		t.buffer.WriteString(".")
		t.buffer.WriteString(v.Name)
		t.buffer.WriteString(":\n")
		t.buffer.WriteString("\t\t\t return ")
		t.buffer.WriteString(strconv.Itoa(i))
		t.buffer.WriteString(";\n")
	}
	t.buffer.WriteString("\t\t default:\n")
	t.buffer.WriteString("\t\t\t return -1;\n")
	t.buffer.WriteString("\t }\n")
	t.buffer.WriteString("}\n\n")

	t.buffer.WriteString("export function get")
	t.buffer.WriteString(enum.Name)
	t.buffer.WriteString("(index: number): ")
	t.buffer.WriteString(enum.Name)
	t.buffer.WriteString(" {\n")
	t.buffer.WriteString("\t switch (index) {\n")
	for i, v := range enum.Body {
		t.buffer.WriteString("\t\t case ")
		t.buffer.WriteString(strconv.Itoa(i))
		t.buffer.WriteString(":\n")
		t.buffer.WriteString("\t\t\t return ")
		t.buffer.WriteString(enum.Name)
		t.buffer.WriteString(".")
		t.buffer.WriteString(v.Name)
		t.buffer.WriteString(";\n")
	}
	t.buffer.WriteString("\t\t default:\n")
	t.buffer.WriteString("\t\t\t throw new Error(\"unknown enum value\");\n")
	t.buffer.WriteString("\t }\n")
	t.buffer.WriteString("}\n\n")
}

func (t *TypeScriptGenerator) generateType(ty *Type) {
	if ty.PrimitiveType != nil {
		t.generatePrimitiveType(ty.PrimitiveType)
	} else if ty.ListType != nil {
		t.generateListType(ty.ListType)
	} else if ty.MapType != nil {
		t.generateMapType(ty.MapType)
	} else if ty.Identity != nil {
		t.buffer.WriteString(*ty.Identity)
	} else {
		panic("unknown type")
	}
}

func (t *TypeScriptGenerator) generatePrimitiveType(ty *PrimitiveType) {
	switch ty.Type {
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		t.buffer.WriteString("number")
	case "string":
		t.buffer.WriteString("string")
	case "bool":
		t.buffer.WriteString("boolean")
	default:
		panic("unknown primitive type")
	}
}

func (t *TypeScriptGenerator) generateListType(ty *ListType) {
	t.buffer.WriteString("Array<")
	t.generateType(&ty.ElementType)
	t.buffer.WriteString(">")
}

func (t *TypeScriptGenerator) generateMapType(ty *MapType) {
	t.buffer.WriteString("Map<")
	t.generatePrimitiveType(&ty.KeyType)
	t.buffer.WriteString(",")
	t.generatePrimitiveType(&ty.ValueType)
	t.buffer.WriteString(">")
}

func (t *TypeScriptGenerator) generateObject(object *Object) {
	t.buffer.WriteString("export interface ")
	t.buffer.WriteString(object.Name)
	t.buffer.WriteString(" {\n")

	for _, f := range object.Fields {
		t.buffer.WriteString("\t")
		t.buffer.WriteString(f.Name)
		t.buffer.WriteString(": ")
		t.generateType(&f.Type)
		t.buffer.WriteString(";\n")
	}

	t.buffer.WriteString("}\n\n")
}
