package cout

// TextMarshaler exposes custom serialization for text values.
// We force user code to manually cast and call this interface --
// that avoids having to check every value for the special text value
// and allows the user code to customize fallback handling.
type TextMarshaler interface {
	TextValue(*string) bool
}

type VariableMarshaler interface {
	VariableValue(*string) bool
}

func (cv Chart) TextValue(pstr *string) (okay bool) {
	str := *pstr
	if len(str) > 0 && str[0] == '@' {
		cv.StrValue(compactString("@" + str))
		okay = true
	}
	return
}

func (cv Chart) VariableValue(pstr *string) (okay bool) {
	str := *pstr
	// a leading ampersand would with @@ escaped text serialization.
	if leadingAmp := len(str) > 0 && str[0] == '@'; !leadingAmp {
		cv.StrValue(compactString("@" + str))
		okay = true
	}
	return
}

type compactString string

func (n compactString) GetType() string {
	return "unexpected type"
}

func (n compactString) SetStr(v string) {
	panic("compact output shouldnt be setting a string")
}

func (n compactString) GetStr() string {
	return string(n)
}
