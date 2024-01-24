package bconst

import "strings"

// yup.
const DefaultColor = "TAP_HUE"

// tapestry typespec markup
const ColorMarkup = "blockly-color"
const StackMarkup = "blockly-stack"

const RootBlockMarkup = "mosaic-root"     // special shapes that have no output.
const InlineBlockMarkup = "mosaic-inline" // a block with no labels

const (
	FieldCheckbox = "field_checkbox"
	FieldDropdown = "field_dropdown"
	FieldText     = "field_input"
	FieldLabel    = "field_label"
	FieldNumber   = "field_number"
	// blocks can stack vertically ( or link horizontally )
	InputStatement = "input_statement"
	// blocks can link horizontally ( or stack vertically )
	InputValue = "input_value"
	// fields are partitioned into rows of inputs;
	// if a row doesnt need connect with other blocks, it can use a dummy input.
	InputDummy = "input_dummy"
)

// https://developers.google.com/blockly/guides/create-custom-blocks/define-blocks
// other types: field_colour, field_number, field_angle, field_variable, field_label, field_image, input_dummy

const (
	// alternative class for str fields so that open strs with choices can use dropdowns.
	MosaicStrField = "mosaic_str_field"
	// alternative class for text fields so all text can have placeholder text.
	MosaicTextField = "mosaic_text_field"
	// multiline variant of the text field
	MosaicMultilineField = "mosaic_multiline_field"
)

// underscore name to something else.
func KeyName(s string) string {
	return "$" + strings.ToUpper(s)
}
