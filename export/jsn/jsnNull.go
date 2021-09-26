package jsn

// NullMarshaler implements Marshal and does nothing.
type NullMarshaler struct{}

func (NullMarshaler) MapValues(lede, kind string)               {}
func (NullMarshaler) MapKey(sig, field string)                  {}
func (NullMarshaler) MapLiteral(field string)                   {}
func (NullMarshaler) SpecifyValue(kind string, val interface{}) {}
func (NullMarshaler) SpecifyEnum(kind string, val Enumeration)  {}
func (NullMarshaler) PickValues(kind, choice string)            {}
func (NullMarshaler) RepeatValues(hint int)                     {}
func (NullMarshaler) EndValues()                                {}
func (NullMarshaler) SetCursor(id string)                       {}
func (NullMarshaler) Warning(err error)                         {}
func (NullMarshaler) Error(err error)                           {}
