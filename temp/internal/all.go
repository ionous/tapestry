package internal

import (
	"git.sr.ht/~ionous/iffy"
)

func MakeCommands() Cmds {
	cmds := make(Cmds)
	for _, slats := range iffy.AllSlats {
		for _, slat := range slats {
			cmd := MakeCommand(slat)
			if cmd == nil {
				println("couldnt make", slat.Compose().Name)
			} else {
				cmds.Add(cmd)
			}
		}
	}
	return cmds
}
