package main

import (
	"bytes"
	"errors"
	"go/format"
	"os"
	"strconv"

	"golang.org/x/tools/imports"
)

type GoGenerator struct {
	buffer *bytes.Buffer
}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{
		buffer: bytes.NewBuffer(nil),
	}
}

func (g *GoGenerator) Generate(outPath string, values *Grammar) error {
	g.buffer.Reset()

	g.buffer.WriteString("package ")
	g.buffer.WriteString(values.Package.Names[len(values.Package.Names)-1])
	g.buffer.WriteString("\n\n")

	// g.buffer.WriteString("import (\n")
	// g.buffer.WriteString("\t\"errors\"\n")
	// g.buffer.WriteString(")\n\n")

	for _, entry := range values.Entries {
		if entry.Const != nil {
			g.generateConst(entry.Const)
		} else if entry.Enum != nil {
			g.generateEnum(entry.Enum)
		} else if entry.Object != nil {
			g.generateObject(entry.Object)
		}
	}

	formatted, err := format.Source(g.buffer.Bytes())
	if err != nil {
		return err
	}

	formatted, err = imports.Process("", formatted, nil)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		return err
	}

	return nil
}

func (g *GoGenerator) generateConst(constant *Const) error {
	g.buffer.WriteString("const (\n")
	for _, f := range constant.Fields {
		g.buffer.WriteString("\t")
		g.buffer.WriteString(constant.Name)
		g.buffer.WriteString("_")
		g.buffer.WriteString(f.Name)
		g.buffer.WriteString(" = ")
		g.generatePrimitiveValue(&f.Value)
		g.buffer.WriteString("\n")
	}
	g.buffer.WriteString(")\n\n")

	return nil
}

func (g *GoGenerator) generateEnum(enum *Enum) error {
	g.buffer.WriteString("type ")
	g.buffer.WriteString(enum.Name)
	g.buffer.WriteString(" ")
	g.buffer.WriteString(enum.Type.Type)
	g.buffer.WriteString("\n\n")

	g.buffer.WriteString("const (\n")
	for _, v := range enum.Body {
		g.buffer.WriteString("\t")
		g.buffer.WriteString(enum.Name)
		g.buffer.WriteString("_")
		g.buffer.WriteString(v.Name)
		g.buffer.WriteString(" = ")
		g.buffer.WriteString(enum.Name)
		g.buffer.WriteString("(")
		g.generatePrimitiveValue(&v.Value)
		g.buffer.WriteString(")")
		g.buffer.WriteString("\n")
	}
	g.buffer.WriteString(")\n\n")

	g.buffer.WriteString("func (e ")
	g.buffer.WriteString(enum.Name)
	g.buffer.WriteString(") String() string {\n")
	g.buffer.WriteString("\tswitch e {\n")
	for _, v := range enum.Body {
		g.buffer.WriteString("\tcase ")
		g.buffer.WriteString(enum.Name)
		g.buffer.WriteString("_")
		g.buffer.WriteString(v.Name)
		g.buffer.WriteString(":\n")
		g.buffer.WriteString("\t\treturn \"")
		g.buffer.WriteString(v.Name)
		g.buffer.WriteString("\"\n")
	}
	g.buffer.WriteString("\tdefault:\n")
	g.buffer.WriteString("\t\treturn \"unknown enum value\"\n")
	g.buffer.WriteString("\t}\n")
	g.buffer.WriteString("}\n\n")

	g.buffer.WriteString("func Get")
	g.buffer.WriteString(enum.Name)
	g.buffer.WriteString("(index int) (value ")
	g.buffer.WriteString(enum.Name)
	g.buffer.WriteString(", ok bool) {\n")
	g.buffer.WriteString("\tswitch index {\n")
	for i, v := range enum.Body {
		g.buffer.WriteString("\tcase ")
		g.buffer.WriteString(strconv.Itoa(i))
		g.buffer.WriteString(":\n")
		g.buffer.WriteString("\t\treturn ")
		g.buffer.WriteString(enum.Name)
		g.buffer.WriteString("_")
		g.buffer.WriteString(v.Name)
		g.buffer.WriteString(", true\n")
	}
	g.buffer.WriteString("\tdefault:\n")
	g.buffer.WriteString("\t\treturn ")
	g.buffer.WriteString("value, false\n")
	g.buffer.WriteString("\t}\n")
	g.buffer.WriteString("}\n\n")

	g.buffer.WriteString("func IndexOf")
	g.buffer.WriteString(enum.Name)
	g.buffer.WriteString("(value ")
	g.buffer.WriteString(enum.Name)
	g.buffer.WriteString(") (index int, ok bool) {\n")
	g.buffer.WriteString("\tswitch value {\n")
	for i, v := range enum.Body {
		g.buffer.WriteString("\tcase ")
		g.buffer.WriteString(enum.Name)
		g.buffer.WriteString("_")
		g.buffer.WriteString(v.Name)
		g.buffer.WriteString(":\n")
		g.buffer.WriteString("\t\treturn ")
		g.buffer.WriteString(strconv.Itoa(i))
		g.buffer.WriteString(", true\n")
	}
	g.buffer.WriteString("\tdefault:\n")
	g.buffer.WriteString("\t\treturn ")
	g.buffer.WriteString("index, false\n")
	g.buffer.WriteString("\t}\n")
	g.buffer.WriteString("}\n\n")

	return nil
}

func (g *GoGenerator) generatePrimitiveValue(value *PrimitiveValue) error {
	if value.StringValue != nil {
		g.buffer.WriteString("\"")
		g.buffer.WriteString(*value.StringValue)
		g.buffer.WriteString("\"")
	} else if value.IntValue != nil {
		g.buffer.WriteString(strconv.FormatInt(*value.IntValue, 10))
	} else if value.FloatValue != nil {
		g.buffer.WriteString(strconv.FormatFloat(*value.FloatValue, 'f', -1, 64))
	} else if value.BoolValue != nil {
		g.buffer.WriteString(strconv.FormatBool(*value.BoolValue))
	} else {
		return errors.New("unknown primitive value")
	}

	return nil
}

func (g *GoGenerator) generateType(t *Type) error {
	if t.PrimitiveType != nil {
		g.generatePrimitiveType(t.PrimitiveType)
	} else if t.ListType != nil {
		g.generateListType(t.ListType)
	} else if t.MapType != nil {
		g.generateMapType(t.MapType)
	} else if t.Identity != nil {
		g.buffer.WriteString(*t.Identity)
	} else {
		return errors.New("unknown type")
	}

	return nil
}

func (g *GoGenerator) generatePrimitiveType(t *PrimitiveType) error {
	g.buffer.WriteString(t.Type)

	return nil
}

func (g *GoGenerator) generateMapType(m *MapType) error {
	g.buffer.WriteString("map[")
	g.generatePrimitiveType(&m.KeyType)
	g.buffer.WriteString("]")
	g.generatePrimitiveType(&m.ValueType)

	return nil
}

func (g *GoGenerator) generateListType(list *ListType) error {
	g.buffer.WriteString("[]")
	g.generateType(&list.ElementType)

	return nil
}

func (g *GoGenerator) generateObject(object *Object) error {
	g.buffer.WriteString("type ")
	g.buffer.WriteString(object.Name)
	g.buffer.WriteString(" struct {\n")
	for _, f := range object.Fields {
		g.buffer.WriteString("\t")
		g.buffer.WriteString(SnakeToPascal(f.Name))
		g.buffer.WriteString(" ")
		g.generateType(&f.Type)
		g.buffer.WriteString(" `")
		g.buffer.WriteString("json:\"")
		g.buffer.WriteString(f.Name)
		g.buffer.WriteString("\"`")
		g.buffer.WriteString("\n")
	}
	g.buffer.WriteString("}\n\n")

	return nil
}
