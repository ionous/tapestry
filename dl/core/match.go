package core

import (
	"regexp"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type MatchCache struct {
	GobNeedsOnePublicField int // https://github.com/golang/go/issues/19969
	err                    error
	exp                    *regexp.Regexp
}

func (op *Matches) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if exp, e := op.getRegexp(); e != nil {
		err = cmdError(op, e)
	} else {
		b := exp.MatchString(text.String())
		ret = g.BoolOf(b)
	}
	return
}

func (op *Matches) getRegexp() (ret *regexp.Regexp, err error) {
	if e := op.Cache.err; e != nil {
		err = e
	} else if exp := op.Cache.exp; exp != nil {
		ret = exp
	} else if exp, e := regexp.Compile(op.Pattern); e != nil {
		op.Cache.err = err
		err = e
	} else {
		op.Cache.exp = exp
		ret = exp
	}
	return
}
