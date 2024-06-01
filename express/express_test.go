package express

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/template"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// var True = B(true)
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
			&math.CompareText{
				A: T("a"), Compare: math.C_Comparison_LessThan, B: T("b"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("num cmp", func(t *testing.T) {
		if e := testExpression(
			"7 >= 8",
			&math.CompareNum{
				A: F(7), Compare: math.C_Comparison_AtLeast, B: F(8),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("math", func(t *testing.T) {
		if e := testExpression(
			"(5+6)*(1+2)",
			&math.MultiplyValue{
				A: &math.AddValue{A: F(5), B: F(6)},
				B: &math.AddValue{A: F(1), B: F(2)},
			}); e != nil {
			t.Fatal(e)
		}
	})
	// isNot requires command parsing
	t.Run("logic", func(t *testing.T) {
		if e := testExpression(
			"true and (false or {not: true})",
			&logic.IsAll{
				Test: []rt.BoolEval{
					B(true),
					&logic.IsAny{
						Test: []rt.BoolEval{
							B(false),
							// isNot requires command parsing
							&logic.Not{
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
			renderRef("A", "num")); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("little dot", func(t *testing.T) {
		if e := testExpression(".a.b.c",
			renderRef("a", "b", "c")); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("binary", func(t *testing.T) {
		if e := testExpression(".A.num * .b.num",
			&math.MultiplyValue{
				A: renderRef("A", "num"),
				B: renderRef("b", "num"),
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
			&text.PrintNumWord{
				Num: &render.RenderRef{
					Name: T("group_size"),
				},
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("cycle", func(t *testing.T) {
		if e := testTemplate("{cycle}a{or}b{or}c{end}",
			&format.CallCycle{
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
			&format.CallTerminal{
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
			&format.CallShuffle{
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
			&logic.ChooseText{
				If: &math.CompareNum{
					A: F(7), Compare: math.C_Comparison_EqualTo, B: F(7),
				},
				True:  T("boop"),
				False: T("beep"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("unless", func(t *testing.T) {
		if e := testTemplate("{unless 7=7}boop{otherwise}beep{end}",
			&logic.ChooseText{
				If: &logic.Not{
					Test: &math.CompareNum{
						A: F(7), Compare: math.C_Comparison_EqualTo, B: F(7),
					}},
				True:  T("boop"),
				False: T("beep"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("filter", func(t *testing.T) {
		if e := testTemplate("{15|print_num!}",
			&text.PrintNum{
				Num: I(15),
			}); e != nil {
			t.Fatal(e)
		}
	})
	// all of the text in a template gets turned into an expression
	// plain text between bracketed sections becomes text evals
	t.Run("span", func(t *testing.T) {
		if e := testTemplate("{15|print_num!} {if 7=7}boop{end}",
			&text.Join{
				Parts: []rt.TextEval{
					&text.PrintNum{Num: F(15)},
					T(" "),
					&logic.ChooseText{
						If: &math.CompareNum{
							A: F(7), Compare: math.C_Comparison_EqualTo, B: F(7),
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
				PatternName: ("hello"), Render: []render.RenderEval{
					&render.RenderValue{Value: &assign.FromText{Value: T("world")}},
				}}); e != nil {
			t.Fatal(e)
		}
	})
	// dotted names standing alone in a template become requests to print its friendly name
	// as a lowercase dotted name, we try to get the actual object name first from a variable named "object"
	t.Run("object", func(t *testing.T) {
		if e := testTemplate("hello {.object}",
			&text.Join{Parts: []rt.TextEval{
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
			renderRef("Object", "prop"),
		); e != nil {
			t.Fatal(e)
		}
	})
}

func renderRef(v string, path ...any) *render.RenderRef {
	return &render.RenderRef{Name: T(v), Dot: object.MakeDot(path...)}
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
