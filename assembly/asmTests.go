package assembly

import (
	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/check"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/tables"
)

func AssembleTests(asm *Assembler) (err error) {
	// todo: doesn't build mdl check
	// doesnt check for conflicts or errors in test definitions
	var name, expect string
	var prog []byte
	type Entry struct {
		Name string
		Cmd  *check.CheckOutput
	}
	var list []Entry
	var curr *check.CheckOutput
	if e := tables.QueryAll(asm.cache.DB(),
		`select name, prog, expect
		from asm_check ek
		join asm_expect ex
			using (name)
		join eph_prog ep
			on (ek.idProg = ep.rowid)
		order by name, progType, idProg`,
		func() (err error) {
			if curr == nil || curr.Name != name {
				curr = &check.CheckOutput{
					Name:   name,
					Expect: expect,
				}
				list = append(list, Entry{name, curr})
			}
			var el core.Activity
			if e := cin.Decode(&el, prog, iffy.AllSignatures); e != nil {
				err = e
			} else {
				curr.Test.Exe = append(curr.Test.Exe, el.Exe...)
			}
			return
		}, &name, &prog, &expect); e != nil {
		err = e
	} else {
		for _, c := range list {
			if e := asm.WriteProgram(c.Name, "CheckOutput", c.Cmd); e != nil {
				err = e
				break
			}
		}
	}
	return
}
