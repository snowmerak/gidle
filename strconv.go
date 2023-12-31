package main

import "strings"

func SnakeToPascal(s string) string {
	sb := strings.Builder{}
	sb.Grow(len(s))

	for i := range s {
		if i == 0 && s[i] >= 'a' && s[i] <= 'z' {
			sb.WriteByte(s[i] - 32)
		} else if s[i] == '_' && i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z' {
			sb.WriteByte(s[i+1] - 32)
			i++
		} else {
			sb.WriteByte(s[i])
		}
	}

	return sb.String()
}

func SnakeToCamel(s string) string {
	sb := strings.Builder{}
	sb.Grow(len(s))

	for i := range s {
		if s[i] == '_' && i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z' {
			sb.WriteByte(s[i+1] - 32)
			i++
		} else {
			sb.WriteByte(s[i])
		}
	}

	return sb.String()
}
