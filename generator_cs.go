package main

import (
	"bytes"
	"os"
	"strconv"
)

type CSharpGenerator struct {
	buffer *bytes.Buffer
}

func NewCSharpGenerator() *CSharpGenerator {
	return &CSharpGenerator{buffer: bytes.NewBuffer(nil)}
}

func (cs *CSharpGenerator) Generate(outPath string, values *Grammar) error {
	cs.buffer.Reset()

	cs.buffer.WriteString("using System.Text.Json;\n")
	cs.buffer.WriteString("using System.Text.Json.Serialization;\n")
	cs.buffer.WriteString("\n")

	for _, name := range values.Package.Names {
		cs.buffer.WriteString("namespace ")
		cs.buffer.WriteString(name)
		cs.buffer.WriteString(" {\n")
	}

	for _, entry := range values.Entries {
		if entry.Const != nil {
			cs.generateConst(entry.Const)
		} else if entry.Enum != nil {
			cs.generateEnum(entry.Enum)
		} else if entry.Object != nil {
			cs.generateObject(entry.Object)
		}
	}

	for range values.Package.Names {
		cs.buffer.WriteString("}\n")
	}

	if err := os.WriteFile(outPath, cs.buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (cs *CSharpGenerator) generatePrimitiveValue(value *PrimitiveValue) {
	if value.StringValue != nil {
		cs.buffer.WriteString("\"")
		cs.buffer.WriteString(*value.StringValue)
		cs.buffer.WriteString("\"")
	} else if value.IntValue != nil {
		cs.buffer.WriteString(strconv.FormatInt(*value.IntValue, 10))
	} else if value.FloatValue != nil {
		cs.buffer.WriteString(strconv.FormatFloat(*value.FloatValue, 'f', -1, 64))
	} else if value.BoolValue != nil {
		cs.buffer.WriteString(strconv.FormatBool(*value.BoolValue))
	} else {
		cs.buffer.WriteString("null")
	}
}

func (cs *CSharpGenerator) generateConst(constant *Const) {
	cs.buffer.WriteString("public static class ")
	cs.buffer.WriteString(constant.Name)
	cs.buffer.WriteString(" {\n")

	for _, f := range constant.Fields {
		cs.buffer.WriteString("public const ")
		cs.generatePrimitiveType(&constant.Type)
		cs.buffer.WriteString(" ")
		cs.buffer.WriteString(f.Name)
		cs.buffer.WriteString(" = ")
		cs.generatePrimitiveValue(&f.Value)
		cs.buffer.WriteString(";\n")
	}

	cs.buffer.WriteString("}\n")
}

func (cs *CSharpGenerator) generateEnum(enum *Enum) {
	cs.buffer.WriteString("public class ")
	cs.buffer.WriteString(enum.Name)
	cs.buffer.WriteString(" {\n")

	for _, v := range enum.Body {
		cs.buffer.WriteString("\tpublic const ")
		cs.generatePrimitiveType(&enum.Type)
		cs.buffer.WriteString(" ")
		cs.buffer.WriteString(v.Name)
		cs.buffer.WriteString(" = ")
		cs.generatePrimitiveValue(&v.Value)
		cs.buffer.WriteString(";\n")
	}

	cs.buffer.WriteString("public static int IndexOf(")
	cs.generatePrimitiveType(&enum.Type)
	cs.buffer.WriteString(" value) {\n")
	cs.buffer.WriteString("return value switch {\n")
	for i, v := range enum.Body {
		cs.buffer.WriteString(v.Name)
		cs.buffer.WriteString(" => ")
		cs.buffer.WriteString(strconv.Itoa(i))
		cs.buffer.WriteString(",\n")
	}
	cs.buffer.WriteString("_ => throw new ArgumentOutOfRangeException(nameof(value), \"Invalid value\")\n")
	cs.buffer.WriteString("};\n")
	cs.buffer.WriteString("}\n")

	cs.buffer.WriteString("public static ")
	cs.generatePrimitiveType(&enum.Type)
	cs.buffer.WriteString(" ValueOf(int index) {\n")
	cs.buffer.WriteString("return index switch {\n")
	for i, v := range enum.Body {
		cs.buffer.WriteString(strconv.Itoa(i))
		cs.buffer.WriteString(" => ")
		cs.buffer.WriteString(v.Name)
		cs.buffer.WriteString(",\n")
	}
	cs.buffer.WriteString("_ => throw new ArgumentOutOfRangeException(nameof(index), \"Invalid index\")\n")
	cs.buffer.WriteString("};\n")

	cs.buffer.WriteString("}\n")

	cs.buffer.WriteString("}\n")
}

func (cs *CSharpGenerator) generateType(t *Type) {
	if t.PrimitiveType != nil {
		cs.generatePrimitiveType(t.PrimitiveType)
	} else if t.ListType != nil {
		cs.generateListType(t.ListType)
	} else if t.MapType != nil {
		cs.generateMapType(t.MapType)
	} else if t.Identity != nil {
		cs.buffer.WriteString(*t.Identity)
	} else {
		cs.buffer.WriteString("unknown type")
	}
}

func (cs *CSharpGenerator) generatePrimitiveType(t *PrimitiveType) {
	switch t.Type {
	case "uint8":
		cs.buffer.WriteString("byte")
	case "uint16":
		cs.buffer.WriteString("ushort")
	case "uint32":
		cs.buffer.WriteString("uint")
	case "uint64":
		cs.buffer.WriteString("ulong")
	case "int8":
		cs.buffer.WriteString("sbyte")
	case "int16":
		cs.buffer.WriteString("short")
	case "int32":
		cs.buffer.WriteString("int")
	case "int64":
		cs.buffer.WriteString("long")
	case "float32":
		cs.buffer.WriteString("float")
	case "float64":
		cs.buffer.WriteString("double")
	case "bool":
		cs.buffer.WriteString("bool")
	case "string":
		cs.buffer.WriteString("string")
	default:
		cs.buffer.WriteString("unknown type")
	}
}

func (cs *CSharpGenerator) generateListType(t *ListType) {
	cs.buffer.WriteString("List<")
	cs.generateType(&t.ElementType)
	cs.buffer.WriteString(">")
}

func (cs *CSharpGenerator) generateMapType(t *MapType) {
	cs.buffer.WriteString("Dictionary<")
	cs.generatePrimitiveType(&t.KeyType)
	cs.buffer.WriteString(", ")
	cs.generatePrimitiveType(&t.ValueType)
	cs.buffer.WriteString(">")
}

func (cs *CSharpGenerator) generateObject(object *Object) {
	cs.buffer.WriteString("public class ")
	cs.buffer.WriteString(object.Name)
	cs.buffer.WriteString(" {\n")

	for _, f := range object.Fields {
		cs.buffer.WriteString("[JsonPropertyName(")
		cs.buffer.WriteString("\"")
		cs.buffer.WriteString(f.Name)
		cs.buffer.WriteString("\")]\n")
		cs.buffer.WriteString("public ")
		cs.generateType(&f.Type)
		cs.buffer.WriteString(" ")
		cs.buffer.WriteString(SnakeToPascal(f.Name))
		cs.buffer.WriteString(" { get; set; }\n")
	}
	cs.buffer.WriteString("\n")

	cs.buffer.WriteString("public ")
	cs.buffer.WriteString(object.Name)
	cs.buffer.WriteString("(")
	for i, f := range object.Fields {
		if i > 0 {
			cs.buffer.WriteString(", ")
		}
		cs.generateType(&f.Type)
		cs.buffer.WriteString(" ")
		cs.buffer.WriteString(f.Name)
	}
	cs.buffer.WriteString(") {\n")
	for _, f := range object.Fields {
		cs.buffer.WriteString("this.")
		cs.buffer.WriteString(SnakeToPascal(f.Name))
		cs.buffer.WriteString(" = ")
		cs.buffer.WriteString(f.Name)
		cs.buffer.WriteString(";\n")
	}
	cs.buffer.WriteString("}\n\n")

	cs.buffer.WriteString("public static ")
	cs.buffer.WriteString(object.Name)
	cs.buffer.WriteString("? FromJson(string json) {\n")
	cs.buffer.WriteString("return JsonSerializer.Deserialize<")
	cs.buffer.WriteString(object.Name)
	cs.buffer.WriteString(">(json);\n")
	cs.buffer.WriteString("}\n\n")

	cs.buffer.WriteString("public string ToJson() {\n")
	cs.buffer.WriteString("return JsonSerializer.Serialize(this);\n")
	cs.buffer.WriteString("}\n\n")

	cs.buffer.WriteString("}\n")
}
