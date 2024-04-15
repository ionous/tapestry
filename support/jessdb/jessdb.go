package jessdb

import (
	"errors"
	"log"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func MakeQuery(m *mdl.Modeler, domain string) jess.Query {
	return dbWrapper{m.Pin(domain, "jess")}
}

// implements jess.Query; returned by dbWrapper.
type dbWrapper struct {
	*mdl.Pen
}

func (d dbWrapper) GetContext() int {
	return 0
}

func (d dbWrapper) FindKind(ws []match.TokenValue, out *kindsOf.Kinds) (ret string, width int) {
	if str, min := match.Normalize(ws); min > 0 {
		if m, e := d.GetPartialKind(str); e != nil {
			log.Println("FindKind", e)
		} else {
			if out != nil {
				*out = m.Base
			}
			ret, width = m.Name, m.WordCount()
		}
	}
	return
}

func (d dbWrapper) FindTrait(ws []match.TokenValue) (ret string, width int) {
	if str, min := match.Normalize(ws); min > 0 {
		if m, e := d.GetPartialTrait(str); e != nil {
			log.Println("FindTrait", e)
		} else {
			ret, width = m.Name, m.WordCount()
		}
	}
	return
}

func (d dbWrapper) FindField(kind string, ws []match.TokenValue) (ret string, width int) {
	if str, min := match.Normalize(ws); min > 0 {
		if m, e := d.GetPartialField(kind, str); e != nil {
			log.Println("FindField", e)
		} else {
			ret, width = m.Name, m.WordCount()
		}
	}
	return
}

func (d dbWrapper) FindNoun(ws []match.TokenValue, pkind *string) (ret string, width int) {
	if m, e := d.findNoun(ws, pkind); e != nil {
		if !errors.Is(e, mdl.Missing) {
			log.Println("FindNoun", e)
		}
	} else {
		ret, width = m.Name, m.WordCount()
		if pkind != nil {
			*pkind = m.Kind
		}
	}
	return
}

func (d dbWrapper) findNoun(ws []match.TokenValue, pkind *string) (ret mdl.MatchedNoun, err error) {
	if str, min := match.Normalize(ws); min > 0 {
		var kind string
		if pkind != nil {
			kind = *pkind
		}
		if len(kind) == 0 {
			ret, err = d.GetClosestNoun(str)
		} else {
			ret, err = d.GetPartialNoun(str, kind)
		}
	}
	return
}
