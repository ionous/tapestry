package weave

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// StoryStatement - import a single story statement.
// used during weave, and expects that the runtime is the importer's own runtime.
// ( as opposed to the story's playtime. )
func StoryStatement(run rt.Runtime, op Schedule) (err error) {
	if k, ok := run.(*Weaver); !ok {
		err = errutil.Fmt("runtime %T doesn't support story statements", run)
	} else {
		err = op.Schedule(k.d.cat)
	}
	return
}

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
	d := ctx.d
	if _, kind := d.UniformDeterminer(kindName); len(kind) == 0 {
		err = InvalidString(kindName)
	} else if field, ok := UniformString(fieldName); !ok {
		err = InvalidString(fieldName)
	} else if class, ok := UniformString(fieldClass); !ok && len(fieldClass) > 0 {
		err = InvalidString(fieldClass)
	} else {
		err = addField(kind, field, class)
	}
	return
}

func fromTiming(timing assert.EventTiming) int {
	var part int
	always := timing&assert.RunAlways != 0
	if always {
		timing ^= assert.RunAlways
	}
	switch timing {
	case assert.Before:
		part = 0
	case assert.During:
		part = 1
	case assert.After:
		part = 2
	case assert.Later:
		part = 3
	}
	flags := part + int(rt.FirstPhase)
	if always {
		flags = -flags // marker for rules that need to always run (ex. counters "every third try" )
	}
	return flags
}
