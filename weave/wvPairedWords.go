package weave

import (
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (cat *Catalog) AssertPlural(opSingular, opPlural string) error {
	return cat.Schedule(assert.PluralPhase, func(ctx *Weaver) (err error) {
		if plural, ok := UniformString(opPlural); !ok {
			err = InvalidString(opPlural)
		} else if singular, ok := UniformString(opSingular); !ok {
			err = InvalidString(opSingular)
		} else if rows, e := cat.db.Query(
			`select one
				from active_plurals
				where many = ?1`, plural); e != nil {
			err = e
		} else {
			var prev string
			var dupe int // log duplicates?
			if e := tables.ScanAll(rows, func() (err error) {
				if prev == singular {
					dupe++
				} else {
					err = errutil.Fmt("conflict: plural %q had singular %q wants %q", plural, prev, opSingular)
				}
				return
			}, &prev); e != nil {
				err = e
			} else if dupe == 0 {
				err = cat.writer.Plural(ctx.d.name, plural, singular, cat.cursor)
			}
		}
		return // from schedule
	})
}

func (cat *Catalog) AssertOpposite(opOpposite, opWord string) error {
	return cat.Schedule(assert.PluralPhase, func(ctx *Weaver) (err error) {
		if a, ok := UniformString(opOpposite); !ok {
			err = InvalidString(opOpposite)
		} else if b, ok := UniformString(opWord); !ok {
			err = InvalidString(opWord)
		} else if rows, e := cat.db.Query(
			`select oneWord, otherWord
				from active_rev
				where oneWord=?1 or otherWord=?2`, a, b); e != nil {
			err = e
		} else {
			var x, y string
			var dupe int // log duplicates?
			if e := tables.ScanAll(rows, func() (err error) {
				if (x == a && y == b) || (y == b && x == a) {
					dupe++
				} else {
					err = errutil.Fmt("conflict: %q had opposite %q wanted %q as %q ", x, y, a, b)
				}
				return
			}, &x, &y); e != nil {
				err = e
			} else if dupe == 0 {
				if e := cat.writer.Opposite(ctx.d.name, a, b, cat.cursor); e != nil {
					err = e
				} else if a != b {
					// write the opposite pairing to; helps simplify queries.
					if e := cat.writer.Opposite(ctx.d.name, b, a, cat.cursor); e != nil {
						err = e
					}
				}
			}
		}
		return // from schedule
	})
}
