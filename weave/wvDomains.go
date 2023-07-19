package weave

import (
	"database/sql"
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

type Domain struct {
	name       string
	cat        *Catalog
	currPhase  Phase                     // updated during weave, ends at NumPhases
	scheduling [RequireAll + 1][]memento // separates commands into phases
	suspended  []memento                 // for missing definitions
}

type memento struct {
	cb    func(*Weaver) error
	at    string
	phase Phase
	err   error
}

func (d *Domain) Name() string {
	return d.name
}

func (op *memento) call(ctx *Weaver) error {
	ctx.At = op.at
	return op.cb(ctx)
}

// have all parent domains been processed?
func (d *Domain) isReadyForProcessing() (okay bool, err error) {
	cat := d.cat
	// get the domain hierarchy: the ancestors ending just before the domain itself.
	// direct parents may not be contiguous ( depending on whether their ancestors overlap. )
	if rows, e := cat.db.Query(`select uses from domain_tree 
		where base = ?1 order by dist desc`, d.name); e != nil {
		err = e
	} else if tree, e := tables.ScanStrings(rows); e != nil {
		err = e
	} else {
		okay = true // provisionally
		for _, name := range tree {
			if uses, ok := cat.domains[name]; !ok {
				okay = false
				break
			} else if (d != uses) && (uses.currPhase != -1) {
				okay = false
				break
			}
		}
	}
	return
}

func (d *Domain) schedule(at string, when Phase, what func(*Weaver) error) (err error) {
	if d.currPhase < 0 {
		err = errutil.Fmt("domain %q already finished", d.name)
	} else if d.currPhase <= when {
		d.scheduling[when] = append(d.scheduling[when], memento{what, at, when, nil})
	} else {
		ctx := Weaver{Catalog: d.cat, Domain: d, Phase: d.currPhase, Runtime: d.cat.run}
		if e := what(&ctx); errors.Is(e, mdl.Missing) {
			d.suspended = append(d.suspended, memento{what, at, when, e})
		} else {
			err = e
		}
	}
	return
}

// find the noun with the closest name in this scope
// skips aliases for the sake of backwards compatibility:
// there should be a difference between "a noun is known as"
// and "understand this word by the player as" -- and currently there's not.
func (d *Domain) GetExactNoun(name string) (ret *ScopedNoun, err error) {
	return d.findNoun(name, `
	select mn.noun, mn.domain  
	from mdl_noun mn
	join domain_tree dt
		on (dt.uses = mn.domain)
	where base = ?1
	and mn.noun = ?2
	limit 1`)
}

// find the noun with the closest name in this scope
// skips aliases for the sake of backwards compatibility:
// there should be a difference between "a noun is known as"
// and "understand this word by the player as" -- and currently there's not.
func (d *Domain) GetClosestNoun(name string) (ret *ScopedNoun, err error) {
	return d.findNoun(name, `
	select mn.noun, mn.domain  
	from mdl_name my 
	join mdl_noun mn
		on (mn.rowid = my.noun)
	join domain_tree dt
		on (dt.uses = my.domain)
	where base = ?1
	and my.name = ?2
	and my.rank >= 0
	order by my.rank, my.rowid asc
	limit 1`)
}

func (d *Domain) findNoun(name, q string) (ret *ScopedNoun, err error) {
	var noun struct{ name, domain string }
	if e := d.cat.db.QueryRow(q, d.name, name).Scan(&noun.name, &noun.domain); e == sql.ErrNoRows {
		err = errutil.Fmt("%w couldn't find a noun named %s", mdl.Missing, name)
	} else if e != nil {
		err = errutil.New("database error", e)
	} else if n, ok := d.cat.domainNouns[domainNoun{noun.domain, noun.name}]; !ok {
		err = errutil.Fmt("unexpected noun %q in domain %q", noun.name, noun.domain)
	} else {
		ret = n
	}
	return
}

// returns nil, no noun if it already existed
func (d *Domain) AddNoun(long, short, kind string) (ret *ScopedNoun, err error) {
	at := d.cat.cursor
	pen := d.cat.Pin(d.name, at)
	if e := pen.AddNoun(short, kind); e != nil {
		err = d.cat.eatDuplicates(e)
	} else if e := d.makeNames(short, long, at); e != nil {
		err = e
	} else {
		out := &ScopedNoun{domain: d, name: short}
		d.cat.domainNouns[domainNoun{d.name, short}] = out
		ret = out
	}
	return
}

func (d *Domain) makeNames(noun, name, at string) (err error) {
	cat := d.cat
	// if the original got transformed into underscores
	// write the original name (ex. "toy boat" vs "toy_boat" )
	var out []string
	if clip := strings.TrimSpace(name); clip != noun {
		out = append(out, clip)
	}
	out = append(out, noun)

	// generate additional names by splitting the name into parts
	split := lang.Fields(noun)
	if cnt := len(split); cnt > 1 {
		// in case the name was reduced due to multiple separators
		if breaks := strings.Join(split, " "); breaks != noun {
			out = append(out, breaks)
		}
		// write individual words in increasing rank ( ex. "boat", then "toy" )
		// note: trailing words are considered "stronger"
		// because adjectives in noun names tend to be first ( ie. "toy boat" )
		for i := len(split) - 1; i >= 0; i-- {
			word := split[i]
			out = append(out, word)
		}
	}
	for i, name := range out {
		// ignore duplicate errors here.
		// since these are generated, there's probably very little the user could do about them.
		if e := cat.Pin(d.name, at).AddName(noun, name, i); e != nil {
			err = e
			break
		}
	}
	return
}

func (d *Domain) runPhase(ctx *Weaver) (err error) {
	phase := ctx.Phase
	d.currPhase = phase // hrmm
	// don't range over the slice since the contents can change during traversal.
	// tbd; may no longer be true.
	els := &d.scheduling[phase]

	for len(*els) > 0 {
		// slice the next element out of the list
		next := (*els)[0]
		(*els) = (*els)[1:]

		switch e := next.call(ctx); {
		case errors.Is(e, mdl.Missing):
			next.err = e
			d.suspended = append(d.suspended, next)

		case e != nil:
			err = errutil.Append(err, e)
		}
	}
	d.currPhase++
	return
}

// tbd: the suspended mnd flush model is a little low bar.
// it'd be better to categorize individual statements -- even from macros
// and trigger them as their needs are satisfied; how exactly? not sure.
// the big mish mash are patterns,params,locals, and returns --
// which need to be written in a specific order.
func (d *Domain) flush(ignore bool) (err error) {
	ctx := Weaver{Catalog: d.cat, Domain: d, Runtime: d.cat.run}
	redo := struct {
		cnt int
		err error
	}{}

Loop:
	for len(d.suspended) > 0 {
		// slice the next element out of the list
		next := d.suspended[0]
		d.suspended = d.suspended[1:]
		ctx.Phase = next.phase

		switch e := next.call(&ctx); {
		case e == nil:
			// every success, abandon all old errors and try everything over again.
			redo.cnt, redo.err = 0, nil

		case errors.Is(e, mdl.Missing):
			// append to rack all that are missing
			redo.err = errutil.Append(redo.err, e)
			// add redo elements back into the list
			next.err = e
			d.suspended = append(d.suspended, next)
			// might still have statements to try?
			// keep going
			if redo.cnt = redo.cnt + 1; redo.cnt > len(d.suspended) {
				// if we have visited every suspended element
				// an haven't progressed; we're done.
				// return all the errors.
				if !ignore {
					err = redo.err
				}
				break Loop
			}

		default:
			// accumulate all errors
			err = errutil.Append(err, e)
		}
	}
	return
}
