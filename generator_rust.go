package main

import (
	"bytes"
	"os"
	"strconv"
)

type RustGenerator struct {
	buffer *bytes.Buffer
}

func NewRustGenerator() *RustGenerator {
	return &RustGenerator{
		buffer: bytes.NewBuffer(nil),
	}
}

func (r *RustGenerator) Generate(outPath string, values *Grammar) error {
	r.buffer.Reset()

	r.buffer.WriteString("use std::collections::HashMap;\n")
	r.buffer.WriteString("use serde::{Deserialize, Serialize};\n")
	r.buffer.WriteString("use serde_json::{to_string, from_str, Result};\n")
	r.buffer.WriteString("\n")

	for _, entry := range values.Entries {
		if entry.Const != nil {
			r.generateConst(entry.Const)
		} else if entry.Enum != nil {
			r.generateEnum(entry.Enum)
		} else if entry.Object != nil {
			r.generateObject(entry.Object)
		}
	}

	if err := os.WriteFile(outPath, r.buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (r *RustGenerator) generatePrimitiveValue(value *PrimitiveValue) {
	if value.StringValue != nil {
		r.buffer.WriteString("\"")
		r.buffer.WriteString(*value.StringValue)
		r.buffer.WriteString("\"")
	} else if value.IntValue != nil {
		r.buffer.WriteString(strconv.FormatInt(*value.IntValue, 10))
	} else if value.FloatValue != nil {
		r.buffer.WriteString(strconv.FormatFloat(*value.FloatValue, 'f', -1, 64))
	} else if value.BoolValue != nil {
		r.buffer.WriteString(strconv.FormatBool(*value.BoolValue))
	} else {
		r.buffer.WriteString("null")
	}
}

func (r *RustGenerator) generateConst(constant *Const) {
	for _, f := range constant.Fields {
		r.buffer.WriteString("pub const ")
		r.buffer.WriteString(constant.Name)
		r.buffer.WriteString("_")
		r.buffer.WriteString(f.Name)
		r.buffer.WriteString(": ")
		r.generatePrimitiveType(&constant.Type)
		r.buffer.WriteString(" = ")
		r.generatePrimitiveValue(&f.Value)
		r.buffer.WriteString(";\n")
	}
}

func (r *RustGenerator) generateEnum(enum *Enum) {
	r.buffer.WriteString("pub enum ")
	r.buffer.WriteString(enum.Name)
	r.buffer.WriteString(" {\n")
	for _, v := range enum.Body {
		r.buffer.WriteString("\t")
		r.buffer.WriteString(v.Name)
		r.buffer.WriteString(" = ")
		r.generatePrimitiveValue(&v.Value)
		r.buffer.WriteString(",\n")
	}
	r.buffer.WriteString("}\n\n")
}

func (r *RustGenerator) generateType(t *Type) {
	if t.PrimitiveType != nil {
		r.generatePrimitiveType(t.PrimitiveType)
	} else if t.ListType != nil {
		r.buffer.WriteString("Vec<")
		r.generateType(&t.ListType.ElementType)
		r.buffer.WriteString(">")
	} else if t.MapType != nil {
		r.buffer.WriteString("HashMap<")
		r.generatePrimitiveType(&t.MapType.KeyType)
		r.buffer.WriteString(", ")
		r.generatePrimitiveType(&t.MapType.ValueType)
		r.buffer.WriteString(">")
	} else if t.Identity != nil {
		r.buffer.WriteString(*t.Identity)
	} else {
		panic("unreachable")
	}
}

func (r *RustGenerator) generatePrimitiveType(t *PrimitiveType) {
	switch t.Type {
	case "int8":
		r.buffer.WriteString("i8")
	case "int16":
		r.buffer.WriteString("i16")
	case "int32":
		r.buffer.WriteString("i32")
	case "int64":
		r.buffer.WriteString("i64")
	case "uint8":
		r.buffer.WriteString("u8")
	case "uint16":
		r.buffer.WriteString("u16")
	case "uint32":
		r.buffer.WriteString("u32")
	case "uint64":
		r.buffer.WriteString("u64")
	case "float32":
		r.buffer.WriteString("f32")
	case "float64":
		r.buffer.WriteString("f64")
	case "string":
		r.buffer.WriteString("String")
	case "bool":
		r.buffer.WriteString("bool")
	default:
		panic("unreachable")
	}
}

func (r *RustGenerator) generateObject(object *Object) {
	r.buffer.WriteString("#[derive(Debug, Serialize, Deserialize)]\n")
	r.buffer.WriteString("pub struct ")
	r.buffer.WriteString(object.Name)
	r.buffer.WriteString(" {\n")
	for _, f := range object.Fields {
		r.buffer.WriteString("\tpub ")
		r.buffer.WriteString(f.Name)
		r.buffer.WriteString(": ")
		r.generateType(&f.Type)
		r.buffer.WriteString(",\n")
	}
	r.buffer.WriteString("}\n\n")

	r.buffer.WriteString("impl ")
	r.buffer.WriteString(object.Name)
	r.buffer.WriteString(" {\n")

	r.buffer.WriteString("\tpub fn new(")
	for i, f := range object.Fields {
		if i > 0 {
			r.buffer.WriteString(", ")
		}
		r.buffer.WriteString(f.Name)
		r.buffer.WriteString(": ")
		r.generateType(&f.Type)
	}
	r.buffer.WriteString(") -> Self {\n")
	r.buffer.WriteString("\t\tSelf {\n")
	for _, f := range object.Fields {
		r.buffer.WriteString("\t\t\t")
		r.buffer.WriteString(f.Name)
		r.buffer.WriteString(",\n")
	}
	r.buffer.WriteString("\t\t}\n")
	r.buffer.WriteString("\t}\n")

	r.buffer.WriteString("\tpub fn to_json(&self) -> Result<String> {\n")
	r.buffer.WriteString("\t\tto_string(self)\n")
	r.buffer.WriteString("\t}\n")

	r.buffer.WriteString("\tpub fn from_json(json: &str) -> Result<Self> {\n")
	r.buffer.WriteString("\t\tfrom_str(json)\n")
	r.buffer.WriteString("\t}\n")

	r.buffer.WriteString("}\n\n")
}
