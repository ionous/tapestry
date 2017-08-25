package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestArticles(t *testing.T) {
	suite.Run(t, new(ArticleSuite))
}

// Regular expression to select test suites specified command-line argument "-run".
// Regular expression to select the methods of test suites specified command-line argument "-m"
type ArticleSuite struct {
	suite.Suite
	run   rt.Runtime
	cmds  *ops.Ops
	lines printer.Lines
}

func (assert *ArticleSuite) Lines() (ret []string) {
	ret = assert.lines.Lines()
	assert.lines = printer.Lines{}
	return
}

func (assert *ArticleSuite) SetupTest() {
	errutil.Panic = false
	//
	cmds := ops.NewOps()

	// FIX? i dont like that article suite depends on patterns
	// can we provide a default print name somehow?
	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*Commands)(nil),
		(*core.Commands)(nil),
		(*rule.Commands)(nil),
	)

	classes := ref.NewClasses()
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*Kind)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		&Kind{Name: "lamp-post"},
		&Kind{Name: "soldiers", IndefiniteArticle: "some"},
		&Kind{Name: "trevor", CommonProper: ProperNamed},
	)
	patterns, e := patbuilder.NewPatternMaster(cmds, classes,
		(*Patterns)(nil)).Build(
		printNamePatterns,
	)
	assert.NoError(e)

	assert.cmds = cmds
	assert.run = rtm.New(classes).Objects(objects).Writer(&assert.lines).Patterns(patterns).Rtm()
}

func (assert *ArticleSuite) match(expected string, run func(c *ops.Builder)) {
	var root struct{ Eval rt.Execute }
	if c, ok := assert.cmds.NewBuilder(&root); ok {
		run(c)
		if e := c.Build(); assert.NoError(e) {
			if e := root.Eval.Execute(assert.run); assert.NoError(e) {
				lines := assert.Lines()
				assert.Equal(sliceOf.String(expected), lines)
			}
		}
	}
}

// lower a/n
func (assert *ArticleSuite) TestATrailingLampPost() {
	assert.match("You can only just make out a lamp-post.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower a/n", "lamp post"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestATrailingTrevor() {
	assert.match("You can only just make out Trevor.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower a/n", "trevor"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestATrailingSoldiers() {
	assert.match("You can only just make out some soldiers.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower a/n", "soldiers"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

// upper a/n
func (assert *ArticleSuite) TestALeadingLampPost() {
	assert.match("A lamp-post can be made out in the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper a/n", "lamp post"))
					c.Cmd("print text", "can be made out in the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestALeadingTrevor() {
	assert.match("Trevor can be made out in the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper a/n", "trevor"))
					c.Cmd("print text", "can be made out in the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestALeadingSoldiers() {
	assert.match("Some soldiers can be made out in the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper a/n", "soldiers"))
					c.Cmd("print text", "can be made out in the mist.")
					c.End()
				}
				c.End()
			}
		})
}

// lower-the
func (assert *ArticleSuite) TestTheTrailingLampPost() {
	assert.match("You can only just make out the lamp-post.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower the", "lamp post"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestTheTrailingTrevor() {
	assert.match("You can only just make out Trevor.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower the", "trevor"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestTheTrailingSoldiers() {
	assert.match("You can only just make out the soldiers.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower the", "soldiers"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

// uppe the
func (assert *ArticleSuite) TestTheLeadingLampPost() {
	assert.match("The lamp-post may be a trick of the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper the", "lamp post"))
					c.Cmd("print text", "may be a trick of the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestTheLeadingTrevor() {
	assert.match("Trevor may be a trick of the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper the", "trevor"))
					c.Cmd("print text", "may be a trick of the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (assert *ArticleSuite) TestTheLeadingSoldiers() {
	assert.match("The soldiers may be a trick of the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper the", "soldiers"))
					c.Cmd("print text", "may be a trick of the mist.")
					c.End()
				}
				c.End()
			}
		})
}

// FIX: should really be separate -- in a "text" test.
func (assert *ArticleSuite) TestPluralize() {
	assert.match("lamps",
		func(c *ops.Builder) {
			if c.Cmd("print span").Begin() {
				c.Cmds(c.Cmd("print text", c.Cmd("pluralize", "lamp")))
				c.End()
			}
		})
}
