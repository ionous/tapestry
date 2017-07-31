package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/initial"
	"github.com/ionous/iffy/dl/locate"
	"github.com/ionous/iffy/index"
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

func TestContents(t *testing.T) {
	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		Thingaverse.objects(sliceOf.String("box", "cake"))...)

	cmds := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*core.Commands)(nil),
		(*patspec.Commands)(nil),
		(*Commands)(nil),
		(*initial.Commands)(nil),
	)
	unique.RegisterBlocks(unique.PanicTypes(cmds.ShadowTypes),
		(*Patterns)(nil),
	)

	//t.Run("Names", func(t *testing.T) {
	assert := testify.New(t)

	patterns, e := patbuilder.NewPatternMaster(cmds, classes,
		(*Patterns)(nil)).Build(printPatterns)
	assert.NoError(e)

	pc := locate.Locale{index.NewTable(index.OneToMany)}
	type OpsCb func(c *ops.Builder)
	type Match func(lines []string) error
	test := func(build, exec OpsCb, match Match) (err error) {
		var src struct{ initial.Statements }
		if c, ok := cmds.NewBuilder(&src); !ok {
			err = errutil.New("no buidler")
		} else {
			if c.Cmds().Begin() {
				build(c)
				c.End()
			}
			if e := c.Build(); e != nil {
				err = e
			} else {
				var facts initial.Facts
				if e := src.Assess(&facts); e != nil {
					err = e
				} else {
					objs := objects.Build()
					// for _, v := range facts.Values {
					// 	if obj, ok := objs.GetObject(v.Obj); !ok {
					// 		t.Fatal("couldnt find", v.Obj)
					// 		break
					// 	} else if e := obj.SetValue(v.Prop, v.Val); e != nil {
					// 		t.Fatal(e)
					// 		break
					// 	}
					// }
					// for _, r := range facts.Relations {
					// 	if e := relations.NewRelation(r.Name, index.NewTable(r.Type)); e != nil {
					// 		t.Fatal(e)
					// 		break
					// 	}
					// }

					for _, l := range facts.Locations {
						// in this case we're probably a command too
						if p, ok := objs.GetObject(l.Parent); !ok {
							err = errutil.New("unknown", l.Parent)
							break
						} else if c, ok := objs.GetObject(l.Child); !ok {
							err = errutil.New("unknown", l.Child)
							break
						} else if e := pc.SetLocation(p, c, l.Relative); e != nil {
							err = e
							break
						}
					}

				}

			}
		}
		if err == nil {
			var root struct{ rt.ExecuteList }
			if c, ok := cmds.NewBuilder(&root); !ok {
				err = errutil.New("no builder")
			} else {
				if c.Cmds().Begin() {
					exec(c)
					c.End()
				}
				if e := c.Build(); e != nil {
					err = e
				} else {
					var lines printer.Lines
					run := rtm.New(classes).Objects(objects).Patterns(patterns).Writer(&lines).Rtm()
					if e := root.Execute(run); e != nil {
						err = e
					} else {
						err = match(lines.Lines())
					}
				}
			}
		}
		return
	}
	//

	t.Run("contains", func(t *testing.T) {
		assert := testify.New(t)
		e := test(func(c *ops.Builder) {
			c.Cmd("Location", "box", locate.Contains, "cake")
		}, func(c *ops.Builder) {
			//
		}, func(lines []string) (nil error) {
			// verify relation:
			if in, ok := pc.GetData("$box", "$cake"); assert.True(ok) {
				assert.EqualValues(locate.Contains, in)
			}
			return
		})
		assert.NoError(e)
	})

}
