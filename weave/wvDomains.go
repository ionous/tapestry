package weave

import (
	"database/sql"
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

type Domain struct {
	name       string
	cat        *Catalog
	currPhase  assert.Phase                     // updated during weave, ends at NumPhases
	scheduling [assert.RequireAll + 1][]memento // separates commands into phases
}

type memento struct {
	cb func(*Weaver) error
	at string
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
			} else if d != uses && uses.currPhase <= assert.RequireAll {
				okay = false
				break
			}
		}
	}
	return
}

func (d *Domain) schedule(at string, when assert.Phase, what func(*Weaver) error) (err error) {
	if d.currPhase > when {
		ctx := Weaver{Domain: d, Phase: d.currPhase, Runtime: d.cat.run}
		err = what(&ctx)
	} else {
		d.scheduling[when] = append(d.scheduling[when], memento{what, at})
	}
	return
}

func (d *Domain) GetClosestNoun(name string) (ret string, err error) {
	if x, e := d.getClosestNoun(name); e != nil {
		err = e
	} else {
		ret = x.name
	}
	return
}

// find the noun with the closest name in this scope
// skips aliases for the sake of backwards compatibility:
// there should be a difference between "a noun is known as"
// and "understand this word by the player as" -- and currently there's not.
func (d *Domain) getClosestNoun(name string) (ret struct{ name, domain string }, err error) {
	if e := d.cat.db.QueryRow(`
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
	limit 1`, d.name, name).Scan(&ret.name, &ret.domain); e == sql.ErrNoRows {
		err = errutil.Fmt("%w couldn't find a noun named %s", mdl.Missing, name)
	} else if e != nil {
		err = errutil.New("database error", e)
	}
	return
}

// use the domain rules ( and hierarchy ) to strip determiners off of the passed word
func (d *Domain) StripDeterminer(word string) (retDet, retWord string) {
	// fix: determiners should be specified by the author ( and libraries )
	return lang.SliceArticle(word)
}

// use the domain rules ( and hierarchy ) to strip determiners off of the passed word
func (d *Domain) UniformDeterminer(word string) (retDet, retWord string) {
	// fix: determiners should be specified by the author ( and libraries )
	det, name := lang.SliceArticle(word)
	if name, ok := UniformString(name); ok {
		retDet, retWord = det, name
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

	// generate additional names by splitting the lowercase uniform name on the underscores:
	split := strings.FieldsFunc(noun, lang.IsBreak)
	if cnt := len(split); cnt > 1 {
		// in case the name was reduced due to multiple separators
		if breaks := strings.Join(split, "_"); breaks != noun {
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
		if e := cat.AddName(d.name, noun, name, i, at); e != nil && !errors.Is(e, mdl.Duplicate) {
			err = e
			break
		}
	}
	return
}

func (d *Domain) runPhase(ctx *Weaver) (err error) {
	w := ctx.Phase
	d.currPhase = w // hrmm
	redo := struct {
		cnt int
		err error
	}{}
	// don't range over the slice since the contents can change during traversal.
	// tbd: have "Schedule" immediately execute the statement if in the correct phase?
	els := &d.scheduling[w]
Loop:
	for len(*els) > 0 {
		// slice the next element out of the list
		next := (*els)[0]
		(*els) = (*els)[1:]

		switch e := next.call(ctx); {
		case e == nil:
			redo.cnt, redo.err = 0, nil
		case errors.Is(e, mdl.Missing):
			redo.err = errutil.Append(redo.err, e)
			if redo.cnt < len((*els)) {
				// add redo elements back into the list
				(*els) = append((*els), next)
				redo.cnt++
			} else {
				if d.cat.warn != nil {
					e := errutil.New(w, "didn't finish")
					d.cat.warn(e)
				}
				err = errutil.Append(err, redo.err)
				break Loop
			}
		case errors.Is(e, mdl.Duplicate):
			if d.cat.warn != nil {
				d.cat.warn(e)
			}
		default:
			err = errutil.Append(err, e)
		}
	}
	d.currPhase++
	return
}
