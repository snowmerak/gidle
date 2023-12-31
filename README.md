# gidle

gidle is a simple IDL (Interface Definition Language) for JSON and YAML.

## Installation

```bash
go install github.com/snowmerak/gidle
```

## IDL

### Syntax

```
const <const-name> for <primitive-type-name> {
    <field-name> = <field-value>
    <field-name> = <field-value>
    <field-name> = <field-value>
}

enum <enum-name> for <primitive-type-name> {
    <field-name> = <field-value>
    <field-name> = <field-value>
    <field-name> = <field-value>
}

object <object-name> {
    <field-type> <field-name>
    <field-type> <field-name>
    <field-type> <field-name>
}
```

### Types

1. `int8`: 8-bit signed integer
2. `int16`: 16-bit signed integer
3. `int32`: 32-bit signed integer
4. `int64`: 64-bit signed integer
5. `uint8`: 8-bit unsigned integer
6. `uint16`: 16-bit unsigned integer
7. `uint32`: 32-bit unsigned integer
8. `uint64`: 64-bit unsigned integer
9. `float32`: 32-bit floating point number
10. `float64`: 64-bit floating point number
11. `string`: UTF-8 string
12. `bool`: boolean value
13. `list of <Type>`: array of type
14. `map <Key-Type> for <Value-Type>`: map of key-type to value-type
15. `<message-name>`: message type

### Example

```
package gidle.test.main

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
```

```go
package main

type Person struct {
	Name       string            `json:"name"`
	Age        int32             `json:"age"`
	Friends    []string          `json:"friends"`
	Properties map[string]string `json:"properties"`
}

type CASE uint8

const (
	CASE_UPPER = CASE(0)
	CASE_LOWER = CASE(1)
)

func (e CASE) String() string {
	switch e {
	case CASE_UPPER:
		return "UPPER"
	case CASE_LOWER:
		return "LOWER"
	default:
		return "unknown enum value"
	}
}

func GetCASE(index int) (value CASE, ok bool) {
	switch index {
	case 0:
		return CASE_UPPER, true
	case 1:
		return CASE_LOWER, true
	default:
		return value, false
	}
}

func IndexOfCASE(value CASE) (index int, ok bool) {
	switch value {
	case CASE_UPPER:
		return 0, true
	case CASE_LOWER:
		return 1, true
	default:
		return index, false
	}
}

const (
	BOUNDARY_MAX = 100
	BOUNDARY_MIN = 0
)
```