package rt

// wrap an existing value as an assigned value
func AssignValue(v Value) Assignment {
	return simpleAssignment{v}
}

type simpleAssignment struct{ v Value }

func (a simpleAssignment) GetAssignedValue(Runtime) (ret Value, _ error) {
	ret = a.v
	return
}
