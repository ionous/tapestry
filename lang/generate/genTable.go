package generate

import (
	"fmt"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables/idl"
)

func writeTable(w DB, g Group) (err error) {
	gc := g.groupContent
	if e := writeTypes(w, g.Name, "slot", gc.Slot); e != nil {
		err = fmt.Errorf("%w writing slot for %s", e, g.Name)

	} else if e := writeTypes(w, g.Name, "flow", gc.Flow); e != nil {
		err = fmt.Errorf("%w writing flow for %s", e, g.Name)

	} else if e := writeTypes(w, g.Name, "str", gc.Str); e != nil {
		err = fmt.Errorf("%w writing str for %s", e, g.Name)

	} else if e := writeTypes(w, g.Name, "num", gc.Num); e != nil {
		err = fmt.Errorf("%w writing num for %s", e, g.Name)

	} else if e := writeStr(w, g.Name, gc.Str); e != nil {
		err = fmt.Errorf("%w writing str for %s", e, g.Name)
	}
	for _, sig := range gc.Reg {
		w.Write(idl.Sig, sig.Type, sig.Slot, strconv.FormatUint(sig.Hash, 16), sig.Body())
	}
	return
}

func writeReferences(w DB, g Group) (err error) {
	gc := g.groupContent
	if e := writeTerms(w, g.Name, gc.Flow); e != nil {
		err = fmt.Errorf("%w writing flow for %s", e, g.Name)
	}
	return
}

func writeTerms(w DB, group string, types []typeData) (err error) {
Break:
	for _, flow := range types {
		flow := flow.(flowData)
		for _, term := range flow.Terms {
			if e := w.Write(idl.Term, flow.Name,
				term.Name, term.Label, term.Type,
				term.Private, term.Optional, term.Repeats); e != nil {
				err = e
				break Break
			}
		}
	}
	return
}

// write str types to the db
func writeStr(w DB, group string, types []typeData) (err error) {
	for _, str := range types {
		str := str.(strData)
		var hasComments bool
		comments := str.OptionComments
		for _, c := range comments {
			if len(c) > 0 {
				hasComments = true
				break
			}
		}
		for i, opt := range str.Options {
			var comment any // let comments be nil if they are all empty
			if hasComments {
				var str string // or string if any of them exist
				if i < len(comments) {
					str = comments[i]
				}
				comment = str
			}
			if e := w.Write(idl.Enum, str.Name, opt, comment); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func writeTypes(w DB, group, spec string, types []typeData) (err error) {
Break:
	for _, t := range types {
		n := t.getName()
		if e := w.Write(idl.Op, n, group, spec); e != nil {
			err = e
			break
		} else {
			m := t.getMarkup()
			for _, k := range m.Keys() {
				v := m[k]
				var b strings.Builder
				if e := files.JsonEncoder(&b, files.RawJson).Encode(v); e != nil {
					err = e
					break Break
				} else if e := w.Write(idl.Markup, n, k, b.String()); e != nil {
					err = fmt.Errorf("%w couldnt write markup %q ", e, k)
					break Break
				}
			}
		}
	}
	return
}
