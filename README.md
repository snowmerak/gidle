# gidle

gidle is a simple IDL (Interface Definition Language) for JSON.

## Supported Languages

1. Go: `gidle -l go`
2. TypeScript: `gidle -l ts`
3. Dart: `gidle -l dart`
4. Rust: `gidle -l rs`
5. C#: `gidle -l cs`

> Rust version depends on [serde_json](https://docs.rs/serde_json/latest/serde_json/)
> and [serde](https://docs.rs/serde/latest/serde/) crate.
> You must add them to your `Cargo.toml` file.
> e.g.
> serde = { version = <VERSION>, features = ["derive"] }
> serde_json = <VERSION>

## Usage

### Install

```bash
go install github.com/snowmerak/gidle@latest
```

### Generate

```bash
gidle -i <gidle-file> -o <output-file> -l <language>
```

e.g.

```
gidle -i test.gidle -o out/test.go -l go
gidle -i test.gidle -o out/test.ts -l ts
gidle -i test.gidle -o out/test.dart -l dart
gidle -i test.gidle -o out/test.rs -l rs
gidle -i test.gidle -o out/test.cs -l cs
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
14. `map <Key-Type> for <Value-Type>`: map of key-type for value-type
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

```dart
import 'dart:convert';

class Person {
	String? name;
	int? age;
	List<String>? friends;
	Map<String, String>? properties;

	Person({
		this.name,
		this.age,
		this.friends,
		this.properties,
	});

	toMap() {
		return {
			"name": name,
			"age": age,
			"friends": friends,
			"properties": properties,
		};
	}

	String toJson() {
		return jsonEncode(toMap());
	}

	Person.fromMap(Map<String, dynamic> map) {
		name = map["name"];
		age = map["age"];
		friends = List<String>.from(map["friends"]);
		properties = Map<String, String>.from(map["properties"]);
	}

	Person.fromJson(String source) : this.fromMap(jsonDecode(source));

}

enum CASE {
	UPPER,
	LOWER,
}

const BOUNDARY_MAX = 100;
const BOUNDARY_MIN = 0;
```

```typescript
export namespace gidle {
  export namespace test {
    export namespace main {
      export interface Person {
        name: string;
        age: number;
        friends: Array<string>;
        properties: Map<string, string>;
      }

      export enum CASE {
        UPPER = 0,
        LOWER = 1,
      }

      export function indexOfCASE(value: CASE): number {
        switch (value) {
          case CASE.UPPER:
            return 0;
          case CASE.LOWER:
            return 1;
          default:
            return -1;
        }
      }

      export function getCASE(index: number): CASE {
        switch (index) {
          case 0:
            return CASE.UPPER;
          case 1:
            return CASE.LOWER;
          default:
            throw new Error("unknown enum value");
        }
      }

      export const BOUNDARY_MAX = 100;
      export const BOUNDARY_MIN = 0;
    }
  }
}
```

```rust
use std::collections::HashMap;
use serde::{Deserialize, Serialize};
use serde_json::{to_string, from_str, Result};

#[derive(Debug, Serialize, Deserialize)]
pub struct Person {
	pub name: String,
	pub age: i32,
	pub friends: Vec<String>,
	pub properties: HashMap<String, String>,
}

impl Person {
	pub fn new(name: String, age: i32, friends: Vec<String>, properties: HashMap<String, String>) -> Self {
		Self {
			name,
			age,
			friends,
			properties,
		}
	}
	pub fn to_json(&self) -> Result<String> {
		to_string(self)
	}
	pub fn from_json(json: &str) -> Result<Self> {
		from_str(json)
	}
}

pub enum CASE {
	UPPER = 0,
	LOWER = 1,
}

pub const BOUNDARY_MAX: i32 = 100;
pub const BOUNDARY_MIN: i32 = 0;
```

```csharp
using System.Text.Json;
using System.Text.Json.Serialization;

namespace gidle {
    namespace test {
        namespace main {
            public class Person {
                [JsonPropertyName("name")]
                public string Name { get; set; }
                [JsonPropertyName("age")]
                public int Age { get; set; }
                [JsonPropertyName("friends")]
                public List<string> Friends { get; set; }
                [JsonPropertyName("properties")]
                public Dictionary<string, string> Properties { get; set; }

                public Person(string name, int age, List<string> friends, Dictionary<string, string> properties) {
                    this.Name = name;
                    this.Age = age;
                    this.Friends = friends;
                    this.Properties = properties;
                }

                public static Person? FromJson(string json) {
                    return JsonSerializer.Deserialize<Person>(json);
                }

                public string ToJson() {
                    return JsonSerializer.Serialize(this);
                }

            }
            public class CASE {
                public const byte UPPER = 0;
                public const byte LOWER = 1;
                public static int IndexOf(byte value) {
                    return value switch {
                        UPPER => 0,
                        LOWER => 1,
                        _ => throw new ArgumentOutOfRangeException(nameof(value), "Invalid value")
                    };
                }
                public static byte ValueOf(int index) {
                    return index switch {
                        0 => UPPER,
                        1 => LOWER,
                        _ => throw new ArgumentOutOfRangeException(nameof(index), "Invalid index")
                    };
                }
            }
            public static class BOUNDARY {
                public const int MAX = 100;
                public const int MIN = 0;
            }
        }
    }
}
```