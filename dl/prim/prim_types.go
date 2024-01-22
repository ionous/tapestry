// Code generated by Tapestry; edit at your own risk.
package prim

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// bool, a type of str enum.
const Z_Bool_Type = "bool"

const (
	W_Bool_True  = "$TRUE"
	W_Bool_False = "$FALSE"
)

var Z_Bool_Info = typeinfo.Str{
	Name: Z_Bool_Type,
	Options: []string{
		W_Bool_True,
		W_Bool_False,
	},
}

// lines, a type of str.
const Z_Lines_Type = "lines"

var Z_Lines_Info = typeinfo.Str{
	Name: Z_Lines_Type,
	Markup: map[string]any{
		"comment": "A sequence of characters of any length spanning multiple lines. See also: text.",
	},
}

// text, a type of str.
const Z_Text_Type = "text"

var Z_Text_Info = typeinfo.Str{
	Name: Z_Text_Type,
	Markup: map[string]any{
		"comment": "A sequence of characters of any length, all on one line. Examples include letters, words, or short sentences. Text is generally something displayed to the player. See also: lines.",
	},
}

// number, a type of num.
const Z_Number_Type = "number"

var Z_Number_Info = typeinfo.Num{
	Name: Z_Number_Type,
}

// a list of all command signatures
// ( for processing and verifying story files )
var Z_Signatures = map[uint64]interface{}{}
