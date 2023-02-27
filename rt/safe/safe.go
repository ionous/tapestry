package safe

import (
	"io"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

// MissingEval error type for unknown variables while processing loops.
type MissingEval string

// Error returns the name of the unknown variable.
func (e MissingEval) Error() string { return "missing " + string(e) }

// Run executes the passed statement using the passed runtime;
// does *not* error if the passed exec is nil.
func RunAll(run rt.Runtime, exes []rt.Execute) (err error) {
	for i, exe := range exes {
		if exe != nil {
			if e := exe.Execute(run); e != nil {
				err = errutil.New("failed statement", i, e)
				break
			}
		}
	}
	return
}

// Run executes the passed statement using the passed runtime;
// does *not* error if the passed exec is nil.
func Run(run rt.Runtime, exe rt.Execute) (err error) {
	if exe == nil {
		err = MissingEval("execute")
	} else {
		err = exe.Execute(run)
	}
	return
}

// WriteText evaluates t and outputs the results to w.
func WriteText(run rt.Runtime, eval rt.TextEval) (err error) {
	if t, e := GetText(run, eval); e != nil {
		err = e
	} else if w := run.Writer(); w == nil {
		err = errutil.New("missing writer")
	} else {
		_, e := io.WriteString(w, t.String())
		err = e
	}
	return
}

// GetBool runs the specified eval, returning an error if the eval is nil.
func GetBool(run rt.Runtime, eval rt.BoolEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("boolean")
	} else {
		ret, err = eval.GetBool(run)
	}
	return
}

// GetNumber runs the specified eval, returning an error if the eval is nil.
func GetNumber(run rt.Runtime, eval rt.NumberEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("number")
	} else {
		ret, err = eval.GetNumber(run)
	}
	return
}

// GetText runs the specified eval, returning an error if the eval is nil.
func GetText(run rt.Runtime, eval rt.TextEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("text")
	} else {
		ret, err = eval.GetText(run)
	}
	return
}

// GetRecord runs the specified eval, returning an error if the eval is nil.
func GetRecord(run rt.Runtime, eval rt.RecordEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("record")
	} else {
		ret, err = eval.GetRecord(run)
	}
	return
}

// GetOptionalBool runs the optionally specified eval.
func GetOptionalBool(run rt.Runtime, eval rt.BoolEval, fallback bool) (ret g.Value, err error) {
	if eval == nil {
		ret = g.BoolOf(fallback)
	} else {
		ret, err = eval.GetBool(run)
	}
	return
}

// GetOptionalNumber runs the optionally specified eval.
func GetOptionalNumber(run rt.Runtime, eval rt.NumberEval, fallback float64) (ret g.Value, err error) {
	if eval == nil {
		ret = g.FloatOf(fallback)
	} else {
		ret, err = eval.GetNumber(run)
	}
	return
}

// GetOptionalText runs the optionally specified eval.
func GetOptionalText(run rt.Runtime, eval rt.TextEval, fallback string) (ret g.Value, err error) {
	if eval == nil {
		ret = g.StringOf(fallback)
	} else {
		ret, err = eval.GetText(run)
	}
	return
}

// GetOptionalNumber runs the optionally specified eval.
func GetOptionalNumbers(run rt.Runtime, eval rt.NumListEval, fallback []float64) (ret g.Value, err error) {
	if eval == nil {
		ret = g.FloatsOf(fallback)
	} else {
		ret, err = GetNumList(run, eval)
	}
	return
}

// GetOptionalText runs the optionally specified eval.
func GetOptionalTexts(run rt.Runtime, eval rt.TextListEval, fallback []string) (ret g.Value, err error) {
	if eval == nil {
		ret = g.StringsOf(fallback)
	} else {
		ret, err = GetTextList(run, eval)
	}
	return
}

// GetNumList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetNumList(run rt.Runtime, eval rt.NumListEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("num list")
	} else {
		ret, err = eval.GetNumList(run)
	}
	return
}

// GetTextList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetTextList(run rt.Runtime, eval rt.TextListEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("text list")
	} else {
		ret, err = eval.GetTextList(run)
	}
	return
}

// GetRecordList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetRecordList(run rt.Runtime, eval rt.RecordListEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("record list")
	} else {
		ret, err = eval.GetRecordList(run)
	}
	return
}

// ObjectText - given an eval producing a name, return a string value of the object's id.
// can return a valid "empty" value for empty strings
func ObjectText(run rt.Runtime, eval rt.TextEval) (ret g.Value, err error) {
	if eval == nil {
		err = MissingEval("object text")
	} else if t, e := eval.GetText(run); e != nil {
		err = e
	} else if n := t.String(); len(n) == 0 {
		ret = g.Empty
	} else {
		ret, err = run.GetField(meta.ObjectId, n)
	}
	return
}

func IsKindOf(run rt.Runtime, obj, kind string) (ret bool, err error) {
	if objectPath, e := run.GetField(meta.ObjectKinds, obj); e != nil {
		err = e
	} else {
		for _, k := range objectPath.Strings() {
			if k == kind {
				ret = true
				break
			}
		}
	}
	return
}
