package text

import (
	"regexp"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

type MatchCache struct {
	err error
	exp *regexp.Regexp
}

func (op *Matches) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if exp, e := op.getRegexp(); e != nil {
		err = cmd.Error(op, e)
	} else {
		b := exp.MatchString(text.String())
		ret = rt.BoolOf(b)
	}
	return
}

func (op *Matches) getRegexp() (ret *regexp.Regexp, err error) {
	if e := op.cache.err; e != nil {
		err = e
	} else if exp := op.cache.exp; exp != nil {
		ret = exp
	} else if exp, e := regexp.Compile(op.Match); e != nil {
		op.cache.err = err
		err = e
	} else {
		op.cache.exp = exp
		ret = exp
	}
	return
}
