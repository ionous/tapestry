package rt

import "github.com/ionous/iffy/dl/composer"

var Slots = []composer.Slot{{
	Name: "execute",
	Type: (*Execute)(nil),
	Desc: "Execute: Run a series of statements.",
}, {
	Name: "bool_eval",
	Type: (*BoolEval)(nil),
	Desc: "Booleans: Statements which return true/false values.",
}, {
	Name: "number_eval",
	Type: (*NumberEval)(nil),
	Desc: "Numbers: Statements which return a number.",
}, {
	Name: "text_eval",
	Type: (*TextEval)(nil),
	Desc: "Texts: Statements which return text.",
}, {
	Name: "num_list_eval",
	Type: (*NumListEval)(nil),
	Desc: "Number List: Statements which return a list of numbers.",
}, {
	Name: "text_list_eval",
	Type: (*TextListEval)(nil),
	Desc: "Text Lists: Statements which return a list of text.",
}}