package rt

// uniform access to objects and variables.
type Address interface {
	GetReference(Runtime) (Reference, error)
}

// a value that might contain other values.
type Reference interface {
	// set the value at the current cursor.
	SetValue(newValue Value) (err error)
	// read the value of the current cursor.
	GetValue() (ret Value, err error)
	// peeks into the current cursor at the requested field or index.
	Dot(next Dotted) (ret Reference, err error)
}

// a "pointer" into the contents of some as yet unspecified target:
// for example, a field for use with objects or records,
// or an index for use in lists.
type Dotted interface {
	// read a value from the contents of the passed target.
	Peek(target Cursor) (Cursor, error)
	// write into the contents of the passed target.
	Poke(target Cursor, newValue Value) error
}

// a value within an object, record, or list.
type Cursor interface {
	CurrentValue() Value
	SetAtIndex(int, Value) error
	GetAtIndex(int) (Cursor, error)
	SetAtField(string, Value) error
	GetAtField(string) (Cursor, error)
}
