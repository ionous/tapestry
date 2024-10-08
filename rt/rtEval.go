package rt

// Execute runs a bit of code that has no return value.
type Execute interface {
	Execute(Runtime) error
}

// Assignment accesses values in a generic way.
type Assignment interface {
	GetAssignedValue(Runtime) (Value, error)
}

// BoolEval represents the result of some true or false expression.
type BoolEval interface {
	GetBool(Runtime) (Value, error)
}

// NumEval represents the result of some numeric expression.
type NumEval interface {
	GetNum(Runtime) (Value, error)
}

// TextEval represents the result of some expression which creates a string.
type TextEval interface {
	GetText(Runtime) (Value, error)
}

// RecordEval yields access to a set of fields and their values.
type RecordEval interface {
	GetRecord(Runtime) (Value, error)
}

// NumListEval represents the computation of a series of numeric values.
type NumListEval interface {
	GetNumList(Runtime) (Value, error)
}

// TextListEval represents the computation of a series of strings.
type TextListEval interface {
	GetTextList(Runtime) (Value, error)
}

// RecordListEval represents the computation of a series of a set of fields.
type RecordListEval interface {
	GetRecordList(Runtime) (Value, error)
}
