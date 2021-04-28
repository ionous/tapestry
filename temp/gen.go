package main

import (
	"os"
	"path/filepath"
	"sort"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/pb/internal"
)

var temp = internal.Cap

func main() {
	path := filepath.Join("auto", "all"+temp.Ext)
	if fp, e := os.Create(path); e != nil {
		panic(e)
	} else {
		temp.Header.Must(fp, nil)
		writeSlots(fp)
		writeSlats(fp)
	}
}

var allCmds internal.Cmds = internal.MakeCommands()

func writeSlots(fp *os.File) {
	// write slots
	var sms []internal.SlotMessage
	for _, slots := range iffy.AllSlots {
		for _, slot := range slots {
			var sigs []internal.Sig
			for _, slat := range internal.ImplementorsOf(slot.Type) {
				cmd := allCmds.Add(internal.MakeCommand(slat))
				if cmd != nil {
					sigs = append(sigs, cmd.Signatures()...)
				}
			}
			sort.Slice(sigs, func(i, j int) bool {
				return sigs[i].Raw < sigs[j].Raw
			})
			//
			// pretty.Println(sigs)
			sm := internal.SlotMessage{
				Name: internal.Pascal(composer.SlotName(slot)),
				Desc: internal.ClipDesc(slot.Desc),
				Sigs: sigs,
			}
			sms = append(sms, sm)
		}
	}
	sort.Slice(sms, func(i, j int) bool {
		return sms[i].Name < sms[j].Name
	})
	temp.Slots.Must(fp, sms)
}

func writeSlats(fp *os.File) {
	// write slats
	var sms []internal.SlatMessage
	for _, cmd := range allCmds {
		sm := internal.SlatMessage{cmd}
		sms = append(sms, sm)
	}
	sort.Slice(sms, func(i, j int) bool {
		return sms[i].Name < sms[j].Name
	})
	temp.Slats.Must(fp, sms)
}
