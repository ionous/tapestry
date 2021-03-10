package assembly

import (
	"reflect"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// reads mdl_kind, mdl_field, mdl_noun
// out of asm_noun and asm_value... via eph_value... via NewValue()
func assembleInitialFields(asm *Assembler) (err error) {
	var store valueStore
	var curr, last valueInfo
	if e := tables.QueryAll(asm.cache.DB(),
		`select asm.noun, mf.field, mf.type, asm.value
			from asm_noun as asm
		join mdl_field mf
			on (asm.prop = mf.field)
			and (mf.type != 'aspect')
		where instr((
			select mk.kind || ',' || mk.path || ','
			from mdl_kind mk 
			join mdl_noun mn
			using (kind)
			where (mn.noun = asm.noun)
		), mf.kind || ',')
		order by noun, field, type`,
		func() (err error) {
			if nv, e := convertField(curr.fieldType, curr.value); e != nil {
				err = errutil.New(e, "assembling", curr.target, curr.field)
			} else if last.target != curr.target || last.field != curr.field {
				store.add(last)
				last, last.value = curr, nv
			} else if !reflect.DeepEqual(last.value, nv) {
				e := errutil.Fmt("conflicting values: %s != %s:%T", last.String(), curr.String(), nv)
				err = errutil.New(e, "assembling", curr.target, curr.field, e)
			}
			return
		},
		&curr.target, &curr.field, &curr.fieldType,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.writeFieldValues(asm)
	}
	return
}
