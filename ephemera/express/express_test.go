package express

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/template"
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
				T("a"), &core.LessThan{}, T("b"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("num cmp", func(t *testing.T) {
		if e := testExpression(
			"7 >= 8",
			&core.CompareNum{
				F(7), &core.AtLeast{}, F(8),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("math", func(t *testing.T) {
		if e := testExpression(
			"(5+6)*(1+2)",
			&core.ProductOf{
				&core.SumOf{F(5), F(6)},
				&core.SumOf{F(1), F(2)},
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
								B(true),
							},
						}},
				}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("global", func(t *testing.T) {
		if e := testExpression(".A",
			&render.RenderName{"A"}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("big dot", func(t *testing.T) {
		if e := testExpression(".A.num",
			// get 'num' out of 'A' ( note: get field at supports any value )
			&core.GetAtField{
				Field: W("num"),
				// get "a" -- some value supporting field access
				// could be a record or an object variable, or a global object.
				// ( b/c its capitalized, we know its going to be a global object )
				From: &render.RenderField{
					Name: T("A"),
				},
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("little dot", func(t *testing.T) {
		if e := testExpression(".a.b.c",
			// c, a value in b, can be anything.
			&core.GetAtField{
				Field: W("c"),
				// to get a value from b, b must have been specifically a record.
				From: &core.FromRec{
					// get b out of a ( note: get field at supports any value )
					Rec: &core.GetAtField{
						Field: W("b"),
						// get "a" -- some value supporting field access
						// could be a record or an object variable, or a global object.
						From: &render.RenderField{
							Name: T("a"),
						},
					},
				}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("binary", func(t *testing.T) {
		if e := testExpression(".A.num * .b.num",
			&core.ProductOf{
				A: &core.GetAtField{
					Field: W("num"),
					From:  &render.RenderField{T("A")},
				},
				B: &core.GetAtField{
					Field: W("num"),
					From:  &render.RenderField{T("b")},
				},
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
		err = errutil.New("failed:", pretty.Sprint(got))
	}
	return
}

// test full templates
func TestTemplates(t *testing.T) {
	t.Run("print", func(t *testing.T) {
		if e := testTemplate("{print_num_word: .group_size}",
			&core.PrintNumWord{
				Num: &render.RenderRef{
					N("group_size"),
					render.RenderFlags{Str: render.RenderFlags_RenderAsAny},
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
					F(7), &core.Equal{}, F(7),
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
					&core.CompareNum{
						F(7), &core.Equal{}, F(7),
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
				Num: &core.NumValue{15},
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
					&core.PrintNum{F(15)},
					T(" "),
					&core.ChooseText{
						If: &core.CompareNum{
							F(7), &core.Equal{}, F(7),
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
				Pattern: P("hello"), Arguments: core.Args(
					&core.FromText{T("world")},
				)}); e != nil {
			t.Fatal(e)
		}
	})
	// dotted names standing alone in a template become requests to print its friendly name
	// as a lowercase dotted name, we try to get the actual object name first from a variable named "object"
	t.Run("object", func(t *testing.T) {
		if e := testTemplate("hello {.object}",
			&core.Join{Parts: []rt.TextEval{
				T("hello "),
				&render.RenderName{"object"}}},
		); e != nil {
			t.Fatal(e)
		}
	})

	// dotted names started with capital letters are requests for objects exactly matching that name
	// note: we do the cap check at runtime now, so there's no difference in the resulting commands b/t .Object and .object
	t.Run("global prop", func(t *testing.T) {
		if e := testTemplate("{.Object.prop}",
			&core.GetAtField{
				Field: W("prop"),
				From: &render.RenderField{
					Name: T("Object"),
				},
			},
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
