package safe

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
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
				err = fmt.Errorf("failed statement %d %T %w", i, exe, e)
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
		err = errors.New("missing writer")
	} else {
		_, e := io.WriteString(w, t.String())
		err = e
	}
	return
}

// handles null assignments by returning "MissingEval" error
// ( cant live in package safe because package assign uses package safe )
func GetAssignment(run rt.Runtime, a rt.Assignment) (ret rt.Value, err error) {
	if a == nil {
		err = MissingEval("assigned value")
	} else {
		ret, err = a.GetAssignedValue(run)
	}
	return
}

// GetBool runs the specified eval, returning an error if the eval is nil.
func GetBool(run rt.Runtime, eval rt.BoolEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("boolean")
	} else {
		ret, err = eval.GetBool(run)
	}
	return
}

// GetNum runs the specified eval, returning an error if the eval is nil.
func GetNum(run rt.Runtime, eval rt.NumEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("num")
	} else {
		ret, err = eval.GetNum(run)
	}
	return
}

// GetText runs the specified eval, returning an error if the eval is nil.
func GetText(run rt.Runtime, eval rt.TextEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("text")
	} else {
		ret, err = eval.GetText(run)
	}
	return
}

// GetRecord runs the specified eval, returning an error if the eval is nil.
func GetRecord(run rt.Runtime, eval rt.RecordEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("record")
	} else {
		ret, err = eval.GetRecord(run)
	}
	return
}

// GetOptionalBool runs the optionally specified eval.
func GetOptionalBool(run rt.Runtime, eval rt.BoolEval, fallback bool) (ret rt.Value, err error) {
	if eval == nil {
		ret = rt.BoolOf(fallback)
	} else {
		ret, err = eval.GetBool(run)
	}
	return
}

// GetOptionalNumber runs the optionally specified eval.
func GetOptionalNumber(run rt.Runtime, eval rt.NumEval, fallback float64) (ret rt.Value, err error) {
	if eval == nil {
		ret = rt.FloatOf(fallback)
	} else {
		ret, err = eval.GetNum(run)
	}
	return
}

// GetOptionalText runs the optionally specified eval.
func GetOptionalText(run rt.Runtime, eval rt.TextEval, fallback string) (ret rt.Value, err error) {
	if eval == nil {
		ret = rt.StringOf(fallback)
	} else {
		ret, err = eval.GetText(run)
	}
	return
}

// GetOptionalNumber runs the optionally specified eval.
func GetOptionalNumbers(run rt.Runtime, eval rt.NumListEval, fallback []float64) (ret rt.Value, err error) {
	if eval == nil {
		ret = rt.FloatsOf(fallback)
	} else {
		ret, err = GetNumList(run, eval)
	}
	return
}

// GetOptionalText runs the optionally specified eval.
func GetOptionalTexts(run rt.Runtime, eval rt.TextListEval, fallback []string) (ret rt.Value, err error) {
	if eval == nil {
		ret = rt.StringsOf(fallback)
	} else {
		ret, err = GetTextList(run, eval)
	}
	return
}

// GetNumList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetNumList(run rt.Runtime, eval rt.NumListEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("num list")
	} else {
		ret, err = eval.GetNumList(run)
	}
	return
}

// GetTextList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetTextList(run rt.Runtime, eval rt.TextListEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("text list")
	} else {
		ret, err = eval.GetTextList(run)
	}
	return
}

// GetRecordList returns an new iterator to walk the passed list,
// or an empty iterator if the value is null.
func GetRecordList(run rt.Runtime, eval rt.RecordListEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("record list")
	} else {
		ret, err = eval.GetRecordList(run)
	}
	return
}

// ObjectText - given an eval producing a name, return a string value of the object's id.
// can return a valid "empty" value for empty strings
func ObjectText(run rt.Runtime, eval rt.TextEval) (ret rt.Value, err error) {
	if eval == nil {
		err = MissingEval("object text")
	} else if t, e := eval.GetText(run); e != nil {
		err = e
	} else if n := t.String(); len(n) == 0 {
		ret = rt.Empty
	} else {
		ret, err = run.GetField(meta.ObjectId, n)
	}
	return
}

func IsKindOf(run rt.Runtime, obj, kind string) (ret bool, err error) {
	if path, e := run.GetField(meta.ObjectKinds, obj); e != nil {
		err = e
	} else {
		ret = slices.Contains(path.Strings(), kind)
	}
	return
}
