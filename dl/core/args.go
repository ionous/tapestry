package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
)

// Argument - a name and value passed to a pattern or other indeterminate method.
type Argument struct {
	Name string
	From rt.Assignment
}

// Arguments - composer wrapper for a list of Argument(s).
type Arguments struct {
	Args []*Argument
}

func (*Argument) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "arg",
		Spec:  " {name:variable_name}: {from:assignment}",
		Group: "patterns",
		Stub:  true,
	}
}

func (*Arguments) Compose() composer.Spec {
	return composer.Spec{
		Spec:  " {arguments%args+argument|commas}",
		Group: "patterns",
		Stub:  true,
	}
}
