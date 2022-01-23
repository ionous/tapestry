package blocks

import "strings"

const (
	BKY_COLOUR_HUE            = "%{BKY_COLOUR_HUE}"
	BKY_LISTS_HUE             = "%{BKY_LISTS_HUE}"
	BKY_LOGIC_HUE             = "%{BKY_LOGIC_HUE}"
	BKY_LOOPS_HUE             = "%{BKY_LOOPS_HUE}"
	BKY_MATH_HUE              = "%{BKY_MATH_HUE}"
	BKY_PROCEDURES_HUE        = "%{BKY_PROCEDURES_HUE}"
	BKY_TEXTS_HUE             = "%{BKY_TEXTS_HUE}"
	BKY_VARIABLES_DYNAMIC_HUE = "%{BKY_VARIABLES_DYNAMIC_HUE}"
	BKY_VARIABLES_HUE         = "%{BKY_VARIABLES_HUE}"
)

const (
	FieldCheckbox  = "field_checkbox"
	FieldDropdown  = "field_dropdown"
	FieldInput     = "field_input"
	FieldLabel     = "field_label"
	FieldNumber    = "field_number"
	InputStatement = "input_statement"
	InputValue     = "input_value"
)

// https://developers.google.com/blockly/guides/create-custom-blocks/define-blocks
// other types: field_colour, field_number, field_angle, field_variable, field_label, field_image, input_dummy

func NewLabel(name string) string {
	return Embrace(obj, func(out *Js) {
		out.
			Kv("type", FieldLabel).
			R(comma).
			Kv("text", name)
	})
}

// other options:
// spellcheck: true/false
// text: the default value
func NewText(name string) string {
	return Embrace(obj, func(out *Js) {
		out.
			Kv("type", FieldInput).
			R(comma).
			Kv("name", strings.ToUpper(name))
	})
}

// "value": 100,"min": 0,"max": 100,"precision": 10
func NewNumber(name string) string {
	return Embrace(obj, func(out *Js) {
		out.
			Kv("type", FieldNumber).
			R(comma).
			Kv("name", strings.ToUpper(name))
	})
}

// {"type": "input_value", "name": "DELTA", "check": "Number"}
func NewInput(name, inputType, check string) string {
	return Embrace(obj, func(out *Js) {
		out.
			Kv("type", inputType).
			R(comma).
			Kv("name", strings.ToUpper(name)).
			If(len(check) > 0, func(out *Js) {
				out.
					R(comma).
					Kv("check", strings.ToUpper(name))
			})
	})
}

func NewDropdown(name string, pairs []string) string {
	return Embrace(obj, func(out *Js) {
		out.
			Kv("type", FieldDropdown).
			R(comma).
			Kv("name", strings.ToUpper(name)).
			Kv("option", Embrace(array, func(out *Js) {
				for i, cnt := 0, len(pairs); i < cnt; i += 2 {
					if i > 0 {
						out.R(comma)
					}
					out.Brace(array, func(out *Js) {
						// "first item", "ITEM1"
						out.
							Q(pairs[i]).
							R(comma).
							Q(pairs[i+1])
					})
				}
			}))
	})
}
