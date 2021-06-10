// Code generated by "makeops"; edit at your own risk.
package value

import (
	"encoding/json"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/reader"
)

// Bool requires a user-specified string.
type Bool struct {
	Str string
}

func (op *Bool) String() (ret string) {
	return op.Str
}

func (op *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "bool",
		"value": op.Str,
	})
}

const Bool_True = "$TRUE"
const Bool_False = "$FALSE"

func (*Bool) Compose() composer.Spec {
	return composer.Spec{
		Name: "bool",
		Uses: "str",
		Choices: []string{
			Bool_True, Bool_False,
		},
		Strings: []string{
			"true", "false",
		},
	}
}

// Lines requires a user-specified string.
type Lines struct {
	Str string
}

func (op *Lines) String() (ret string) {
	return op.Str
}

func (op *Lines) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "lines",
		"value": op.Str,
	})
}

func (*Lines) Compose() composer.Spec {
	return composer.Spec{
		Name:        "lines",
		Uses:        "str",
		OpenStrings: true,
	}
}

// Number requires a user-specified number.
type Number struct {
	Value float64
}

func (op *Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "number",
		"value": op.Value,
	})
}

func (*Number) Compose() composer.Spec {
	return composer.Spec{
		Name: "number",
		Uses: "num",
	}
}

// PatternName requires a user-specified string.
type PatternName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *PatternName) String() (ret string) {
	return op.Str
}

func (op *PatternName) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{"id": op.At.Offset,
		"type":  "pattern_name",
		"value": op.Str,
	})
}

func (*PatternName) Compose() composer.Spec {
	return composer.Spec{
		Name:        "pattern_name",
		Uses:        "str",
		OpenStrings: true,
	}
}

// RelationName requires a user-specified string.
type RelationName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *RelationName) String() (ret string) {
	return op.Str
}

func (op *RelationName) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{"id": op.At.Offset,
		"type":  "relation_name",
		"value": op.Str,
	})
}

func (*RelationName) Compose() composer.Spec {
	return composer.Spec{
		Name:        "relation_name",
		Uses:        "str",
		OpenStrings: true,
	}
}

// Text requires a user-specified string.
type Text struct {
	Str string
}

func (op *Text) String() (ret string) {
	return op.Str
}

func (op *Text) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "text",
		"value": op.Str,
	})
}

func (*Text) Compose() composer.Spec {
	return composer.Spec{
		Name:        "text",
		Uses:        "str",
		OpenStrings: true,
	}
}

// VariableName requires a user-specified string.
type VariableName struct {
	At  reader.Position `if:"internal"`
	Str string
}

func (op *VariableName) String() (ret string) {
	return op.Str
}

func (op *VariableName) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{"id": op.At.Offset,
		"type":  "variable_name",
		"value": op.Str,
	})
}

func (*VariableName) Compose() composer.Spec {
	return composer.Spec{
		Name:        "variable_name",
		Uses:        "str",
		OpenStrings: true,
	}
}

var Slats = []composer.Composer{
	(*Bool)(nil),
	(*Lines)(nil),
	(*Number)(nil),
	(*PatternName)(nil),
	(*RelationName)(nil),
	(*Text)(nil),
	(*VariableName)(nil),
}
