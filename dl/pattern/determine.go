package pattern

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// FromPattern helps runs a pattern
type FromPattern struct {
	Pattern   string     // pattern name, i guess a text eval here would be like a function pointer.
	Arguments *Arguments // arguments passed to the pattern. kept as a pointer for composer formatting...
	// each is a name targeting some parameter, and an "assignment"
}

type DetermineAct FromPattern
type DetermineNum FromPattern
type DetermineText FromPattern
type DetermineBool FromPattern
type DetermineNumList FromPattern
type DetermineTextList FromPattern

// Stitch finds the pattern, builds the scope, and executes the passed callback to generate a result.
// It's an adapter from the the specific DetermineActivity, DetermineNumber, etc. statements.
func (op *FromPattern) Stitch(run rt.Runtime, pat Pattern, fn func() error) (err error) {
	// find the pattern (p), qna's implementation assembles the rules by querying the db.
	if e := run.GetEvalByName(op.Pattern, pat); e != nil {
		err = e
	} else {
		// bake the parameters down
		parms := make(Parameters)
		pat.Prepare(parms)

		// read from each argument
		if op.Arguments != nil {
			for _, arg := range op.Arguments.Args {
				if name, e := getParamName(pat, arg); e != nil {
					err = errutil.Append(err, e)
				} else if e := arg.From.Assign(run, func(val interface{}) error {
					return parms.SetVariable(name, val)
				}); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		if err == nil {
			run.PushScope(&parms)
			err = fn()
			run.PopScope()
		}
	}
	return
}

// change a param name ( which could be an index ) into a valid param name
func getParamName(pat Pattern, arg *Argument) (ret string, err error) {
	param := arg.Name
	if usesIndex := len(param) > 1 && param[:1] == "$"; !usesIndex {
		ret = param
	} else if idx, e := strconv.Atoi(param[1:]); e != nil {
		err = e
	} else {
		// parameters are 1 indexed right now
		ret, err = pat.GetParameterName(idx - 1)
	}
	return
}

func (*DetermineAct) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_act",
		Group: "patterns",
		Desc:  "Determine an activity",
	}
}

func (op *DetermineAct) Execute(run rt.Runtime) (err error) {
	var pat ActivityPattern
	err = (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		err = pat.Execute(run)
		return
	})
	return
}

func (*DetermineNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num",
		Spec:  "the {number pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a number",
	}
}

// GetNumber returns the first matching num evaluation.
func (op *DetermineNum) GetNumber(run rt.Runtime) (ret float64, err error) {
	var pat NumberPattern
	err = (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetNumber(run)
		return
	})
	return
}

func (*DetermineText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text",
		Spec:  "the {text pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine some text",
	}
}

// GetText returns the first matching text evaluation.
func (op *DetermineText) GetText(run rt.Runtime) (ret string, err error) {
	var pat TextPattern
	err = (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetText(run)
		return
	})
	return
}

func (*DetermineBool) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_bool",
		Spec:  "the {true/false pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a true/false value",
	}
}

// GetBool returns the first matching bool evaluation.
func (op *DetermineBool) GetBool(run rt.Runtime) (ret bool, err error) {
	var pat BoolPattern
	err = (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetBool(run)
		return
	})
	return
}

func (*DetermineNumList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_num_list",
		Spec:  "the {number list pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a list of numbers",
	}
}

func (op *DetermineNumList) GetNumberStream(run rt.Runtime) (ret rt.Iterator, err error) {
	var pat NumListPattern
	err = (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetNumberStream(run)
		return
	})
	return
}

func (*DetermineTextList) Compose() composer.Spec {
	return composer.Spec{
		Name:  "determine_text_list",
		Spec:  "the {text list pattern%name:pattern_name}{?parameters}",
		Group: "patterns",
		Desc:  "Determine a list of text",
	}
}

func (op *DetermineTextList) GetTextStream(run rt.Runtime) (ret rt.Iterator, err error) {
	var pat TextListPattern
	err = (*FromPattern)(op).Stitch(run, &pat, func() (err error) {
		ret, err = pat.GetTextStream(run)
		return
	})
	return
}