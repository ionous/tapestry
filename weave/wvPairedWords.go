package weave

import (
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (cat *Catalog) AssertPlural(opSingular, opPlural string) error {
	return cat.Schedule(assert.PluralPhase, func(ctx *Weaver) (err error) {
		if plural, ok := UniformString(opPlural); !ok {
			err = InvalidString(opPlural)
		} else if singular, ok := UniformString(opSingular); !ok {
			err = InvalidString(opSingular)
		} else {
			err = cat.writer.Plural(ctx.d.name, plural, singular, cat.cursor)
		}
		return
	})
}

func (cat *Catalog) AssertOpposite(opOpposite, opWord string) error {
	return cat.Schedule(assert.PluralPhase, func(ctx *Weaver) (err error) {
		if a, ok := UniformString(opOpposite); !ok {
			err = InvalidString(opOpposite)
		} else if b, ok := UniformString(opWord); !ok {
			err = InvalidString(opWord)
		} else {
			err = cat.writer.Opposite(ctx.d.name, a, b, cat.cursor)
		}
		return
	})
}
