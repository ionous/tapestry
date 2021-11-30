package composer

func SpecName(c Composer) (ret string) {
	if c == nil {
		ret = "<nil>"
	} else if spec := c.Compose(); len(spec.Name) > 0 {
		ret = spec.Name
	} else {
		ret = "???" // all generated types have names now.
	}
	return
}

// translate a choice, typically a $TOKEN, to a value.
func FindChoice(op Composer, choice string) (ret string, found bool) {
	spec := op.Compose()
	if s, i := spec.IndexOfChoice(choice); i >= 0 {
		ret = s
		found = true
	} else if spec.OpenStrings {
		ret = choice
	}
	return
}
