package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
)

type Argument struct {
	Name string // argument name
	From rt.Assignment
}

type Arguments struct {
	Args []*Argument
}

func (*Argument) Compose() composer.Spec {
	return composer.Spec{
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
