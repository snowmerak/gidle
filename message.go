package main

type ListType struct {
	ElementType Type `"list" "of" @@`
}

type ListValue struct {
	Values []Value `"[" (@@ ("," @@)*)? "]"`
}

type MapType struct {
	KeyType   PrimitiveType `"map" @@`
	ValueType PrimitiveType `"for" @@`
}

type MapEntry struct {
	Key   Value `@@`
	Value Value `"of" @@`
}

type MapValue struct {
	Entries []MapEntry `"{" @@* "}"`
}

type PrimitiveType struct {
	Type string `@("int8"|"int16"|"int32"|"int64"|"uint8"|"uint16"|"uint32"|"uint64"|"float32"|"float64"|"string"|"bool")`
}

type PrimitiveValue struct {
	IntValue    *int64   `@Int`
	FloatValue  *float64 `| @Float`
	StringValue *string  `| @String`
	BoolValue   *bool    `| @("true"|"false")`
}

type Value struct {
	PrimitiveValue *PrimitiveValue `@@`
	ListValue      *ListValue      `| @@`
	MapValue       *MapValue       `| @@`
}

type Type struct {
	PrimitiveType *PrimitiveType `@@`
	ListType      *ListType      `| @@`
	MapType       *MapType       `| @@`
	Identity      *string        `| @Ident`
}

type ObjectField struct {
	Type Type   `@@`
	Name string `@Ident`
}

type Object struct {
	Name   string        `@"object" @Ident`
	Fields []ObjectField `"{" @@* "}"`
}

type EnumValue struct {
	Name  string        `@Ident`
	Value PrimitiveType `"=" @@`
}

type Enum struct {
	Name string        `"enum" @Ident`
	Type PrimitiveType `"for" @@`
	Body []EnumValue   `"{" @@* "}"`
}

type ConstField struct {
	Name  string         `@Ident`
	Value PrimitiveValue `"=" @@`
}

type Const struct {
	Type   PrimitiveType `"const" @@`
	Fields []ConstField  `"{" @@* "}"`
}

type Entry struct {
	Const  *Const  `@@`
	Enum   *Enum   `| @@`
	Object *Object `| @@`
}

type Grammar struct {
	Entries []Entry `@@*`
}
