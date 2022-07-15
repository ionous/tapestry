package bconst

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
