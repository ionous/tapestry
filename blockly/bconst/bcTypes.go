package bconst

// return the passed type name formatted for use as a stacked block
func StackedName(name string) string {
	return "_" + name + "_stack"
}

func MutatorName(name string) string {
	return "_" + name + "_mutator"
}
