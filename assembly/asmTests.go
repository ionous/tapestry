package assembly

import (
	"git.sr.ht/~ionous/iffy/dl/check"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/tables"
)

func AssembleTests(asm *Assembler) (err error) {
	// todo: doesn't build mdl check
	// doesnt check for conflicts or errors in test definitions
	var name, expect string
	var prog []byte
	list := make(map[string]interface{})
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
					Test:   &core.Activity{},
				}
				list[name] = curr
			}
			var el core.Activity
			if e := tables.DecodeGob(prog, &el); e != nil {
				err = e
			} else {
				curr.Test.Exe = append(curr.Test.Exe, &el)
			}
			return
		}, &name, &prog, &expect); e != nil {
		err = e
	} else {
		err = asm.WriteGobs(list)
	}
	return
}
