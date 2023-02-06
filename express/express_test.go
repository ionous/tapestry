package express

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/template"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

var True = B(true)
var False = B(false)

// TestExpressions single expressions within a template.
// ( the parts that normally appear inside curly brackets {here} ).
func TestExpressions(t *testing.T) {
	t.Run("num", func(t *testing.T) {
		if e := testExpression("5", F(5)); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("text", func(t *testing.T) {
		if e := testExpression("'5'", T("5")); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("bool", func(t *testing.T) {
		if e := testExpression("false", False); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("T cmp", func(t *testing.T) {
		if e := testExpression(
			"'a' < 'b'",
			&core.CompareText{
				A: T("a"), Is: core.LessThan, B: T("b"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("num cmp", func(t *testing.T) {
		if e := testExpression(
			"7 >= 8",
			&core.CompareNum{
				A: F(7), Is: core.AtLeast, B: F(8),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("math", func(t *testing.T) {
		if e := testExpression(
			"(5+6)*(1+2)",
			&core.ProductOf{
				A: &core.SumOf{A: F(5), B: F(6)},
				B: &core.SumOf{A: F(1), B: F(2)},
			}); e != nil {
			t.Fatal(e)
		}
	})
	// isNot requires command parsing
	t.Run("logic", func(t *testing.T) {
		if e := testExpression(
			"true and (false or {not: true})",
			&core.AllTrue{
				Test: []rt.BoolEval{
					B(true),
					&core.AnyTrue{
						Test: []rt.BoolEval{
							B(false),
							// isNot requires command parsing
							&core.Not{
								Test: B(true),
							},
						}},
				}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("global", func(t *testing.T) {
		if e := testExpression(".A",
			&render.RenderName{Name: "A"}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("big dot", func(t *testing.T) {
		// get 'num' out of 'A' ( which is in this case an object )
		if e := testExpression(".A.num",
			core.GetName("A", "num")); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("little dot", func(t *testing.T) {
		if e := testExpression(".a.b.c",
			core.GetName("a", "b", "c")); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("binary", func(t *testing.T) {
		if e := testExpression(".A.num * .b.num",
			&core.ProductOf{
				A: core.GetName("A", "num"),
				B: core.GetName("b", "num"),
			}); e != nil {
			t.Fatal(e)
		}
	})
}

func testExpression(str string, want interface{}) (err error) {
	if xs, e := template.ParseExpression(str); e != nil {
		err = errutil.New(e)
	} else if got, e := Convert(xs); e != nil {
		err = errutil.New(e)
	} else if diff := pretty.Diff(got, want); len(diff) > 0 {
		err = errutil.New("have:", pretty.Sprint(got), "want:", pretty.Sprint(want))
	}
	return
}

// test full templates
func TestTemplates(t *testing.T) {
	t.Run("print", func(t *testing.T) {
		if e := testTemplate("{print_num_word: .group_size}",
			&core.PrintNumWord{
				Num: &render.RenderRef{
					Name:  N("group_size"),
					Flags: render.RenderFlags{Str: render.RenderFlags_RenderAsAny},
				},
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("cycle", func(t *testing.T) {
		if e := testTemplate("{cycle}a{or}b{or}c{end}",
			&core.CallCycle{
				Name: "autoexp1",
				Parts: []rt.TextEval{
					T("a"), T("b"), T("c"),
				},
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("once", func(t *testing.T) {
		if e := testTemplate("{once}a{or}b{or}c{end}",
			&core.CallTerminal{
				Name: "autoexp1",
				Parts: []rt.TextEval{
					T("a"), T("b"), T("c"),
				},
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("shuffle", func(t *testing.T) {
		if e := testTemplate("{shuffle}a{or}b{or}c{end}",
			&core.CallShuffle{
				Name: "autoexp1",
				Parts: []rt.TextEval{
					T("a"), T("b"), T("c"),
				},
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("if", func(t *testing.T) {
		if e := testTemplate("{if 7=7}boop{else}beep{end}",
			&core.ChooseText{
				If: &core.CompareNum{
					A: F(7), Is: core.Equal, B: F(7),
				},
				True:  T("boop"),
				False: T("beep"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("unless", func(t *testing.T) {
		if e := testTemplate("{unless 7=7}boop{otherwise}beep{end}",
			&core.ChooseText{
				If: &core.Not{
					Test: &core.CompareNum{
						A: F(7), Is: core.Equal, B: F(7),
					}},
				True:  T("boop"),
				False: T("beep"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("filter", func(t *testing.T) {
		if e := testTemplate("{15|print_num!}",
			&core.PrintNum{
				Num: I(15),
			}); e != nil {
			t.Fatal(e)
		}
	})
	// all of the text in a template gets turned into an expression
	// plain text between bracketed sections becomes text evals
	t.Run("span", func(t *testing.T) {
		if e := testTemplate("{15|print_num!} {if 7=7}boop{end}",
			&core.Join{
				Parts: []rt.TextEval{
					&core.PrintNum{Num: F(15)},
					T(" "),
					&core.ChooseText{
						If: &core.CompareNum{
							A: F(7), Is: core.Equal, B: F(7),
						},
						True: T("boop"),
					},
				},
			}); e != nil {
			t.Fatal(e)
		}
	})
	// parameters to template calls become indexed parameter assignments
	t.Run("indexed", func(t *testing.T) {
		if e := testTemplate("{'world'|hello!}",
			&render.RenderPattern{
				Call: core.CallPattern{
					Pattern: P("hello"), Arguments: core.Args(
						&core.FromText{Val: T("world")},
					)}}); e != nil {
			t.Fatal(e)
		}
	})
	// dotted names standing alone in a template become requests to print its friendly name
	// as a lowercase dotted name, we try to get the actual object name first from a variable named "object"
	t.Run("object", func(t *testing.T) {
		if e := testTemplate("hello {.object}",
			&core.Join{Parts: []rt.TextEval{
				T("hello "),
				&render.RenderName{Name: "object"}}},
		); e != nil {
			t.Fatal(e)
		}
	})

	// dotted names started with capital letters are requests for objects exactly matching that name
	// note: this does the case check at runtime now, so there's no difference in the resulting commands b/t .Object and .object
	t.Run("global prop", func(t *testing.T) {
		if e := testTemplate("{.Object.prop}",
			core.GetName("Object", "prop"),
		); e != nil {
			t.Fatal(e)
		}
	})
}
func testTemplate(str string, want interface{}) (err error) {
	if xs, e := template.Parse(str); e != nil {
		err = e
	} else if got, e := Convert(xs); e != nil {
		err = errutil.New(e, xs)
	} else if diff := pretty.Diff(got, want); len(diff) > 0 {
		err = errutil.New("mismatch:", "got", pretty.Sprint(got),
			"want", pretty.Sprint(want))
	}
	return
}
