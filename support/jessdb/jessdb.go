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

func countWords(str string) (ret int) {
	if len(str) > 0 {
		ret = 1 + strings.Count(str, " ")
	}
	return
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
		ret, width = m.Name, countWords(m.Match)
	}
	return
}

func (d dbWrapper) FindTrait(ws match.Span) (ret string, width int) {
	str := strings.ToLower(ws.String())
	if trait, e := d.GetPartialTrait(str); e != nil {
		log.Println("FindTrait", e)
	} else {
		// the returned name is the name of the trait from the db
		// it was used to match the front of the passed string
		// so the words in the trait are the words in the string.
		ret, width = trait, countWords(trait)
	}
	return
}

func (d dbWrapper) FindField(kind string, field match.Span) (ret string, width int) {
	str := strings.ToLower(field.String())
	if field, e := d.GetPartialField(kind, str); e != nil {
		log.Println("FindField", e)
	} else {
		// re: countWords, same logic as find trait.
		ret, width = field, countWords(field)
	}
	return
}

func (d dbWrapper) FindNoun(ws match.Span, pkind *string) (ret string, width int) {
	if n, e := d.findNoun(ws, pkind); e == nil {
		ret, width = n, countWords(n)
	} else if !errors.Is(e, mdl.Missing) {
		log.Println("FindNoun", e)
	}
	return
}

func (d dbWrapper) findNoun(ws match.Span, pkind *string) (ret string, err error) {
	str := strings.ToLower(ws.String())
	var kind string
	if pkind != nil {
		kind = *pkind
	}
	if len(kind) == 0 {
		if n, k, e := d.GetClosestNoun(str); e != nil {
			err = e
		} else {
			ret = n
			if pkind != nil {
				*pkind = k
			}
		}
	} else {
		if m, e := d.GetPartialNoun(str, kind); e != nil {
			err = e
		} else {
			ret = m.Name
			if pkind != nil {
				*pkind = m.Kind
			}
		}
	}
	return
}
