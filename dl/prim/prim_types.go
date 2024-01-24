// Code generated by Tapestry; edit at your own risk.
package prim

import (
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// bool, a type of str enum.
const Z_Bool_Name = "bool"

const (
	W_Bool_True  = "$TRUE"
	W_Bool_False = "$FALSE"
)

var Z_Bool_T = typeinfo.Str{
	Name: Z_Bool_Name,
	Options: []string{
		W_Bool_True,
		W_Bool_False,
	},
}

// lines, a type of str.
const Z_Lines_Name = "lines"

var Z_Lines_T = typeinfo.Str{
	Name: Z_Lines_Name,
	Markup: map[string]any{
		"comment": "A sequence of characters of any length spanning multiple lines. See also: text.",
	},
}

// text, a type of str.
const Z_Text_Name = "text"

var Z_Text_T = typeinfo.Str{
	Name: Z_Text_Name,
	Markup: map[string]any{
		"comment": "A sequence of characters of any length, all on one line. Examples include letters, words, or short sentences. Text is generally something displayed to the player. See also: lines.",
	},
}

// number, a type of num.
const Z_Number_Name = "number"

var Z_Number_T = typeinfo.Num{
	Name: Z_Number_Name,
}

// package listing of type data
var Z_Types = typeinfo.TypeSet{
	Name: "prim",
	Str:  z_str_list,
	Num:  z_num_list,
}

// a list of all strs in this this package
var z_str_list = []*typeinfo.Str{
	&Z_Bool_T,
	&Z_Lines_T,
	&Z_Text_T,
}

// a list of all nums in this this package
var z_num_list = []*typeinfo.Num{
	&Z_Number_T,
}