package debug

import (
	"strings"

	"git.sr.ht/~ionous/iffy/dl/composer"
)

type Level int

//go:generate stringer -type=Level
const (
	Note Level = iota + 1
	ToDo
	Fix
	Warning
	Error
)

func (*Level) Compose() composer.Spec {
	return composer.Spec{
		Name:    "debug_level",
		Group:   "list",
		Strings: []string{"note", "to do", "warning", "fix"},
		Desc:    "Debug level.",
		Stub:    true, // the stub parse the flag
	}
}

func (lvl *Level) Header() string {
	return strings.Repeat("#", 1+int(*lvl)) + " " + lvl.String() + ":"
}
