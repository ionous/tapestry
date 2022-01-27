package spec

// return english human friendly text.
func (op *TermSpec) Label() string {
	return op.Key
}

// return a programmatic name, unique in the flow: usually no spaces, etc.
func (op *TermSpec) Field() (ret string) {
	if len(op.Name) > 0 {
		ret = op.Name
	} else {
		ret = op.Key
	}
	return
}

// return the name of the spec which describes the content of this term.
func (op *TermSpec) TypeName() (ret string) {
	if len(op.Type) > 0 {
		ret = op.Type
	} else {
		ret = op.Field()
	}
	return
}
