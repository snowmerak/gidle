package main

import (
	"bytes"
	"os"
	"strconv"
)

type DartGenerator struct {
	buffer *bytes.Buffer
}

func NewDartGenerator() *DartGenerator {
	return &DartGenerator{
		buffer: bytes.NewBuffer(nil),
	}
}

func (d *DartGenerator) Generate(outPath string, values *Grammar) error {
	d.buffer.Reset()

	for _, entry := range values.Entries {
		if entry.Const != nil {
			d.generateConst(entry.Const)
		} else if entry.Enum != nil {
			d.generateEnum(entry.Enum)
		} else if entry.Object != nil {
			d.generateObject(entry.Object)
		}
	}

	if err := os.WriteFile(outPath, d.buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (d *DartGenerator) generatePrimitiveValue(value *PrimitiveValue) {
	if value.StringValue != nil {
		d.buffer.WriteString("\"")
		d.buffer.WriteString(*value.StringValue)
		d.buffer.WriteString("\"")
	} else if value.IntValue != nil {
		d.buffer.WriteString(strconv.FormatInt(*value.IntValue, 10))
	} else if value.FloatValue != nil {
		d.buffer.WriteString(strconv.FormatFloat(*value.FloatValue, 'f', -1, 64))
	} else if value.BoolValue != nil {
		d.buffer.WriteString(strconv.FormatBool(*value.BoolValue))
	} else {
		d.buffer.WriteString("null")
	}
}

func (d *DartGenerator) generateConst(constant *Const) {
	for _, f := range constant.Fields {
		d.buffer.WriteString("const ")
		d.buffer.WriteString(constant.Name)
		d.buffer.WriteString("_")
		d.buffer.WriteString(f.Name)
		d.buffer.WriteString(" = ")
		d.generatePrimitiveValue(&f.Value)
		d.buffer.WriteString(";\n")
	}
	d.buffer.WriteString("\n")
}

func (d *DartGenerator) generateEnum(enum *Enum) {
	d.buffer.WriteString("enum ")
	d.buffer.WriteString(enum.Name)
	d.buffer.WriteString(" {\n")
	for _, v := range enum.Body {
		d.buffer.WriteString("\t")
		d.buffer.WriteString(v.Name)
		d.buffer.WriteString(",\n")
	}
	d.buffer.WriteString("}\n\n")
}

func (d *DartGenerator) generateType(t *Type) {
	if t.PrimitiveType != nil {
		d.generatePrimitiveType(t.PrimitiveType)
	} else if t.ListType != nil {
		d.generateListType(t.ListType)
	} else if t.MapType != nil {
		d.generateMapType(t.MapType)
	} else if t.Identity != nil {
		d.buffer.WriteString(*t.Identity)
	} else {
		d.buffer.WriteString("unknown type")
	}
}

func (d *DartGenerator) generatePrimitiveType(t *PrimitiveType) {
	switch t.Type {
	case "string":
		d.buffer.WriteString("String")
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
		d.buffer.WriteString("int")
	case "float32", "float64":
		d.buffer.WriteString("double")
	case "bool":
		d.buffer.WriteString("bool")
	default:
		d.buffer.WriteString("unknown type")
	}
}

func (d *DartGenerator) generateListType(t *ListType) {
	d.buffer.WriteString("List<")
	d.generateType(&t.ElementType)
	d.buffer.WriteString(">")
}

func (d *DartGenerator) generateMapType(t *MapType) {
	d.buffer.WriteString("Map<")
	d.generatePrimitiveType(&t.KeyType)
	d.buffer.WriteString(", ")
	d.generatePrimitiveType(&t.ValueType)
	d.buffer.WriteString(">")
}

func (d *DartGenerator) generateObject(object *Object) {
	d.buffer.WriteString("import 'dart:convert';\n\n")

	d.buffer.WriteString("class ")
	d.buffer.WriteString(object.Name)
	d.buffer.WriteString(" {\n")

	for _, f := range object.Fields {
		d.buffer.WriteString("\t")
		d.generateType(&f.Type)
		d.buffer.WriteString("? ")
		d.buffer.WriteString(SnakeToCamel(f.Name))
		d.buffer.WriteString(";\n")
	}
	d.buffer.WriteString("\n")

	d.buffer.WriteString("\t")
	d.buffer.WriteString(object.Name)
	d.buffer.WriteString("({\n")
	for _, f := range object.Fields {
		d.buffer.WriteString("\t\tthis.")
		d.buffer.WriteString(SnakeToCamel(f.Name))
		d.buffer.WriteString(",\n")
	}
	d.buffer.WriteString("\t});\n\n")

	d.buffer.WriteString("\ttoMap() {\n")
	d.buffer.WriteString("\t\treturn {\n")
	for _, f := range object.Fields {
		d.buffer.WriteString("\t\t\t\"")
		d.buffer.WriteString(f.Name)
		d.buffer.WriteString("\": ")
		d.buffer.WriteString(SnakeToCamel(f.Name))
		d.buffer.WriteString(",\n")
	}
	d.buffer.WriteString("\t\t};\n")
	d.buffer.WriteString("\t}\n\n")

	d.buffer.WriteString("\tString toJson() {\n")
	d.buffer.WriteString("\t\treturn jsonEncode(toMap());\n")
	d.buffer.WriteString("\t}\n\n")

	d.buffer.WriteString("\t")
	d.buffer.WriteString(object.Name)
	d.buffer.WriteString(".fromMap(Map<String, dynamic> map) {\n")
	for _, f := range object.Fields {
		d.buffer.WriteString("\t\t")
		d.buffer.WriteString(SnakeToCamel(f.Name))
		d.buffer.WriteString(" = ")
		switch IsObjectType(&f.Type) {
		case false:
			switch IsPrimitiveType(&f.Type) {
			case true:
				d.buffer.WriteString("map[\"")
				d.buffer.WriteString(f.Name)
				d.buffer.WriteString("\"]")
			case false:
				d.generateType(&f.Type)
				d.buffer.WriteString(".from(")
				d.buffer.WriteString("map[\"")
				d.buffer.WriteString(f.Name)
				d.buffer.WriteString("\"]")
				d.buffer.WriteString(")")
			}

		default:
			d.buffer.WriteString(*f.Type.Identity)
			d.buffer.WriteString(".fromMap(map[\"")
			d.buffer.WriteString(f.Name)
			d.buffer.WriteString("\"] ?? {})")
		}
		d.buffer.WriteString(";\n")
	}
	d.buffer.WriteString("\t}\n\n")

	d.buffer.WriteString("\t")
	d.buffer.WriteString(object.Name)
	d.buffer.WriteString(".fromJson(String source) : this.fromMap(jsonDecode(source));\n\n")

	d.buffer.WriteString("}\n\n")
}
