package std_test

import (
	"github.com/ionous/iffy/dl/core"
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestStd(t *testing.T) {
	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		Thingaverse.objects(sliceOf.String("apple", "pen", "thing#1", "thing#2"))...)

	cmds := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*core.Commands)(nil),
		(*patspec.Commands)(nil),
		(*Commands)(nil),
	)

	unique.RegisterBlocks(unique.PanicTypes(cmds.ShadowTypes),
		(*Patterns)(nil),
	)

	t.Run("Names", func(t *testing.T) {
		assert := testify.New(t)

		//
		patterns, e := patbuilder.NewPatternMaster(cmds, classes,
			(*Patterns)(nil)).Build(PrintNamePatterns)
		assert.NoError(e)

		// TODO: add test for: Rule for printing the name of the pen while taking inventory: say "useful pen".
		// TODO: add test for: A novel is a kind of thing. Dr Zhivago and Persuasion are novels. Before printing the name of a novel, say "[italic type]". After printing the name of a novel, say "[roman type]".”
		run := rtm.New(classes).Objects(objects).Patterns(patterns).Rtm()
		//
		t.Run("printed name", func(t *testing.T) {
			apple := Thingaverse["apple"].(*Thing)
			match(run, testify.New(t), "empire apple", &PrintName{&apple.Kind})
		})
		t.Run("named", func(t *testing.T) {
			genericPen := Thingaverse["pen"].(*Thing)
			match(run, testify.New(t), "pen", &PrintName{&genericPen.Kind})
		})
		t.Run("unnamed", func(t *testing.T) {
			unnamedThing := Thingaverse["thing#1"].(*Thing)
			match(run, testify.New(t), "thing", &PrintName{&unnamedThing.Kind})
		})
		//
		t.Run("plural printed name", func(t *testing.T) {
			apple := Thingaverse["apple"].(*Thing)
			match(run, testify.New(t), "empire apples", &PrintPluralName{&apple.Kind})
		})
		t.Run("plural named", func(t *testing.T) {
			genericPen := Thingaverse["pen"].(*Thing)
			match(run, testify.New(t), "pens", &PrintPluralName{&genericPen.Kind})
		})
		t.Run("plural unnamed", func(t *testing.T) {
			unnamedThing := Thingaverse["thing#1"].(*Thing)
			match(run, testify.New(t), "things", &PrintPluralName{&unnamedThing.Kind})
		})
		//
		t.Run("printed plural name", func(t *testing.T) {
			forced := "party favors"
			apple := Thingaverse["apple"].(*Thing)
			apple.PrintedPluralName = forced
			match(run, testify.New(t), forced, &PrintPluralName{&apple.Kind})
		})
		//
	})
}

func match(run rt.Runtime, assert *testify.Assertions, match string, op interface{}) (okay bool) {
	var lines printer.Lines
	run = rt.Writer(run, &lines)
	if patdata, e := run.Emplace(op); assert.NoError(e) {
		if e := run.ExecuteMatching(run, patdata); assert.NoError(e) {
			okay = assert.EqualValues(sliceOf.String(match), lines.Lines())
		}
	}
	return
}
