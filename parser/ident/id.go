package ident

import (
	"git.sr.ht/~ionous/iffy/lang"
)

type Id string

func (id Id) String() (ret string) {
	return string(id)
}

// creates a new string id from the passed raw string.
func IdOf(str string) (ret Id) {
	return Id(nameOf(str))
}

func nameOf(str string) (ret string) {
	return lang.Underscore(str)
}
