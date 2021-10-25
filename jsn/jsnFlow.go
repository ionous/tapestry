package jsn

// Flow wraps str-like values used by the ifspec code generator.
type Flow struct {
	lede, typeName string
	op             interface{}
}

// MakeFlow indicates the start of a set of key-value pairs.
// Unlike the other block types, the block itself is not mutable -- only its values.
func MakeFlow(lede, typeName string, op interface{}) Flow {
	return Flow{lede, typeName, op}
}

func (n Flow) GetType() string {
	return n.typeName
}
func (n Flow) GetLede() string {
	return n.lede
}
func (n Flow) GetValue() interface{} {
	return n.op
}
