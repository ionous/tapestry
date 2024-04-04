package jessdb

import (
	"errors"
	"log"
	"strings"

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

func (d dbWrapper) FindKind(ws match.Span, out *kindsOf.Kinds) (ret string, width int) {
	str := strings.ToLower(ws.String()) // fix: so not a fan of these string lowers
	if m, e := d.GetPartialKind(str); e != nil {
		log.Println("FindKind", e)
	} else {
		if out != nil {
			*out = m.Base
		}
		ret, width = m.Name, m.WordCount()
	}
	return
}

func (d dbWrapper) FindTrait(ws match.Span) (ret string, width int) {
	str := strings.ToLower(ws.String())
	if m, e := d.GetPartialTrait(str); e != nil {
		log.Println("FindTrait", e)
	} else {
		ret, width = m.Name, m.WordCount()
	}
	return
}

func (d dbWrapper) FindField(kind string, field match.Span) (ret string, width int) {
	str := strings.ToLower(field.String())
	if m, e := d.GetPartialField(kind, str); e != nil {
		log.Println("FindField", e)
	} else {
		ret, width = m.Name, m.WordCount()
	}
	return
}

func (d dbWrapper) FindNoun(ws match.Span, pkind *string) (ret string, width int) {
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

func (d dbWrapper) findNoun(ws match.Span, pkind *string) (ret mdl.MatchedNoun, err error) {
	str := strings.ToLower(ws.String())
	var kind string
	if pkind != nil {
		kind = *pkind
	}
	if len(kind) == 0 {
		ret, err = d.GetClosestNoun(str)
	} else {
		ret, err = d.GetPartialNoun(str, kind)
	}
	return
}
