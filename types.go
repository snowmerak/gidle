package main

func IsPrimitiveType(t *Type) bool {
	switch t.PrimitiveType != nil {
	case true:
		return true
	default:
		return false
	}
}

func IsObjectType(t *Type) bool {
	if t.PrimitiveType == nil && t.ListType == nil && t.MapType == nil {
		return true
	}

	return false
}

func IsListType(t *Type) bool {
	switch t.ListType != nil {
	case true:
		return true
	default:
		return false
	}
}

func IsMapType(t *Type) bool {
	switch t.MapType != nil {
	case true:
		return true
	default:
		return false
	}
}
