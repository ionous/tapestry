package gestalt

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

// adapt the results of an input into grok's results
func Reduce(in InputState) (ret grok.Results, err error) {
	var d reducer
	d.setPrimary()

	handler := onMatch
	for _, res := range in.res {
		if next, e := handler(&d, res); e != nil {
			err = e
			break
		} else if next != nil {
			handler = next
		} else {
			handler = onMatch
		}
	}
	ret.Primary = append(ret.Primary, d.primary)
	if d.usedSecondary {
		ret.Secondary = append(ret.Secondary, d.secondary)
	}
	ret.Macro = d.macro
	return
}

type reducer struct {
	name               *grok.Name
	macro              grok.Macro
	usedSecondary      bool
	primary, secondary grok.Name
}

func (d *reducer) setPrimary() {
	d.name = &d.primary
}
func (d *reducer) setSecondary() {
	d.usedSecondary = true
	d.name = &d.secondary
}

type resultState func(*reducer, grok.Match) (resultState, error)

// default handler
func onMatch(d *reducer, m grok.Match) (ret resultState, err error) {
	switch m := m.(type) {
	case grok.Article:
		d.name.Article = m
		if m.Count > 0 {
			ret = onCounted
		}
	case grok.Macro:
		d.macro = m

	case MatchedKind:
		d.name.Kinds = append(d.name.Kinds, m.Span)

	case MatchedName:
		d.name.Span = m.Span

	case MatchedTrait:
		d.name.Traits = append(d.name.Traits, m)

	case MatchedTarget:
		if m {
			d.setPrimary()
		} else {
			d.setSecondary()
		}

	case MatchedVerb:
		if m.Action == "implies" {
			// fix: this is properly kinds, but grok produces a generic name
			// ( and the tests test for it. )
			k := d.name.Kinds[0]
			d.name.Span = k.(grok.Span)
			d.name.Kinds = nil
		}
		d.macro = grok.Macro{
			Name:  m.Action,
			Match: m,
			// Type: Primary, ManyPrimary, ManySecondary, ManyMany
			// Reversed bool
		}

	default:
		err = fmt.Errorf("unhandled result %T", m)
	}
	return
}

func onCounted(d *reducer, res grok.Match) (next resultState, err error) {
	switch m := res.(type) {
	case MatchedKind:
		d.name.Kinds = append(d.name.Kinds, m)
	default:
		err = fmt.Errorf("expected a name to follow 'called' %s", m)
	}
	return
}
