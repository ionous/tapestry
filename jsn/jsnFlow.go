package jsn

// Flow wraps str-like values used by the ifspec code generator.
// it alleviates some redundant code generation.
type Flow struct {
	lede, typeName string
}

// MarkFlow ( not MakeFlow ) indicates the start of a set of key-value pairs.
// Unlike the other block types, the block itself is not mutable -- only its values.
func MarkFlow(lede, typeName string) Flow {
	return Flow{lede, typeName}
}

func (n Flow) GetType() string {
	return n.typeName
}
func (n Flow) GetLede() string {
	return n.lede
}
