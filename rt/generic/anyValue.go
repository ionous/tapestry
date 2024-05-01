package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
)

// Value represents any one of Tapestry's built in types.
type Value interface {
	// identifies the general category of the value.
	Affinity() affine.Affinity
	// identifies how the passed value is represented internally.
	// ex. a number can be represented as float or as an int,
	// a record might be one of several different kinds,
	// text might represent an object id of a specific kind, an aspect, trait, or other value.
	Type() string
	// return this value as a bool, or panic if the value isn't a bool.
	Bool() bool
	// return this value as a float, or panic if the value isn't a number.
	Float() float64
	// return this value as an int, or panic if the value isn't a number.
	Int() int
	// return this value as a string; it doesn't panic.
	// for non-string values, similar to package reflect, it returns a string of the form "<type value>".
	String() string
	// return this value as a record, or panic if the value isn't a record.
	// warning: can return a nil.
	Record() (*Record, bool)
	// return this value as a slice of floats, or panic if this isn't a float slice.
	// note: while primitive values support both ints and floats, slices can only be floats.
	Floats() []float64
	// return this value as a slice of strings, or panic if this isn't a string slice.
	Strings() []string
	// return this value as a slice of records, or panic if not a record slice.
	// note: every value in the returned slice is expected to be record of this value's Type().
	Records() []*Record
	// return a value representing a field inside this record.
	// if the field holds a record or a slice, the returned value shares its memory with the named field.
	// errors if the field doesn't exist.
	// panics if this isn't a record.
	FieldByName(string) (Value, error)
	// writes a *copy* of the passed value into a record.
	// errors if the field doesn't exist or if its affinity cant support the passed value.
	// panics if this isn't a record.
	SetFieldByName(string, Value) error
	// the number of elements in the value.
	// panics if this cant have a length: isnt a slice or a string.
	Len() int
	// return the nth element of this value, where 0 is the first value.
	// panics if this isn't a slice.
	Index(int) Value
	// writes a *copy* of the passed value into a slice
	// panics if this isn't a slice, if the index is out of range, or if the affinities are mismatched.
	// errors if the types are mismatched
	SetIndex(int, Value) error
	// adds a *copy* of a value, or a copy of a slice of values, to the end of this slice.
	// panics if this value isn't an appropriate kind of slice; errors on subtype.
	// In golang, this is a package level function, presumably to mirror the built-in append()
	Appends(Value) error
	// return a *copy* of this slice and its values containing the first index up to (and not including) the second index.
	// panics if this value isn't a slice.
	Slice(i, j int) (Value, error)
	// cut elements out of this slice from start to end,
	// adding copies of the passed additional elements (if any) at the start of the cut point.
	// Returns the cut elements, or an error if the start and end indices are bad.
	// panics if this value isn't a slice, or if additional element(s) are of an incompatible affinity.
	Splice(start, end int, add Value) (Value, error)
}
