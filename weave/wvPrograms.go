package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
func (cat *Catalog) AssertGrammar(opName string, prog *grammar.Directive) error {
	return cat.Schedule(assert.RequireRules /*GrammarPhase*/, func(ctx *Weaver) error {
		d, at := ctx.d, ctx.at
		return cat.writer.Grammar(d.name, opName, prog, at)
	})
}
