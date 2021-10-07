package cout

// TextMarshaler exposes custom serialization for text values.
// We force user code to manually cast and call this interface --
// that avoids having to check every value for the special text value
// and allows the user code to customize fallback handling.
type TextMarshaler interface {
	TextValue(string, *string) bool
}

type VariableMarshaler interface {
	VariableValue(string, *string) bool
}

func (cv Chart) TextValue(typeName string, pstr *string) (okay bool) {
	str := *pstr
	if len(str) > 0 && str[0] == '@' {
		cv.GenericValue(typeName, "@"+str)
		okay = true
	}
	return
}

func (cv Chart) VariableValue(typeName string, pstr *string) (okay bool) {
	str := *pstr
	// a leading ampersand would with @@ escaped text serialization.
	if leadingAmp := len(str) > 0 && str[0] == '@'; !leadingAmp {
		cv.GenericValue(typeName, "@"+str)
		okay = true
	}
	return
}
