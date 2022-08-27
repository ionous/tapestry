package bconst

// defaults for colors are in blockly-master\msg\messages.js
const (
	COLOUR_HUE            = "%{BKY_COLOUR_HUE}"
	LISTS_HUE             = "%{BKY_LISTS_HUE}"
	LOGIC_HUE             = "%{BKY_LOGIC_HUE}"
	LOOPS_HUE             = "%{BKY_LOOPS_HUE}"
	MATH_HUE              = "%{BKY_MATH_HUE}"
	PROCEDURES_HUE        = "%{BKY_PROCEDURES_HUE}"
	TEXTS_HUE             = "%{BKY_TEXTS_HUE}"
	VARIABLES_DYNAMIC_HUE = "%{BKY_VARIABLES_DYNAMIC_HUE}"
	VARIABLES_HUE         = "%{BKY_VARIABLES_HUE}"
)

// "BKY_DEFAULT_COLOR"
type Colour int

//go:generate stringer -type=Colour -linecomment
const (
	DefaultColor Colour = iota // %{BKY_TAPCOLOR_DEFAULT}
)

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
)
