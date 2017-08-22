package define_test

import (
	"github.com/ionous/iffy/dl/define"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {
	var reg define.Registry
	var facts define.Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
}

func TestGrammar(t *testing.T) {
	var reg define.Registry
	reg.Register(defineGrammar)
	//
	var facts define.Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	if testify.Len(t, facts.Grammar.Match, 1) {
		x, ok := facts.Grammar.Match[0].(*parser.AllOf)
		testify.True(t, ok)
		testify.Len(t, x.Match, 2) // l/look;action
	}
}

func TestLocation(t *testing.T) {
	var reg define.Registry
	reg.Register(defineLocation)
	//
	var facts define.Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	testify.Len(t, facts.Locations, 1)
}

func TestRules(t *testing.T) {
	var reg define.Registry
	mandates := []string{"bool", "number", "text", "object", "num list", "text list", "obj list", "run"}
	reg.Register(func(c *ops.Builder) {
		defineRules(c, mandates)
	})
	//
	var facts define.Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	testify.Len(t, facts.Mandates, len(mandates))
}

func TestEvents(t *testing.T) {
	var reg define.Registry
	reg.Register(defineEventHandler)
	//
	var facts define.Facts
	e := reg.Define(&facts)
	testify.NoError(t, e)
	testify.Len(t, facts.Listeners, 1)
}

func defineGrammar(c *ops.Builder) {
	if c.Cmd("grammar").Begin() {
		if c.Cmd("all of").Begin() {
			if c.Cmds().Begin() {
				if c.Cmd("any of").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("word", "l")
						c.Cmd("word", "look")
						c.End()
					}
					c.End()
				}
				if c.Cmd("any of").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("action", "look")
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.End()
		}
		c.End()
	}
}

func defineLocation(c *ops.Builder) {
	c.Cmd("location", "parent", c.Cmd("supports"), "child")
}

func defineRules(c *ops.Builder, mandates []string) {
	for _, k := range mandates {
		if c.Cmd("mandate").Begin() {
			if c.Cmd(k + " rule").Begin() {
				c.Param("name").Val(k)
				c.End()
			}
			c.End()
		}
	}
}

// future: defineEvent, defineInstance, defineRelative, defineClass, definePattern, defineRelation.

func defineEventHandler(c *ops.Builder) {
	if c.Cmd("listen to", "bogart", "jump").Begin() {
		if c.Param("go").Cmds().Begin() {
			c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
			c.Cmd("print text", "jumping!")
			c.End()
		}
		if c.Param("options").Cmds().Begin() {
			c.Cmd("capture")
			c.Cmd("target only")
			c.End()
		}
		c.End()
	}
	if c.Cmd("mandate").Begin() {
		if c.Cmd("run rule", "jump").Begin() {
			if c.Param("decide").Cmds().Begin() {
				if c.Cmd("print span").Begin() {
					if c.Cmds().Begin() {
						c.Cmd("determine", c.Cmd("print name", c.Cmd("get", "@", "target")))
						c.Cmd("print text", "jumped!")
						c.End()
					}
					c.End()
				}
				c.End()
			}
			c.Param("continue").Cmd("continue after")
			c.End()
		}
		c.End()
	}
}

// future: defineEvent, defineInstance, defineRelative, defineClass, definePattern, defineRelation, definePlural
