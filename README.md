# gidle

gidle is a simple IDL (Interface Definition Language) for JSON and YAML.

## Installation

```bash
go install github.com/snowmerak/gidle/cmd/gidle
```

## IDL

```
type <message-name> {
    <field-type> <field-name> [ = <default-value> ];
    <field-type> <field-name> [ = <default-value> ];
    <field-type> <field-name> [ = <default-value> ];
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
type Person {
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
