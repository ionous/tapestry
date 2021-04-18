package parser

var Scanners = []Scanner{
	(*Action)(nil),
	(*AllOf)(nil),
	(*AnyOf)(nil),
	(*Focus)(nil),
	(*Multi)(nil),
	(*Noun)(nil),
	(*Reverse)(nil),
	(*Target)(nil),
	(*Word)(nil),
}

// // Series matches any set of words.
// type Series struct {
// 	Next Rule
// }

// // Number matches any series of digits
// // FUTURE: commas? and words.
// type Number struct {
// }

// // Actor handles addressing actors by name. For example: "name, "
// type Actor struct {
// }
