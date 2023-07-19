package weave

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func makeCard(amany, bmany bool) (ret string) {
	switch {
	case !amany && !bmany:
		ret = tables.ONE_TO_ONE
	case !amany && bmany:
		ret = tables.ONE_TO_MANY
	case amany && !bmany:
		ret = tables.MANY_TO_ONE
	case amany && bmany:
		ret = tables.MANY_TO_MANY
	}
	return
}

// ugly way to normalize field names
func addField(ctx *Weaver, kindName, fieldName, fieldClass string,
	addField func(k, f, c string) error) (err error) {
	if kind, e := grok.StripArticle(kindName); e != nil {
		err = e
	} else {
		kind, field, class := lang.Normalize(kind), lang.Normalize(fieldName), lang.Normalize(fieldClass)
		err = addField(kind, field, class)
	}
	return
}

func fromTiming(timing mdl.EventTiming) int {
	var part int
	always := timing&mdl.RunAlways != 0
	if always {
		timing ^= mdl.RunAlways
	}
	switch timing {
	case mdl.Before:
		part = 0
	case mdl.During:
		part = 1
	case mdl.After:
		part = 2
	case mdl.Later:
		part = 3
	}
	flags := part + int(rt.FirstPhase)
	if always {
		flags = -flags // marker for rules that need to always run (ex. counters "every third try" )
	}
	return flags
}
