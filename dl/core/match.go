package core

import (
	"regexp"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type Matches struct {
	Text    rt.TextEval
	Pattern string
	// should transform into a different command probably during compile
	exp *regexp.Regexp `if:"internal"`
	err error
}

type MatchLike struct {
	Text    rt.TextEval
	Pattern rt.TextEval
}

// Compose defines a spec for the composer editor.
func (*Matches) Compose() composer.Spec {
	return composer.Spec{
		Name:  "matches",
		Group: "matching",
		Desc:  `Matches: Determine whether the specified text is similar to the specified regular expression.`,
		Spec:  "{text:text_eval} matches {pattern:text}",
	}
}

// Compose defines a spec for the composer editor.
func (*MatchLike) Compose() composer.Spec {
	return composer.Spec{
		Name:  "match_like",
		Group: "matching",
		Desc: `Like: Determine whether the specified text is similar to the specified pattern.
		Matching is case-insensitive ( meaning, "A" matches "a" ) and there are two symbols with special meaning. 
		A percent sign ("%") in the pattern matches any series of zero or more characters in the original text, 
		while an underscore matches ("_") any one single character. `,
		Spec: "{text:text_eval} is like {pattern:text_eval}",
	}
}

func (op *Matches) GetBool(run rt.Runtime) (ret bool, err error) {
	if text, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if exp, e := op.getRegexp(); e != nil {
		err = cmdError(op, e)
	} else {
		ret = exp.MatchString(text)
	}
	return
}

func (op *Matches) getRegexp() (ret *regexp.Regexp, err error) {
	if e := op.err; e != nil {
		err = e
	} else if exp := op.exp; exp != nil {
		ret = exp
	} else if exp, e = regexp.Compile(op.Pattern); e != nil {
		op.err = err
		err = e
	} else {
		op.exp = exp
		ret = exp
	}
	return
}

func (op *MatchLike) GetBool(run rt.Runtime) (ret bool, err error) {
	type isLike interface {
		IsLike(text, pattern string) (bool, error)
	}
	if isLike, ok := run.(isLike); !ok {
		e := errutil.New("the runtime doesnt implement the like matcher interface")
		err = cmdError(op, e)
	} else if text, e := rt.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if pattern, e := rt.GetText(run, op.Pattern); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = isLike.IsLike(text, pattern)
	}
	return
}