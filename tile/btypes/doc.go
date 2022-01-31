package btypes

// every workspace block that needs a mutation uses
// the "tapestry_generic_mutation" mutation,
// the "tapestry_generic_mixin" mixi
// and "tapestry_generic_extension" extension.
//
// they use some "customData" stored in the workspace blocks's jsondef,
// and rely on a mutation user interface ( mui ) block: one mui block per workspace block.
//
// the custom data generates all of the fields for the workspace block ( except for the block's label )
// and it contains information on which fields are optional and repeating to direct the mutation.
//
// the mui block is a relatively normal block named "_<type>_mutator".
// unlike blockly's mutations: there aren't smaller blocks to drag and drop.
// instead, the mui contains checkboxes and number fields ( sliders ) to control the mutation.
//
// the fields of the workspace block are classified into a few different types:
//
//  	- optional, non-repeating;
//  	- optional, stackable repeats:
//  			the mui contains a checkbox;
//  			custom data contains whatever field it would have had it been a required field.
//
//  	- optional, non-stacking repeats
//  			the mui contains a number; the minimum value is 0
//  			custom data contains whatever input needs to repeat.
//
//  	- required non-stacking repeats:
//  			same, except the minimum value is 1.
//
// the mui controls ( the checkboxes and numbers ) and the are embedded in dummy inputs inside the mui.
// the dummy inputs have the same name as the corresponding inputs in the workspace block.
